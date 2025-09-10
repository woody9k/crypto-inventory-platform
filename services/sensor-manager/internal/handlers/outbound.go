// Package handlers provides HTTP handlers for the sensor-manager service.
// This file contains handlers for outbound-only communication patterns,
// allowing sensors to poll for commands and submit data without requiring
// inbound firewall rules.
package handlers

import (
	"net/http"

	"github.com/democorp/crypto-inventory/services/sensor-manager/internal/models"
	"github.com/gin-gonic/gin"
)

// Heartbeat handles sensor heartbeat and returns commands (outbound-only).
// This endpoint allows sensors to report their status and receive commands
// without requiring inbound firewall rules. The sensor initiates the connection.
func (h *Handler) Heartbeat(c *gin.Context) {
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
	// TODO: Check for pending commands for this sensor
	// TODO: Generate commands based on sensor state

	// For now, return empty commands
	commands := models.SensorCommands{
		SensorID: sensorID,
		Commands: []models.Command{},
	}

	c.JSON(http.StatusOK, commands)
}

// PollCommands handles sensor polling for commands (outbound-only).
// Sensors can poll this endpoint to check for pending commands without
// requiring the control plane to initiate connections.
func (h *Handler) PollCommands(c *gin.Context) {
	sensorID := c.Param("sensor_id")

	// TODO: Retrieve pending commands for this sensor from database
	// TODO: Mark commands as delivered
	// TODO: Generate new commands based on sensor state

	// For now, return empty commands
	commands := models.SensorCommands{
		SensorID: sensorID,
		Commands: []models.Command{},
	}

	c.JSON(http.StatusOK, commands)
}

// AcknowledgeCommand handles command acknowledgments from sensors.
// This allows sensors to confirm they have received and processed commands,
// enabling proper command lifecycle management.
func (h *Handler) AcknowledgeCommand(c *gin.Context) {
	sensorID := c.Param("sensor_id")

	var response models.CommandResponse
	if err := c.ShouldBindJSON(&response); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate sensor ID
	if response.SensorID != sensorID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sensor ID mismatch"})
		return
	}

	// TODO: Update command status in database
	// TODO: Process command response
	// TODO: Generate follow-up commands if needed

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Command acknowledgment received",
	})
}

// GetWebhookConfig returns webhook configuration for sensors
func (h *Handler) GetWebhookConfig(c *gin.Context) {
	sensorID := c.Param("sensor_id")

	// TODO: Load webhook configuration from database
	// For now, return disabled webhook config
	webhookConfig := models.WebhookConfig{
		SensorID:   sensorID, // Use sensorID to avoid unused variable warning
		Enabled:    false,
		WebhookURL: "",
		Secret:     "",
		Events:     []string{},
		RetryCount: 3,
		Timeout:    30,
	}

	c.JSON(http.StatusOK, webhookConfig)
}

// SubmitAirGappedExport handles air-gapped export submissions
func (h *Handler) SubmitAirGappedExport(c *gin.Context) {
	sensorID := c.Param("sensor_id")

	var export models.AirGappedExport
	if err := c.ShouldBindJSON(&export); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate sensor ID
	if export.SensorID != sensorID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sensor ID mismatch"})
		return
	}

	// TODO: Validate export signature and checksum
	// TODO: Decrypt and process export data
	// TODO: Store discoveries in database

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "Air-gapped export received",
		"export_id": export.ExportID,
		"records":   len(export.Data),
	})
}
