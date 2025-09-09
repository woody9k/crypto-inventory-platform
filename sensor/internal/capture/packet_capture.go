package capture

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/democorp/crypto-inventory/sensor/internal/config"
	"github.com/democorp/crypto-inventory/sensor/internal/models"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// PacketCapture handles network packet capture
type PacketCapture struct {
	config      *config.Config
	interfaces  []string
	handles     []*pcap.Handle
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	discoveries chan *models.CryptoDiscovery
	errors      chan error
}

// NewPacketCapture creates a new packet capture instance
func NewPacketCapture(cfg *config.Config) *PacketCapture {
	ctx, cancel := context.WithCancel(context.Background())

	return &PacketCapture{
		config:      cfg,
		interfaces:  cfg.Capture.Interfaces,
		ctx:         ctx,
		cancel:      cancel,
		discoveries: make(chan *models.CryptoDiscovery, 1000),
		errors:      make(chan error, 100),
	}
}

// Start begins packet capture on all configured interfaces
func (pc *PacketCapture) Start() error {
	log.Printf("Starting packet capture on interfaces: %v", pc.interfaces)

	for _, iface := range pc.interfaces {
		if err := pc.startInterfaceCapture(iface); err != nil {
			log.Printf("Failed to start capture on interface %s: %v", iface, err)
			continue
		}
	}

	if len(pc.handles) == 0 {
		return fmt.Errorf("no interfaces available for capture")
	}

	// Start packet processing goroutines
	for i := 0; i < runtime.NumCPU(); i++ {
		pc.wg.Add(1)
		go pc.processPackets()
	}

	log.Printf("Packet capture started on %d interfaces", len(pc.handles))
	return nil
}

// Stop stops packet capture
func (pc *PacketCapture) Stop() {
	log.Println("Stopping packet capture...")
	pc.cancel()

	// Close all handles
	for _, handle := range pc.handles {
		handle.Close()
	}

	pc.wg.Wait()
	close(pc.discoveries)
	close(pc.errors)
	log.Println("Packet capture stopped")
}

// GetDiscoveries returns the discoveries channel
func (pc *PacketCapture) GetDiscoveries() <-chan *models.CryptoDiscovery {
	return pc.discoveries
}

// GetErrors returns the errors channel
func (pc *PacketCapture) GetErrors() <-chan error {
	return pc.errors
}

