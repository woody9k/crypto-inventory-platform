package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/democorp/crypto-inventory/sensor/internal/api"
	"github.com/democorp/crypto-inventory/sensor/internal/capture"
	"github.com/democorp/crypto-inventory/sensor/internal/config"
	"github.com/democorp/crypto-inventory/sensor/internal/models"
	"github.com/democorp/crypto-inventory/sensor/internal/storage"
)

const Version = "1.0.0"

type Sensor struct {
	config        *config.Config
	packetCapture *capture.PacketCapture
	storage       *storage.EncryptedStorage
	apiClient     *api.OutboundClient
	discoveries   []*models.CryptoDiscovery
	mu            sync.RWMutex
}

func main() {
	// Command line flags
	var (
		version    = flag.Bool("version", false, "Show version information")
		configFile = flag.String("config", "config.yaml", "Path to configuration file")
		verbose    = flag.Bool("verbose", false, "Enable verbose logging")
		register   = flag.Bool("register", false, "Register with control plane")
	)
	flag.Parse()

	// Silence unused flag until file-based config is implemented
	_ = configFile

	// Show version and exit
	if *version {
		fmt.Printf("Crypto Inventory Network Sensor v%s\n", Version)
		fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		fmt.Printf("Go version: %s\n", runtime.Version())
		os.Exit(0)
	}

	// Initialize logging
	if *verbose {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("Verbose logging enabled")
	}

	log.Printf("🚀 Starting Crypto Inventory Network Sensor v%s", Version)
	log.Printf("Platform: %s/%s", runtime.GOOS, runtime.GOARCH)

	// Load configuration
	cfg := config.Load()
	log.Printf("Configuration loaded")

	// Create sensor instance
	sensor := &Sensor{
		config:      cfg,
		discoveries: make([]*models.CryptoDiscovery, 0),
	}

	// Initialize components
	if err := sensor.initialize(); err != nil {
		log.Fatalf("Failed to initialize sensor: %v", err)
	}

	// Register with control plane if requested
	if *register || cfg.RegistrationKey != "" {
		if err := sensor.register(); err != nil {
			log.Fatalf("Failed to register sensor: %v", err)
		}
	}

	// Start sensor
	if err := sensor.start(); err != nil {
		log.Fatalf("Failed to start sensor: %v", err)
	}

	// Handle graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Main sensor loop
	ticker := time.NewTicker(cfg.ReportingInterval)
	defer ticker.Stop()

	log.Println("✅ Sensor started successfully")
	log.Println("📡 Monitoring network traffic for cryptographic implementations...")

	for {
		select {
		case <-ticker.C:
			sensor.processDiscoveries()
		case discovery := <-sensor.packetCapture.GetDiscoveries():
			sensor.handleDiscovery(discovery)
		case err := <-sensor.packetCapture.GetErrors():
			log.Printf("❌ Capture error: %v", err)
		case sig := <-signalChan:
			log.Printf("Received signal %v, shutting down...", sig)
			sensor.cleanup()
			return
		}
	}
}

// initialize initializes all sensor components
func (s *Sensor) initialize() error {
	log.Println("🔧 Initializing sensor components...")

	// Initialize encrypted storage
	storage, err := storage.NewEncryptedStorage(s.config)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %v", err)
	}
	s.storage = storage

	// Initialize packet capture
	packetCapture := capture.NewPacketCapture(s.config)
	s.packetCapture = packetCapture

	// Initialize outbound-only API client
	apiClient := api.NewOutboundClient(s.config)
	s.apiClient = apiClient

	log.Println("✅ Sensor components initialized")
	return nil
}

// register registers the sensor with the control plane
func (s *Sensor) register() error {
	log.Println("📝 Registering with control plane...")

	config, err := s.apiClient.Register()
	if err != nil {
		return fmt.Errorf("registration failed: %v", err)
	}

	// Update sensor configuration with received config
	s.updateConfig(config)

	log.Printf("✅ Sensor registered successfully with ID: %s", s.config.SensorID)
	return nil
}

// start starts the sensor
func (s *Sensor) start() error {
	log.Println("🚀 Starting sensor...")

	// Start packet capture
	if err := s.packetCapture.Start(); err != nil {
		return fmt.Errorf("failed to start packet capture: %v", err)
	}

	log.Println("✅ Sensor started")
	return nil
}

// handleDiscovery handles a new discovery
func (s *Sensor) handleDiscovery(discovery *models.CryptoDiscovery) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Add to in-memory list
	s.discoveries = append(s.discoveries, discovery)

	// Store in encrypted storage
	if err := s.storage.StoreDiscovery(discovery); err != nil {
		log.Printf("❌ Failed to store discovery: %v", err)
	}

	log.Printf("🔍 Discovery: %s on %s:%d (confidence: %.2f)",
		discovery.Protocol, discovery.DestIP, discovery.Port, discovery.Confidence)
}

