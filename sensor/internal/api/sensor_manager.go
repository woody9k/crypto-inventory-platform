package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/democorp/crypto-inventory/sensor/internal/config"
	"github.com/democorp/crypto-inventory/sensor/internal/models"
)

// SensorManagerClient handles communication with the sensor manager service
type SensorManagerClient struct {
	config     *config.Config
	httpClient *http.Client
	baseURL    string
}

// NewSensorManagerClient creates a new sensor manager client
func NewSensorManagerClient(cfg *config.Config) *SensorManagerClient {
	// Configure HTTP client with TLS if configured
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	if cfg.Security.UseTLS && cfg.Security.ClientCert != "" && cfg.Security.ClientKey != "" {
		// Load client certificate
		cert, err := tls.LoadX509KeyPair(cfg.Security.ClientCert, cfg.Security.ClientKey)
		if err != nil {
			fmt.Printf("Warning: Failed to load client certificate: %v\n", err)
		} else {
			// Configure TLS
			tlsConfig := &tls.Config{
				Certificates: []tls.Certificate{cert},
			}

			// Add server CA if provided
			if cfg.Security.ServerCACert != "" {
				// In a real implementation, you would load the CA certificate
				// and add it to the RootCAs
			}

			httpClient.Transport = &http.Transport{
				TLSClientConfig: tlsConfig,
			}
		}
	}

	return &SensorManagerClient{
		config:     cfg,
		httpClient: httpClient,
		baseURL:    cfg.ControlPlaneURL,
	}
}

// Register registers the sensor with the control plane
func (c *SensorManagerClient) Register() (*models.SensorConfig, error) {
	registration := models.SensorRegistration{
		RegistrationKey:   c.config.RegistrationKey,
		Name:              c.config.Name,
		Description:       c.config.Description,
		Platform:          c.config.Platform,
		Version:           c.config.Version,
		Profile:           c.config.Profile,
		NetworkInterfaces: c.config.Capture.Interfaces,
	}

	jsonData, err := json.Marshal(registration)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal registration: %v", err)
	}

	url := fmt.Sprintf("%s/api/v1/sensors/register", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send registration request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("registration failed with status %d: %s", resp.StatusCode, string(body))
	}

	var registrationResp struct {
		SensorID     string              `json:"sensor_id"`
		ClientCert   string              `json:"client_cert"`
		ClientKey    string              `json:"client_key"`
		ServerCACert string              `json:"server_ca_cert"`
		Config       models.SensorConfig `json:"config"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&registrationResp); err != nil {
		return nil, fmt.Errorf("failed to decode registration response: %v", err)
	}

	// Update sensor ID in config
	c.config.SensorID = registrationResp.SensorID

	// Update security config if certificates were provided
	if registrationResp.ClientCert != "" {
		c.config.Security.ClientCert = registrationResp.ClientCert
	}
	if registrationResp.ClientKey != "" {
		c.config.Security.ClientKey = registrationResp.ClientKey
	}
	if registrationResp.ServerCACert != "" {
		c.config.Security.ServerCACert = registrationResp.ServerCACert
	}

	return &registrationResp.Config, nil
}

// SubmitDiscoveries submits a batch of discoveries to the control plane
func (c *SensorManagerClient) SubmitDiscoveries(discoveries []*models.CryptoDiscovery) error {
	if len(discoveries) == 0 {
		return nil
	}

	batch := models.DiscoveryBatch{
		SensorID:    c.config.SensorID,
		Discoveries: make([]models.CryptoDiscovery, len(discoveries)),
		BatchID:     generateBatchID(),
		Timestamp:   time.Now(),
		Count:       len(discoveries),
	}

	// Convert pointers to values
	for i, discovery := range discoveries {
		batch.Discoveries[i] = *discovery
	}

	jsonData, err := json.Marshal(batch)
	if err != nil {
		return fmt.Errorf("failed to marshal discoveries: %v", err)
	}

	url := fmt.Sprintf("%s/api/v1/sensors/%s/discoveries", c.baseURL, c.config.SensorID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send discoveries: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("submission failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// ReportHealth reports sensor health to the control plane
func (c *SensorManagerClient) ReportHealth(health *models.SensorHealth) error {
	jsonData, err := json.Marshal(health)
	if err != nil {
		return fmt.Errorf("failed to marshal health: %v", err)
	}

	url := fmt.Sprintf("%s/api/v1/sensors/%s/health", c.baseURL, c.config.SensorID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send health report: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("health report failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetConfig retrieves sensor configuration from the control plane
func (c *SensorManagerClient) GetConfig() (*models.SensorConfig, error) {
	url := fmt.Sprintf("%s/api/v1/sensors/%s/config", c.baseURL, c.config.SensorID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("config request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var config models.SensorConfig
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %v", err)
	}

	return &config, nil
}

// Helper function to generate batch ID
