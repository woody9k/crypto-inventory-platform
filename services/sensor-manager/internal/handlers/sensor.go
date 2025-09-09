package handlers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/democorp/crypto-inventory/services/sensor-manager/internal/models"
	"github.com/gin-gonic/gin"
)

// RegisterSensor handles sensor registration
func (h *Handler) RegisterSensor(c *gin.Context) {
	var req models.SensorRegistration
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Validate registration key against tenant
	// For now, we'll accept any registration key

	// Generate sensor ID
	sensorID := fmt.Sprintf("sensor-%s-%d", req.Platform, time.Now().Unix())

	// Generate client certificate
	clientCert, clientKey, err := h.generateClientCertificate(sensorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate client certificate"})
		return
	}

	// Generate server CA certificate (simplified for development)
	serverCACert, err := h.generateServerCACertificate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate server CA certificate"})
		return
	}

	// Create sensor record
	sensor := models.Sensor{
		ID:                sensorID,
		TenantID:          "default-tenant", // TODO: Extract from registration key
		Name:              req.Name,
		Description:       req.Description,
		Status:            "active",
		LastSeen:          time.Now(),
		IPAddress:         c.ClientIP(),
		Platform:          req.Platform,
		Version:           req.Version,
		Profile:           req.Profile,
		NetworkInterfaces: req.NetworkInterfaces,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// TODO: Save sensor to database
	// For now, we'll just return the response

	// Create configuration
	config := models.SensorConfig{
		ControlPlaneURL:   "https://localhost:8085",
		ReportingInterval: 30,
		StorageConfig: models.StorageConfig{
			MaxStorageSize: 100 * 1024 * 1024, // 100MB
			RotationSize:   10 * 1024 * 1024,  // 10MB
			RetentionDays:  7,
			EncryptionKey:  h.generateEncryptionKey(),
		},
		CaptureConfig: models.CaptureConfig{
			Interfaces:       req.NetworkInterfaces,
			ActiveProbing:    true,
			NetworkDiscovery: true,
			MaxConnections:   1000,
			TimeoutSeconds:   30,
		},
		Features: map[string]bool{
			"tls_analysis":         true,
			"ssh_analysis":         true,
			"certificate_analysis": true,
			"active_probing":       true,
			"network_discovery":    true,
		},
	}

	response := models.SensorRegistrationResponse{
		SensorID:     sensorID,
		ClientCert:   clientCert,
		ClientKey:    clientKey,
		ServerCACert: serverCACert,
		Config:       config,
	}

	c.JSON(http.StatusOK, response)
}

// SubmitDiscoveries handles discovery data submission from sensors
func (h *Handler) SubmitDiscoveries(c *gin.Context) {
	sensorID := c.Param("sensor_id")

	var batch models.DiscoveryBatch
	if err := c.ShouldBindJSON(&batch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate sensor ID
	if batch.SensorID != sensorID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sensor ID mismatch"})
		return
	}

	// TODO: Process discoveries and store in database
	// For now, we'll just acknowledge receipt

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"message":  "Discoveries received",
		"count":    len(batch.Discoveries),
		"batch_id": batch.BatchID,
	})
}

// ReportHealth handles sensor health reporting
func (h *Handler) ReportHealth(c *gin.Context) {
	sensorID := c.Param("sensor_id")

	var health models.SensorHealth
	if err := c.ShouldBindJSON(&health); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate sensor ID
	if health.SensorID != sensorID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sensor ID mismatch"})
		return
	}

	// TODO: Update sensor health in database
	// For now, we'll just acknowledge receipt

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Health report received",
	})
}

// GetSensorConfig handles sensor configuration requests
func (h *Handler) GetSensorConfig(c *gin.Context) {
	sensorID := c.Param("sensor_id")

	// TODO: Load configuration from database
	// For now, return default configuration

	config := models.SensorConfig{
		ControlPlaneURL:   "https://localhost:8085",
		ReportingInterval: 30,
		StorageConfig: models.StorageConfig{
			MaxStorageSize: 100 * 1024 * 1024, // 100MB
			RotationSize:   10 * 1024 * 1024,  // 10MB
			RetentionDays:  7,
			EncryptionKey:  h.generateEncryptionKey(),
		},
		CaptureConfig: models.CaptureConfig{
			Interfaces:       []string{"eth0", "wlan0"},
			ActiveProbing:    true,
			NetworkDiscovery: true,
			MaxConnections:   1000,
			TimeoutSeconds:   30,
		},
		Features: map[string]bool{
			"tls_analysis":         true,
			"ssh_analysis":         true,
			"certificate_analysis": true,
			"active_probing":       true,
			"network_discovery":    true,
		},
	}

	c.JSON(http.StatusOK, config)
}

// Helper function to generate client certificate
func (h *Handler) generateClientCertificate(sensorID string) (string, string, error) {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   sensorID,
			Organization: []string{"Crypto Inventory Sensor"},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(365 * 24 * time.Hour), // 1 year
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", "", err
	}

	// Encode certificate
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	// Encode private key
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	return string(certPEM), string(keyPEM), nil
}

// Helper function to generate server CA certificate
func (h *Handler) generateServerCACertificate() (string, error) {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", err
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "Crypto Inventory CA",
			Organization: []string{"Crypto Inventory Platform"},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(10 * 365 * 24 * time.Hour), // 10 years
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IsCA:        true,
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", err
	}

	// Encode certificate
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	return string(certPEM), nil
}

// Helper function to generate encryption key
func (h *Handler) generateEncryptionKey() string {
	key := make([]byte, 32) // 256 bits
	rand.Read(key)
	return fmt.Sprintf("%x", key)
}
