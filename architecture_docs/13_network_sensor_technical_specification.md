# Network Sensor Technical Specification

## Document Purpose

This document provides comprehensive technical specifications for the Network Sensor agent, including deployment models, authentication mechanisms, resource constraints, and implementation architecture. This serves as the definitive reference for sensor development and integration with the primary platform interface.

## Executive Summary

The Network Sensor is a **cross-platform Go binary** deployed to enterprise networks for real-time cryptographic asset discovery. It operates under strict resource constraints (<100MB memory, <100MB storage) and supports multiple deployment scenarios from end-user machines to air-gapped data centers.

### Key Design Principles
- **Accuracy over Volume**: Prioritize accurate crypto analysis over capturing every network packet
- **Graceful Degradation**: Continue operating with reduced functionality rather than failing completely
- **Security First**: All data encrypted, minimal privileges, secure authentication
- **Cross-Platform**: Single codebase supporting Windows, Linux, macOS deployments
- **Flexible Deployment**: Support connected, periodic, and air-gapped operation modes

## Core Capabilities

### 1. Cryptographic Discovery Methods

#### Passive Detection (Primary Method)
- **TLS 1.0-1.2 Analysis**: Extract cipher suites, protocol versions, certificates from handshakes
- **Certificate Collection**: Capture full certificate chains during TLS negotiation
- **SSH Protocol Detection**: Identify SSH protocol versions and key exchange algorithms
- **Multi-Protocol Support**: IPSec, VPN, database encryption protocols
- **Partial Handshake Analysis**: Analyze incomplete captures rather than waiting for complete data

#### Active Interrogation (TLS 1.3 Fallback)
- **Service Discovery**: When passive detection reveals encrypted listeners without crypto details
- **Targeted Probing**: Connect to discovered services to capture:
  - Supported cipher suites via ClientHello manipulation
  - Certificate chains and protocol capabilities
  - Minimal footprint connections (immediate disconnect after analysis)
- **Configurable Aggressiveness**: Adjustable connection frequency and timeout values

### 2. Network Topology Discovery

#### Local Network Analysis
- **ARP Table Inspection**: Discover devices on local subnet
- **Routing Table Analysis**: Identify reachable network segments
- **DHCP Lease Detection**: Understand network ranges (where accessible)
- **DNS Traffic Analysis**: Passive DNS monitoring for service discovery

#### Multi-Network Detection
- **VLAN Awareness**: Detect VLAN tagging and multiple broadcast domains
- **Gateway Discovery**: Identify routing infrastructure and network boundaries
- **Segment Mapping**: Build topology maps for control plane optimization
- **Expansion Recommendations**: Suggest additional sensor placements for coverage gaps

## Resource Constraints & Optimization

### Memory Budget (<100MB Total)
```
Component Allocation:
- Packet Ring Buffer:     8MB  (streaming packet processing)
- Connection Tracker:    10MB  (active connection state)
- Certificate Cache:      5MB  (deduplication and analysis)
- Analysis Workspace:     5MB  (crypto processing buffer)
- Export Staging:         5MB  (compressed data preparation)
- Go Runtime + OS:       67MB  (system overhead)
```

### Storage Budget (<100MB Total)
```
Storage Allocation:
- Rotating Discovery Logs: 50MB (encrypted, auto-rotating)
- Certificate Cache:       30MB (metadata and cert chains)
- Network Topology Data:   15MB (discovered network information)
- Configuration & State:    5MB (sensor config and runtime state)
```

### CPU Optimization (Dual Core)
- **Core 1**: Packet capture and initial parsing
- **Core 2**: Crypto analysis, active probing, export processing
- **Shared Background**: Network discovery (low priority)

### Performance Strategies
- **Object Pooling**: Reuse packet and analysis objects to minimize allocations
- **Streaming Analysis**: Process packets immediately without buffering
- **LRU Caches**: Bounded caches for certificates and analysis results
- **Lock-free Queues**: Avoid contention between capture and analysis threads

