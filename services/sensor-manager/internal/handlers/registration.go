// Package handlers provides HTTP handlers for the sensor-manager service.
// This file contains handlers for sensor registration and pending sensor management,
// including registration key generation, IP validation, and mTLS certificate creation.
package handlers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RegistrationRequest represents a sensor registration request
type RegistrationRequest struct {
	RegistrationKey   string   `json:"registration_key" binding:"required"`
	Name              string   `json:"name" binding:"required"`
	Description       string   `json:"description"`
	Platform          string   `json:"platform" binding:"required"`
	Version           string   `json:"version" binding:"required"`
	Profile           string   `json:"profile" binding:"required"`
	NetworkInterfaces []string `json:"network_interfaces" binding:"required"`
	IPAddress         string   `json:"ip_address" binding:"required"`
	Tags              []string `json:"tags"`
}

// RegistrationResponse represents the response to a registration request
type RegistrationResponse struct {
	SensorID          string          `json:"sensor_id"`
	RegistrationKey   string          `json:"registration_key"`
	ClientCert        string          `json:"client_cert"`
	ClientKey         string          `json:"client_key"`
	ServerCACert      string          `json:"server_ca_cert"`
	ControlPlaneURL   string          `json:"control_plane_url"`
	ReportingInterval int             `json:"reporting_interval"`
	Features          map[string]bool `json:"features"`
	Message           string          `json:"message"`
}

// PendingSensor represents a pending sensor registration
type PendingSensor struct {
	ID                string    `json:"id"`
	RegistrationKey   string    `json:"registration_key"`
	Name              string    `json:"name"`
	IPAddress         string    `json:"ip_address"`
	Tags              []string  `json:"tags"`
	Profile           string    `json:"profile"`
	NetworkInterfaces []string  `json:"network_interfaces"`
	CreatedAt         time.Time `json:"created_at"`
	ExpiresAt         time.Time `json:"expires_at"`
	Status            string    `json:"status"` // pending, used, expired
}

// AdminSettings represents admin configuration
type AdminSettings struct {
	KeyExpirationMinutes int  `json:"key_expiration_minutes"`
	MaxPendingSensors    int  `json:"max_pending_sensors"`
	RequireIPValidation  bool `json:"require_ip_validation"`
}

// In-memory storage for demonstration (in production, use database)
var (
	pendingSensors = make(map[string]*PendingSensor)
	adminSettings  = AdminSettings{
		KeyExpirationMinutes: 60,
		MaxPendingSensors:    50,
		RequireIPValidation:  true,
	}
)