// startInterfaceCapture starts packet capture on a specific interface
func (pc *PacketCapture) startInterfaceCapture(iface string) error {
	// Check if interface exists
	interfaces, err := pcap.FindAllDevs()
	if err != nil {
		return fmt.Errorf("failed to find network interfaces: %v", err)
	}

	var found bool
	for _, dev := range interfaces {
		if dev.Name == iface {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("interface %s not found", iface)
	}

	// Open interface for capture
	handle, err := pcap.OpenLive(iface, int32(pc.config.Capture.BufferSize), true, pcap.BlockForever)
	if err != nil {
		return fmt.Errorf("failed to open interface %s: %v", iface, err)
	}

	// Set BPF filter for crypto-related traffic
	filter := "tcp port 443 or tcp port 22 or tcp port 993 or tcp port 995 or tcp port 465 or tcp port 587"
	if err := handle.SetBPFFilter(filter); err != nil {
		log.Printf("Warning: Failed to set BPF filter on %s: %v", iface, err)
	}

	pc.handles = append(pc.handles, handle)

	// Start capture goroutine for this interface
	pc.wg.Add(1)
	go pc.captureInterface(handle, iface)

	return nil
}

// captureInterface captures packets from a specific interface
func (pc *PacketCapture) captureInterface(handle *pcap.Handle, iface string) {
	defer pc.wg.Done()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for {
		select {
		case <-pc.ctx.Done():
			return
		case packet := <-packetSource.Packets():
			if packet == nil {
				continue
			}

			// Process packet in a separate goroutine to avoid blocking
			go pc.analyzePacket(packet, iface)
		}
	}
}

// processPackets processes captured packets
func (pc *PacketCapture) processPackets() {
	defer pc.wg.Done()

	for {
		select {
		case <-pc.ctx.Done():
			return
		case discovery := <-pc.discoveries:
			// Process discovery (this would typically send to storage or API)
			log.Printf("Discovery: %s on %s:%d", discovery.Protocol, discovery.DestIP, discovery.Port)
		}
	}
}

// analyzePacket analyzes a captured packet for crypto information
func (pc *PacketCapture) analyzePacket(packet gopacket.Packet, iface string) {
	// Extract network layer information
	networkLayer := packet.NetworkLayer()
	if networkLayer == nil {
		return
	}

	// Extract transport layer information
	transportLayer := packet.TransportLayer()
	if transportLayer == nil {
		return
	}

	// Get source and destination information
	srcIP := networkLayer.NetworkFlow().Src().String()
	dstIP := networkLayer.NetworkFlow().Dst().String()
	srcPort := transportLayer.TransportFlow().Src().String()
	dstPort := transportLayer.TransportFlow().Dst().String()

	// Analyze based on port
	port := getPortNumber(dstPort)
	protocol := getProtocolFromPort(port)

	if protocol == "" {
		return
	}

	// Create discovery record
	discovery := &models.CryptoDiscovery{
		ID:              generateDiscoveryID(),
		SensorID:        pc.config.SensorID,
		Timestamp:       time.Now(),
		SourceIP:        srcIP,
		DestIP:          dstIP,
		Port:            port,
		Protocol:        protocol,
		DiscoveryMethod: "passive",
		Confidence:      0.8, // Default confidence for passive detection
		RawMetadata: map[string]interface{}{
			"interface":   iface,
			"src_port":    srcPort,
			"packet_size": len(packet.Data()),
		},
		CreatedAt: time.Now(),
	}

	// Analyze packet content for crypto details
	pc.analyzeCryptoDetails(packet, discovery)

	// Send discovery to channel
	select {
	case pc.discoveries <- discovery:
	default:
		// Channel is full, log warning
		log.Printf("Warning: Discovery channel full, dropping discovery")
	}
}

// analyzeCryptoDetails analyzes packet content for cryptographic details
func (pc *PacketCapture) analyzeCryptoDetails(packet gopacket.Packet, discovery *models.CryptoDiscovery) {
	// Get application layer
	applicationLayer := packet.ApplicationLayer()
	if applicationLayer == nil {
		return
	}

	payload := applicationLayer.Payload()

	// Analyze based on protocol
	switch discovery.Protocol {
	case "TLS":
		pc.analyzeTLS(payload, discovery)
	case "SSH":
		pc.analyzeSSH(payload, discovery)
	}
}

// analyzeTLS analyzes TLS handshake data
func (pc *PacketCapture) analyzeTLS(payload []byte, discovery *models.CryptoDiscovery) {
	if len(payload) < 5 {
		return
	}

	// Check for TLS handshake
	if payload[0] != 0x16 { // TLS handshake
		return
	}

	// Extract TLS version
	if len(payload) >= 3 {
		version := (uint16(payload[1]) << 8) | uint16(payload[2])
		discovery.Version = getTLSVersion(version)
	}

	// Look for ClientHello or ServerHello
	if len(payload) >= 5 {
		handshakeType := payload[5]
		switch handshakeType {
		case 0x01: // ClientHello
			pc.analyzeClientHello(payload, discovery)
		case 0x02: // ServerHello
			pc.analyzeServerHello(payload, discovery)
		case 0x0B: // Certificate
			pc.analyzeCertificate(payload, discovery)
		}
	}
}

// analyzeSSH analyzes SSH protocol data
func (pc *PacketCapture) analyzeSSH(payload []byte, discovery *models.CryptoDiscovery) {
	if len(payload) < 4 {
		return
	}

	// Check for SSH protocol identifier
	if string(payload[:4]) == "SSH-" {
		// Extract SSH version
		end := 0
		for i, b := range payload {
			if b == '\n' || b == '\r' {
				end = i
				break
			}
		}
		if end > 4 {
			discovery.Version = string(payload[4:end])
		}
	}
}

// analyzeClientHello analyzes TLS ClientHello message
func (pc *PacketCapture) analyzeClientHello(payload []byte, discovery *models.CryptoDiscovery) {
	// This is a simplified analysis
	// In a real implementation, you would parse the full ClientHello structure
	discovery.RawMetadata["handshake_type"] = "ClientHello"
	discovery.Confidence = 0.9
}

// analyzeServerHello analyzes TLS ServerHello message
func (pc *PacketCapture) analyzeServerHello(payload []byte, discovery *models.CryptoDiscovery) {
	// This is a simplified analysis
	// In a real implementation, you would parse the full ServerHello structure
	discovery.RawMetadata["handshake_type"] = "ServerHello"
	discovery.Confidence = 0.9
}

// analyzeCertificate analyzes TLS Certificate message
func (pc *PacketCapture) analyzeCertificate(payload []byte, discovery *models.CryptoDiscovery) {
	// This is a simplified analysis
	// In a real implementation, you would parse the full Certificate structure
	discovery.RawMetadata["handshake_type"] = "Certificate"
	discovery.Confidence = 0.95
}

// Helper functions
func getPortNumber(portStr string) int {
	// Simple port parsing - in real implementation, use strconv.Atoi
	switch portStr {
	case "443":
		return 443
	case "22":
		return 22
	case "993":
		return 993
	case "995":
		return 995
	case "465":
		return 465
	case "587":
		return 587
	default:
		return 0
	}
}

func getProtocolFromPort(port int) string {
	switch port {
	case 443, 993, 995, 465, 587:
		return "TLS"
	case 22:
		return "SSH"
	default:
		return ""
	}
}

func getTLSVersion(version uint16) string {
	switch version {
	case 0x0301:
		return "TLS 1.0"
	case 0x0302:
		return "TLS 1.1"
	case 0x0303:
		return "TLS 1.2"
	case 0x0304:
		return "TLS 1.3"
	default:
		return "Unknown TLS"
	}
}

func generateDiscoveryID() string {
	// Simple ID generation - in real implementation, use proper UUID
	return fmt.Sprintf("discovery-%d", time.Now().UnixNano())
}
