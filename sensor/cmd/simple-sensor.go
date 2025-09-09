package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Simple sensor for demonstration purposes
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./simple-sensor --register")
		os.Exit(1)
	}

	if os.Args[1] == "--register" {
		registerSensor()
	} else {
		fmt.Println("Unknown command. Use --register")
		os.Exit(1)
	}
}

func registerSensor() {
	// Get registration details from environment or command line
	registrationKey := os.Getenv("REGISTRATION_KEY")
	if registrationKey == "" {
		registrationKey = "REG-tenant-1-20250909-910e7e" // Default for demo
	}

	controlPlaneURL := os.Getenv("CONTROL_PLANE_URL")
	if controlPlaneURL == "" {
		controlPlaneURL = "http://localhost:8080"
	}

	// Create registration payload
	payload := map[string]interface{}{
		"registration_key":   registrationKey,
		"name":               "sensor-demo",
		"description":        "Demo sensor for testing",
		"platform":           "linux",
		"version":            "1.0.0",
		"profile":            "datacenter_host",
		"network_interfaces": []string{"eth0"},
		"ip_address":         "192.168.1.100",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Failed to marshal payload: %v", err)
	}

	// Send registration request
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/sensors/register", controlPlaneURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		log.Fatalf("Failed to register sensor: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅ Sensor registered successfully!")

		// Parse response
		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
			if sensorID, ok := result["sensor_id"].(string); ok {
				fmt.Printf("Sensor ID: %s\n", sensorID)
			}
			if message, ok := result["message"].(string); ok {
				fmt.Printf("Message: %s\n", message)
			}
		}
	} else {
		fmt.Printf("❌ Registration failed with status: %d\n", resp.StatusCode)
		// Print error response
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			if errorMsg, ok := errorResp["error"].(string); ok {
				fmt.Printf("Error: %s\n", errorMsg)
			}
		}
	}
}