## Deployment Architecture

### Deployment Profiles

#### End User Machine Profile
```yaml
profile: end_user_machine
network_interfaces: [primary]
active_probing: false          # Minimal footprint
storage_mode: connected
reporting_interval: 5m
network_discovery: false       # Don't scan from user machines
privilege_level: minimal
```

#### Datacenter Host Profile
```yaml
profile: datacenter_host
network_interfaces: [all]
active_probing: true
storage_mode: connected
reporting_interval: 1m
network_discovery: true        # Full network awareness
privilege_level: full
```

#### Cloud Instance Profile
```yaml
profile: cloud_instance
network_interfaces: [eth0, ens*]
active_probing: true
storage_mode: periodic         # Intermittent connectivity
reporting_interval: 15m
network_discovery: true
privilege_level: full
```

#### Air-Gapped Profile
```yaml
profile: air_gapped
network_interfaces: [all]
active_probing: true
storage_mode: airgapped
reporting_interval: 24h        # Daily export files
network_discovery: true
privilege_level: full
```

### Environment Auto-Detection

The sensor automatically detects deployment environment using:

#### Cloud Instance Detection
- Metadata endpoint availability (169.254.169.254)
- Hypervisor UUID presence (/sys/hypervisor/uuid)
- VM-specific filesystem indicators (/proc/xen/)

#### Datacenter Host Detection
- Multiple network interfaces (>2 NICs)
- Server-class hardware indicators
- Complex routing table entries
- Enterprise network segment patterns

#### End User Machine Detection
- Single primary network interface
- Consumer hardware characteristics
- Simple network configuration
- Desktop/laptop system indicators

#### Air-Gapped Detection
- No internet connectivity to known endpoints
- Missing cloud metadata services
- Isolated network characteristics

## Installation & Authentication

### Admin-Required Installation Process

Installation requires administrative privileges and includes:

1. **Environment Detection**: Automatic detection with admin override capability
2. **Interactive Configuration**: Profile selection and network interface configuration
3. **Registration Process**: Secure registration using UI-generated keys
4. **Service Installation**: Platform-specific service/daemon installation
5. **Certificate Management**: Mutual TLS certificate generation and storage

### Installation Command Interface

```bash
# Windows (requires Administrator)
crypto-sensor-installer.exe --install

# Linux (requires root or sudo)
sudo ./crypto-sensor-installer --install

# macOS (requires admin privileges)
sudo ./crypto-sensor-installer --install
```

### Interactive Installation Flow

```
Crypto Inventory Network Sensor Installer v1.0.0
===================================================

[1/6] Detecting Environment...
✓ Detected: Datacenter Host (Linux, multiple NICs)
  - Network Interfaces: eth0, eth1, eth2
  - Recommended Profile: datacenter_host
  
[2/6] Configuration
  Profile: datacenter_host (detected) [y/N to change]: 
  
[3/6] Network Configuration
  Monitor all interfaces (eth0, eth1, eth2)? [Y/n]: 
  Enable active probing? [Y/n]: 
  Enable network discovery? [Y/n]: 

[4/6] Control Plane Connection
  Control plane URL: https://crypto-inventory.company.com
  Registration key (from admin UI): REG-550e8400-20241215-A7B3C9

[5/6] Installing...
✓ Creating service user: crypto-sensor
✓ Installing binary: /opt/crypto-sensor/crypto-sensor
✓ Registering with control plane...
✓ Storing certificates: /opt/crypto-sensor/certs/
✓ Creating systemd service: crypto-sensor.service

[6/6] Starting Service...
✓ Service started successfully
✓ Sensor registered as: sensor-dc01-eth0-20241215
```

## Security & Authentication

### Registration Flow

#### Registration Key Format
```
Format: REG-{tenant_id}-{timestamp}-{checksum}
Example: REG-550e8400-20241215-A7B3C9

Components:
- tenant_id: 8-character tenant identifier
- timestamp: YYYYMMDD format  
- checksum: 6-character verification code
```