// processDiscoveries processes and sends discoveries to control plane
func (s *Sensor) processDiscoveries() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.discoveries) == 0 {
		return
	}

	// Send discoveries to control plane
	if err := s.apiClient.SubmitDiscoveries(s.discoveries); err != nil {
		log.Printf("❌ Failed to submit discoveries: %v", err)
		return
	}

	log.Printf("📤 Submitted %d discoveries to control plane", len(s.discoveries))

	// Clear discoveries after successful submission
	s.discoveries = s.discoveries[:0]

	// Send heartbeat and receive commands
	health := &models.SensorHealth{
		SensorID:        s.config.SensorID,
		Status:          "active",
		LastHeartbeat:   time.Now(),
		Uptime:          int64(time.Since(time.Now()).Seconds()),
		MemoryUsage:     getMemoryUsage(),
		CPUUsage:        getCPUUsage(),
		PacketsCaptured: 0, // TODO: Track actual packet count
		DiscoveriesMade: int64(len(s.discoveries)),
		Errors:          0, // TODO: Track actual error count
		Metrics:         make(map[string]interface{}),
		Timestamp:       time.Now(),
	}

	commands, err := s.apiClient.Heartbeat(health)
	if err != nil {
		log.Printf("❌ Failed to send heartbeat: %v", err)
	} else {
		// Process received commands
		s.processCommands(commands)
	}
}

// updateConfig updates sensor configuration
func (s *Sensor) updateConfig(config *models.SensorConfig) {
	// Update reporting interval
	s.config.ReportingInterval = time.Duration(config.ReportingInterval) * time.Second

	// Update storage config
	s.config.Storage.MaxStorageSize = config.StorageConfig.MaxStorageSize
	s.config.Storage.RotationSize = config.StorageConfig.RotationSize
	s.config.Storage.RetentionDays = config.StorageConfig.RetentionDays

	// Update capture config
	s.config.Capture.ActiveProbing = config.CaptureConfig.ActiveProbing
	s.config.Capture.NetworkDiscovery = config.CaptureConfig.NetworkDiscovery
	s.config.Capture.MaxConnections = config.CaptureConfig.MaxConnections
	s.config.Capture.TimeoutSeconds = config.CaptureConfig.TimeoutSeconds

	// Update features
	for feature, enabled := range config.Features {
		s.config.Features[feature] = enabled
	}
}

// cleanup performs cleanup operations
func (s *Sensor) cleanup() {
	log.Println("🧹 Performing cleanup...")

	// Stop packet capture
	if s.packetCapture != nil {
		s.packetCapture.Stop()
	}

	// Submit remaining discoveries
	if len(s.discoveries) > 0 {
		log.Printf("📤 Submitting %d remaining discoveries...", len(s.discoveries))
		if err := s.apiClient.SubmitDiscoveries(s.discoveries); err != nil {
			log.Printf("❌ Failed to submit remaining discoveries: %v", err)
		}
	}

	// Close storage
	if s.storage != nil {
		s.storage.Close()
	}

	log.Println("✅ Cleanup completed")
}

// processCommands processes commands received from control plane
func (s *Sensor) processCommands(commands *models.SensorCommands) {
	if len(commands.Commands) == 0 {
		return
	}

	log.Printf("📋 Processing %d commands from control plane", len(commands.Commands))

	for _, command := range commands.Commands {
		s.processCommand(command)
	}
}

// processCommand processes a single command
func (s *Sensor) processCommand(command models.Command) {
	log.Printf("🔧 Processing command: %s (type: %s, priority: %d)",
		command.ID, command.Type, command.Priority)

	switch command.Type {
	case "update_config":
		s.handleUpdateConfigCommand(command)
	case "restart":
		s.handleRestartCommand(command)
	case "stop":
		s.handleStopCommand(command)
	case "start_capture":
		s.handleStartCaptureCommand(command)
	case "stop_capture":
		s.handleStopCaptureCommand(command)
	default:
		log.Printf("⚠️ Unknown command type: %s", command.Type)
	}

	// Acknowledge command if required
	if command.RequiresAck {
		s.acknowledgeCommand(command)
	}
}

// handleUpdateConfigCommand handles configuration update commands
func (s *Sensor) handleUpdateConfigCommand(command models.Command) {
	if configData, ok := command.Payload["config"].(map[string]interface{}); ok {
		_ = configData
		// TODO: Parse and apply configuration updates
		log.Printf("📝 Configuration update command received")
	}
}

// handleRestartCommand handles restart commands
func (s *Sensor) handleRestartCommand(command models.Command) {
	log.Printf("🔄 Restart command received, scheduling restart...")
	// TODO: Implement graceful restart
}

// handleStopCommand handles stop commands
func (s *Sensor) handleStopCommand(command models.Command) {
	log.Printf("🛑 Stop command received, scheduling shutdown...")
	// TODO: Implement graceful shutdown
}

// handleStartCaptureCommand handles start capture commands
func (s *Sensor) handleStartCaptureCommand(command models.Command) {
	log.Printf("▶️ Start capture command received")
	// TODO: Start packet capture if not already running
}

// handleStopCaptureCommand handles stop capture commands
func (s *Sensor) handleStopCaptureCommand(command models.Command) {
	log.Printf("⏹️ Stop capture command received")
	// TODO: Stop packet capture
}

// acknowledgeCommand acknowledges a command to the control plane
func (s *Sensor) acknowledgeCommand(command models.Command) {
	// TODO: Send acknowledgment to control plane
	log.Printf("✅ Acknowledged command: %s", command.ID)
}

// Helper functions for system metrics
func getMemoryUsage() int64 {
	// TODO: Implement actual memory usage tracking
	return 0
}

func getCPUUsage() float64 {
	// TODO: Implement actual CPU usage tracking
	return 0.0
}