// CreatePendingSensor creates a new pending sensor registration
// This endpoint is called from the web UI to generate a registration key
// for a sensor that will be installed later. The key is time-limited
// and bound to a specific IP address for security.
func (h *Handler) CreatePendingSensor(c *gin.Context) {
	var req struct {
		Name              string   `json:"name" binding:"required"`       // Human-readable sensor name
		IPAddress         string   `json:"ip_address" binding:"required"` // IP address for validation
		Tags              []string `json:"tags"`                          // Optional tags for grouping
		Profile           string   `json:"profile" binding:"required"`    // Deployment profile
		NetworkInterfaces []string `json:"network_interfaces"`            // Interfaces to monitor
		Description       string   `json:"description"`                   // Optional description
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validate IP address format
	if net.ParseIP(req.IPAddress) == nil {
		c.JSON(400, gin.H{"error": "Invalid IP address format"})
		return
	}

	// Check if we've reached the maximum pending sensors
	if len(pendingSensors) >= adminSettings.MaxPendingSensors {
		c.JSON(400, gin.H{"error": "Maximum pending sensors reached"})
		return
	}

	// Generate unique registration key
	keySuffix := make([]byte, 3)
	rand.Read(keySuffix)
	key := fmt.Sprintf("REG-tenant-123-%s-%s",
		time.Now().Format("20060102"),
		hex.EncodeToString(keySuffix))

	// Create pending sensor
	pendingSensor := &PendingSensor{
		ID:                uuid.New().String(),
		RegistrationKey:   key,
		Name:              req.Name,
		IPAddress:         req.IPAddress,
		Tags:              req.Tags,
		Profile:           req.Profile,
		NetworkInterfaces: req.NetworkInterfaces,
		CreatedAt:         time.Now(),
		ExpiresAt:         time.Now().Add(time.Duration(adminSettings.KeyExpirationMinutes) * time.Minute),
		Status:            "pending",
	}

	pendingSensors[key] = pendingSensor

	c.JSON(200, gin.H{
		"pending_sensor": pendingSensor,
		"installation_command": fmt.Sprintf(
			"sudo ./install-sensor.sh --key %s --ip %s --name %s --profile %s",
			key, req.IPAddress, req.Name, req.Profile,
		),
	})
}

// RegisterSensor handles sensor registration with IP validation
// This endpoint is called by the sensor during installation to register
// with the control plane. It validates the registration key, checks IP
// address binding, and returns mTLS certificates for secure communication.
func (h *Handler) RegisterSensor(c *gin.Context) {
	var req RegistrationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validate IP address format
	if net.ParseIP(req.IPAddress) == nil {
		c.JSON(400, gin.H{"error": "Invalid IP address format"})
		return
	}

	// Check if registration key exists and is valid
	pendingSensor, exists := pendingSensors[req.RegistrationKey]
	if !exists {
		c.JSON(400, gin.H{"error": "Invalid or expired registration key"})
		return
	}

	// Check if key has expired
	if time.Now().After(pendingSensor.ExpiresAt) {
		pendingSensor.Status = "expired"
		c.JSON(400, gin.H{"error": "Registration key has expired"})
		return
	}

	// Check if key has already been used
	if pendingSensor.Status == "used" {
		c.JSON(400, gin.H{"error": "Registration key has already been used"})
		return
	}

	// Validate IP address matches the pending sensor
	if adminSettings.RequireIPValidation && pendingSensor.IPAddress != req.IPAddress {
		c.JSON(400, gin.H{"error": "IP address does not match the registered IP address"})
		return
	}

	// Additional IP validation: Check if the requesting IP matches the registered IP
	clientIP := c.ClientIP()
	if adminSettings.RequireIPValidation && clientIP != req.IPAddress {
		// Check if the client IP is in the same subnet as the registered IP
		if !isIPInSameSubnet(clientIP, req.IPAddress) {
			c.JSON(400, gin.H{"error": "Requesting IP address does not match the registered IP address"})
			return
		}
	}

	// Mark key as used
	pendingSensor.Status = "used"

	// Generate mTLS certificates
	caCertPEM, caKeyPEM, err := h.generateCACertificate()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate CA certificate"})
		return
	}

	sensorCertPEM, sensorKeyPEM, err := h.generateSensorCertificate(req.Name, caCertPEM, caKeyPEM)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate sensor certificate"})
		return
	}

	// Create sensor configuration based on profile
	features := getProfileFeatures(req.Profile)
	reportingInterval := getProfileReportingInterval(req.Profile)

	// Create sensor record (in production, save to database)
	sensorID := uuid.New().String()

	response := RegistrationResponse{
		SensorID:          sensorID,
		RegistrationKey:   req.RegistrationKey,
		ClientCert:        string(sensorCertPEM),
		ClientKey:         string(sensorKeyPEM),
		ServerCACert:      string(caCertPEM),
		ControlPlaneURL:   "https://crypto-inventory.company.com",
		ReportingInterval: reportingInterval,
		Features:          features,
		Message:           "Sensor registered successfully",
	}

	c.JSON(200, response)
}

// GetPendingSensors returns all pending sensor registrations
func (h *Handler) GetPendingSensors(c *gin.Context) {
	var sensors []*PendingSensor
	for _, sensor := range pendingSensors {
		// Check if expired
		if time.Now().After(sensor.ExpiresAt) && sensor.Status == "pending" {
			sensor.Status = "expired"
		}
		sensors = append(sensors, sensor)
	}
	c.JSON(200, gin.H{"pending_sensors": sensors})
}

// DeletePendingSensor deletes a pending sensor registration
func (h *Handler) DeletePendingSensor(c *gin.Context) {
	key := c.Param("key")
	if _, exists := pendingSensors[key]; !exists {
		c.JSON(404, gin.H{"error": "Pending sensor not found"})
		return
	}
	delete(pendingSensors, key)
	c.JSON(200, gin.H{"message": "Pending sensor deleted successfully"})
}