#### Initial Registration Process

1. **Admin generates registration key** in primary UI
2. **Installation process** includes registration key
3. **Sensor generates** unique identity and keypair
4. **Registration API call** with key validation
5. **Control plane issues** sensor-specific certificates
6. **Certificates stored** securely for ongoing authentication

### Mutual Authentication (mTLS)

#### Certificate-Based Authentication
- **Client Certificate**: Sensor identity certificate (issued during registration)
- **Server CA Certificate**: Control plane certificate authority
- **Mutual Verification**: Both sensor and control plane verify each other
- **Minimum TLS 1.2**: No support for older TLS versions
- **Approved Cipher Suites**: Limited to enterprise-approved encryption

#### Connection Security
```go
TLS Configuration:
- Minimum Version: TLS 1.2
- Cipher Suites: 
  - TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
  - TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
- Client Authentication: Required
- Certificate Verification: Full chain validation
```

## Encrypted Storage System

### Ephemeral Security Model

#### Encryption Strategy
- **Boot-time Key Generation**: New encryption key generated at each sensor restart
- **Memory-only Keys**: Encryption keys never written to persistent storage
- **AES-256-GCM**: All data encrypted with authenticated encryption
- **Secure Deletion**: Overwrite files before deletion (where OS permits)

#### File Rotation System
```go
Storage Structure:
- Current File: Active discovery log (up to 10MB)
- Historical Files: Up to 10 rotated files
- Auto-rotation: Every 10MB or 1 hour (whichever first)
- Total Limit: 100MB maximum storage usage
```

### Air-Gapped Export Format

#### Export File Structure
```go
type OfflineExport struct {
    EncryptedData []byte      // AES-256 encrypted discovery data
    Metadata      ExportMeta  // Export metadata and checksums
    Signature     []byte      // Digital signature for integrity
}

type ExportMeta struct {
    SensorID       string
    TimeRange      TimeRange
    RecordCount    int
    CompressionAlg string
    EncryptionAlg  string
    ExportVersion  string
}
```

#### Export Process
- **Daily Export Generation**: 24-hour batches for air-gapped environments
- **Compression**: Data compressed before encryption to minimize file size
- **Integrity Protection**: Digital signatures ensure export integrity
- **Manual Transfer**: Encrypted files designed for manual/removable media transfer

## Packet Capture & Analysis

### Cross-Platform Capture Strategy

#### Primary Approach: Pure Go Raw Sockets
- **No CGO Dependencies**: Single binary deployment without external libraries
- **Cross-Platform**: Native support for Windows, Linux, macOS, ARM
- **Privilege Requirements**: Requires administrative privileges for raw socket access
- **Fallback Support**: Alternative methods if raw sockets unavailable

#### Packet Processing Pipeline
```
Network Interface → Raw Socket → Packet Parser → Protocol Analyzer → Crypto Extractor → Discovery Record
```

### Accuracy-First Capture Philosophy

#### Priority-Based Processing
```yaml
Capture Priorities:
  Priority 1: HTTPS (port 443) - Web traffic analysis
  Priority 2: IMAPS (993), POP3S (995) - Email encryption
  Priority 3: SSH (port 22) - Secure shell analysis  
  Priority 4: TLS (any port) - Generic TLS detection
  Priority 5: Other encrypted protocols
```

#### Intelligent Filtering
- **High-Value Traffic**: Prioritize crypto handshakes over bulk data transfer
- **Deduplication**: Avoid analyzing identical crypto configurations repeatedly
- **Resource Protection**: Drop low-priority traffic under high load
- **Quality Metrics**: Maintain accuracy metrics for captured vs. analyzed traffic

### Partial Handshake Analysis

#### Analysis Strategies
```yaml
Partial Data Scenarios:
  Client Hello Only:
    - Extract supported cipher suites
    - Identify client capabilities
    - Confidence: 0.3
    
  Server Certificate Only:
    - Extract certificate details
    - Analyze certificate chain
    - Confidence: 0.6
    
  Server Hello + Certificate:
    - Full crypto configuration analysis
    - High confidence assessment
    - Confidence: 0.9
```

#### Timeout Management
- **Connection Timeout**: 30 seconds to capture complete handshake
- **Analysis Timeout**: Analyze partial data after timeout expiration
- **Memory Cleanup**: Purge incomplete connections after 5 minutes
- **Resource Protection**: Limit concurrent partial connection tracking

## Error Handling & Resilience

### Graceful Degradation Strategy

#### Component Failure Responses

**Packet Capture Failure**:
```yaml
Response: Switch to active-only mode
Actions:
  - Disable passive packet capture
  - Increase active probing frequency (2-minute intervals)
  - Log degradation state
  - Continue operation with reduced capability
```

**Control Plane Connectivity Loss**:
```yaml
Response: Switch to offline mode
Actions:
  - Enable local storage of all discoveries
  - Switch to air-gapped export format
  - Retry connection every 15 minutes
  - Maintain full functionality offline
```

**Storage System Errors**:
```yaml
Response: Emergency log rotation
Actions:
  - Force immediate log rotation
  - Compress and archive older data
  - Free space for continued operation
  - Alert control plane when connectivity available
```

#### Retry Mechanisms
- **Exponential Backoff**: Progressive delay increases for repeated failures
- **Circuit Breaker**: Temporary suspension of failing operations
- **Health Monitoring**: Continuous component health assessment
- **Automatic Recovery**: Resume normal operation when components recover

## Data Models & Discovery Records

### Core Discovery Record Structure

```go
type CryptoDiscovery struct {
    ID               string                 // Unique discovery identifier
    Timestamp        time.Time             // Discovery timestamp
    SensorID         string                // Source sensor identifier
    SourceIP         net.IP                // Source IP address
    DestIP           net.IP                // Destination IP address  
    Port             uint16                // Network port
    Protocol         string                // "TLS", "SSH", "IPSec"
    Version          string                // Protocol version
    CipherSuite      string                // Negotiated cipher suite
    Certificates     []CertificateInfo     // Certificate chain details
    DiscoveryMethod  string                // "passive" or "active"
    Confidence       float64               // Confidence score (0.0-1.0)
    RawMetadata      map[string]interface{} // Additional protocol-specific data
}

type CertificateInfo struct {
    SerialNumber     string
    Subject          string
    Issuer           string
    NotBefore        time.Time
    NotAfter         time.Time
    KeyAlgorithm     string
    KeySize          int
    SignatureAlg     string
    Fingerprint      string
    IsCA             bool
    Extensions       []string
}
```

### Network Topology Data

```go
type NetworkTopology struct {
    SensorID          string
    DiscoveredNetworks []NetworkSegment
    LocalInterfaces   []InterfaceInfo
    RoutingTable      []RouteInfo
    RecommendedSensors []SensorRecommendation
}

type NetworkSegment struct {
    Network     string    // CIDR notation (e.g., "192.168.1.0/24")
    Role        string    // "primary", "routed", "detected"
    DeviceCount int       // Estimated devices on segment
    AccessLevel string    // "full", "limited", "unknown"
}

type SensorRecommendation struct {
    Network        string
    Reason         string
    Priority       int
    EstimatedValue string
}
```

## Integration Points

### Control Plane Communication

#### API Endpoints
```yaml
Registration:
  POST /api/v1/sensors/register
  - Registration key validation
  - Certificate issuance
  - Initial configuration

Discovery Upload:
  POST /api/v1/sensors/{sensor_id}/discoveries
  - Batch discovery upload
  - Real-time discovery streaming
  - Compressed data support

Health Reporting:
  POST /api/v1/sensors/{sensor_id}/health
  - Sensor status updates
  - Performance metrics
  - Error reporting

Configuration:
  GET /api/v1/sensors/{sensor_id}/config
  - Dynamic configuration updates
  - Profile changes
  - Feature flag updates
```