// GetAdminSettings returns current admin settings
func (h *Handler) GetAdminSettings(c *gin.Context) {
	c.JSON(200, adminSettings)
}

// UpdateAdminSettings updates admin settings
func (h *Handler) UpdateAdminSettings(c *gin.Context) {
	var newSettings AdminSettings
	if err := c.ShouldBindJSON(&newSettings); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validate settings
	if newSettings.KeyExpirationMinutes < 5 || newSettings.KeyExpirationMinutes > 1440 {
		c.JSON(400, gin.H{"error": "Key expiration must be between 5 and 1440 minutes"})
		return
	}

	if newSettings.MaxPendingSensors < 1 || newSettings.MaxPendingSensors > 1000 {
		c.JSON(400, gin.H{"error": "Max pending sensors must be between 1 and 1000"})
		return
	}

	adminSettings = newSettings
	c.JSON(200, gin.H{"message": "Admin settings updated successfully"})
}

// Helper functions

func isIPInSameSubnet(ip1, ip2 string) bool {
	// Simple check - in production, implement proper subnet validation
	// For now, just check if they're the same
	return ip1 == ip2
}

func getProfileFeatures(profile string) map[string]bool {
	features := map[string]bool{
		"tls_analysis":         true,
		"ssh_analysis":         true,
		"certificate_analysis": true,
		"active_probing":       false,
		"network_discovery":    false,
		"air_gapped_export":    false,
	}

	switch profile {
	case "datacenter_host":
		features["active_probing"] = true
		features["network_discovery"] = true
	case "cloud_instance":
		features["active_probing"] = true
	case "end_user_machine":
		// Minimal features
	case "air_gapped":
		features["air_gapped_export"] = true
		features["active_probing"] = false
		features["network_discovery"] = false
	}

	return features
}

func getProfileReportingInterval(profile string) int {
	switch profile {
	case "datacenter_host":
		return 30 // 30 seconds
	case "cloud_instance":
		return 60 // 1 minute
	case "end_user_machine":
		return 300 // 5 minutes
	case "air_gapped":
		return 3600 // 1 hour
	default:
		return 60
	}
}

// generateCACertificate generates a CA certificate for mTLS
func (h *Handler) generateCACertificate() (string, string, error) {
	// Generate private key
	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	// Create CA certificate template
	caTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"Crypto Inventory CA"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0), // 10 years
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// Create CA certificate
	caCertDER, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		return "", "", err
	}

	// Encode CA certificate
	caCertPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caCertDER,
	})

	// Encode CA private key
	caKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caKey),
	})

	return string(caCertPEM), string(caKeyPEM), nil
}

// generateSensorCertificate generates a client certificate for a sensor
func (h *Handler) generateSensorCertificate(sensorName, caCertPEM, caKeyPEM string) (string, string, error) {
	// Parse CA certificate
	caBlock, _ := pem.Decode([]byte(caCertPEM))
	caCert, err := x509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		return "", "", err
	}

	// Parse CA private key
	caKeyBlock, _ := pem.Decode([]byte(caKeyPEM))
	caKey, err := x509.ParsePKCS1PrivateKey(caKeyBlock.Bytes)
	if err != nil {
		return "", "", err
	}

	// Generate sensor private key
	sensorKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	// Create sensor certificate template
	sensorTemplate := x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			Organization:  []string{"Crypto Inventory Sensor"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
			CommonName:    sensorName,
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(1, 0, 0), // 1 year
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		IPAddresses: []net.IP{},
		DNSNames:    []string{sensorName},
	}

	// Create sensor certificate
	sensorCertDER, err := x509.CreateCertificate(rand.Reader, &sensorTemplate, caCert, &sensorKey.PublicKey, caKey)
	if err != nil {
		return "", "", err
	}

	// Encode sensor certificate
	sensorCertPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: sensorCertDER,
	})

	// Encode sensor private key
	sensorKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(sensorKey),
	})

	return string(sensorCertPEM), string(sensorKeyPEM), nil
}