#### Data Formats
- **JSON**: Primary data exchange format
- **Protocol Buffers**: High-performance binary format for large discovery batches
- **Compression**: gzip compression for all data transfers
- **Encryption**: All API communication over mTLS

### Primary Interface Integration

#### PCAP Ingestion Support
**Note**: PCAP file ingestion capability to be implemented in primary interface (separate from sensor functionality)

**Sensor Export Compatibility**: 
- Sensor export files designed for ingestion by primary interface
- Common data format between live sensor feeds and offline PCAP analysis
- Unified discovery record structure for all crypto intelligence sources

#### UI Integration Requirements

**Sensor Management Interface**:
- Sensor fleet visualization and status monitoring
- Deployment wizard for new sensor installation
- Configuration management for deployed sensors
- Performance metrics and health dashboards

**Discovery Data Presentation**:
- Real-time discovery feeds from active sensors
- Historical discovery data analysis and reporting
- Network topology visualization based on sensor data
- Risk assessment integration with sensor-discovered assets

## Implementation Roadmap

### Phase 1: Core Sensor Framework
- [ ] Cross-platform packet capture implementation
- [ ] Basic TLS handshake analysis
- [ ] Registration and authentication system
- [ ] Encrypted storage system
- [ ] Service/daemon installation framework

### Phase 2: Advanced Discovery
- [ ] Multi-protocol crypto detection (SSH, IPSec, VPN)
- [ ] Active probing for TLS 1.3
- [ ] Partial handshake analysis
- [ ] Certificate chain analysis and validation
- [ ] Network topology discovery

### Phase 3: Enterprise Features
- [ ] Air-gapped deployment support
- [ ] Performance optimization and resource management
- [ ] Advanced error handling and resilience
- [ ] Compliance and audit logging
- [ ] Integration with primary interface

### Phase 4: Intelligence Features
- [ ] Edge AI anomaly detection
- [ ] Advanced network mapping
- [ ] Predictive sensor placement recommendations
- [ ] Integration with threat intelligence feeds

## Security Considerations

### Data Protection
- **Minimal Data Collection**: Only collect cryptographic metadata, never payload data
- **Encryption Everywhere**: All data encrypted at rest and in transit
- **Secure Key Management**: Keys generated and managed securely
- **Data Retention**: Configurable retention policies with automatic cleanup

### Network Security
- **Minimal Network Footprint**: Only necessary connections to control plane
- **Firewall Friendly**: Standard HTTPS outbound connections
- **Certificate Pinning**: Additional security for control plane connections
- **Intrusion Detection**: Sensor activities designed to be IDS-friendly

### Compliance Support
- **Audit Logging**: Comprehensive logs of all sensor activities
- **Data Sovereignty**: Support for regional data residency requirements
- **Access Controls**: Role-based access to sensor management
- **Compliance Frameworks**: Support for SOC2, ISO 27001, NIST frameworks

## Conclusion

This technical specification provides the foundation for developing a robust, secure, and scalable network sensor system. The design balances comprehensive crypto discovery capabilities with strict resource constraints and diverse deployment requirements.

The sensor architecture supports the full spectrum of enterprise deployment scenarios while maintaining security, performance, and operational simplicity. Integration with the primary platform interface enables comprehensive crypto asset management across both connected and air-gapped environments.

---

**Document Status**: ✅ **COMPLETE** - Ready for implementation and primary interface integration

**Last Updated**: 2024-12-15

**Related Documents**:
- `02_system_architecture.md` - Overall platform architecture
- `04_data_models.md` - Database schemas and data structures  
- `06_deployment_guide.md` - Deployment and operations
- `07_ai_agent_handoff_guide.md` - Complete project context
