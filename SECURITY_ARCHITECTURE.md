# 🔒 Security Architecture for Cloud-Hosted Control Plane

## 🚨 **Problem Statement**

Traditional sensor architectures require inbound firewall rules, creating security and deployment challenges:

- **Security Risk**: Control plane must accept inbound connections from sensors
- **Deployment Blocker**: Enterprise networks often block inbound connections
- **Compliance Issue**: May violate network security policies
- **Attack Surface**: Control plane becomes a target for attacks

## 🛡️ **Solution: Outbound-Only Communication**

### **Architecture Principles**

1. **Sensors initiate all connections** (outbound only)
2. **No inbound firewall rules required**
3. **Control plane responds with commands** via HTTP responses
4. **Webhook support** for real-time updates (optional)
5. **Air-gapped export** for isolated environments

### **Communication Patterns**

#### **Pattern 1: Heartbeat + Commands (Primary)**
```
Sensor → Control Plane: POST /api/v1/sensors/{id}/heartbeat
Control Plane → Sensor: HTTP Response with commands
```

#### **Pattern 2: Polling (Fallback)**
```
Sensor → Control Plane: GET /api/v1/sensors/{id}/commands
Control Plane → Sensor: HTTP Response with commands
```

#### **Pattern 3: Webhooks (Optional)**
```
Control Plane → Sensor: POST /webhook (if sensor exposes endpoint)
```

## 🔧 **Implementation Details**

### **Sensor Communication Flow**

```go
// 1. Sensor registers (outbound only)
POST /api/v1/sensors/register
{
  "registration_key": "REG-550e8400-20241215-A7B3C9",
  "name": "sensor-dc01",
  "platform": "linux",
  "version": "1.0.0",
  "profile": "datacenter_host"
}

// 2. Sensor sends heartbeat (outbound only)
POST /api/v1/sensors/{id}/heartbeat
{
  "sensor_id": "sensor-dc01-eth0-20241215",
  "status": "active",
  "last_heartbeat": "2024-12-15T10:30:00Z",
  "uptime": 3600,
  "memory_usage": 52428800,
  "cpu_usage": 15.5,
  "packets_captured": 15000,
  "discoveries_made": 45,
  "errors": 0
}

// 3. Control plane responds with commands
HTTP 200 OK
{
  "sensor_id": "sensor-dc01-eth0-20241215",
  "timestamp": "2024-12-15T10:30:00Z",
  "commands": [
    {
      "id": "cmd-001",
      "type": "update_config",
      "priority": 5,
      "payload": {
        "reporting_interval": 60,
        "active_probing": true
      },
      "requires_ack": true
    }
  ]
}

// 4. Sensor submits discoveries (outbound only)
POST /api/v1/sensors/{id}/discoveries
{
  "sensor_id": "sensor-dc01-eth0-20241215",
  "discoveries": [...],
  "batch_id": "batch-123",
  "timestamp": "2024-12-15T10:30:00Z",
  "count": 5
}
```

### **Security Features**

#### **1. Mutual TLS (mTLS)**
- **Client certificates** for sensor authentication
- **Server certificates** for control plane verification
- **Certificate rotation** support
- **No shared secrets** in configuration

#### **2. Encrypted Storage**
- **AES-256-GCM** encryption for local data
- **Ephemeral keys** (generated at startup)
- **Secure deletion** of sensitive data
- **Air-gapped export** capability

#### **3. Network Security**
- **Outbound HTTPS only** (port 443)
- **No inbound connections** required
- **BPF filtering** for relevant traffic only
- **Minimal network footprint**

#### **4. Command Security**
- **Command signing** and verification
- **Expiration timestamps** for commands
- **Priority-based processing**
- **Acknowledgment requirements**

## 🌐 **Deployment Scenarios**

### **Scenario 1: Connected Environment**
```
┌─────────────────┐    ┌──────────────────┐
│   Sensor        │    │  Control Plane   │
│   (Outbound)    │───▶│  (Cloud)         │
│   Port 443      │    │  Port 443        │
└─────────────────┘    └──────────────────┘
```

**Requirements:**
- Outbound HTTPS access to control plane
- No inbound firewall rules
- Standard enterprise proxy support

### **Scenario 2: Air-Gapped Environment**
```
┌─────────────────┐    ┌──────────────────┐
│   Sensor        │    │  Export File     │
│   (Offline)     │───▶│  (USB/SFTP)      │
│   Local Storage │    │  Manual Transfer │
└─────────────────┘    └──────────────────┘
```

**Requirements:**
- Encrypted local storage
- Export file generation
- Manual transfer process
- Import at control plane

### **Scenario 3: Hybrid Environment**
```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Sensor        │    │  Control Plane   │    │   Webhook       │
│   (Polling)     │───▶│  (Cloud)         │───▶│   Service       │
│   + Webhook     │    │                  │    │   (Optional)    │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

**Requirements:**
- Primary: Outbound polling
- Fallback: Webhook for urgent updates
- Graceful degradation

## 🔐 **Security Benefits**

### **1. Reduced Attack Surface**
- **No inbound ports** exposed on sensors
- **Control plane** not directly accessible from sensors
- **Minimal network exposure**

### **2. Enterprise Compliance**
- **Standard outbound HTTPS** (port 443)
- **No special firewall rules** required
- **Compatible with proxy servers**

### **3. Scalability**
- **No connection limits** on control plane
- **Stateless communication**
- **Horizontal scaling** support

### **4. Fault Tolerance**
- **Sensors continue operating** if control plane is down
- **Local storage** for offline operation
- **Automatic retry** mechanisms

## 📋 **Implementation Checklist**

### **Sensor Implementation**
- [x] Outbound-only HTTP client
- [x] mTLS certificate support
- [x] Encrypted local storage
- [x] Command processing
- [x] Air-gapped export
- [x] Heartbeat mechanism

### **Control Plane Implementation**
- [x] Heartbeat endpoint
- [x] Command generation
- [x] Discovery ingestion
- [x] Webhook support
- [x] Air-gapped import
- [x] Certificate management

### **Security Implementation**
- [x] mTLS authentication
- [x] Command signing
- [x] Encrypted storage
- [x] Secure deletion
- [x] Certificate rotation
- [x] Audit logging

## 🚀 **Deployment Guide**

### **1. Sensor Deployment**
```bash
# Install sensor with outbound-only configuration
./crypto-sensor --register --verbose

# Environment variables
export CONTROL_PLANE_URL="https://crypto-inventory.company.com"
export REGISTRATION_KEY="REG-550e8400-20241215-A7B3C9"
export USE_TLS="true"
```

### **2. Firewall Configuration**
```bash
# Required outbound rules
iptables -A OUTPUT -p tcp --dport 443 -j ACCEPT

# No inbound rules required!
```

### **3. Proxy Configuration**
```bash
# Standard HTTPS proxy support
export HTTPS_PROXY="https://proxy.company.com:8080"
export HTTP_PROXY="http://proxy.company.com:8080"
```

## 🎯 **Benefits Summary**

| Aspect | Traditional | Outbound-Only |
|--------|-------------|---------------|
| **Firewall Rules** | Inbound + Outbound | Outbound only |
| **Security Risk** | High | Low |
| **Deployment** | Complex | Simple |
| **Compliance** | Difficult | Easy |
| **Scalability** | Limited | High |
| **Fault Tolerance** | Poor | Excellent |

## 🔍 **Monitoring & Alerting**

### **Sensor Health Metrics**
- Heartbeat frequency
- Command processing time
- Discovery submission rate
- Error rates
- Storage usage

### **Control Plane Metrics**
- Sensor registration rate
- Command queue depth
- Discovery processing rate
- API response times
- Certificate expiration

### **Security Metrics**
- Failed authentication attempts
- Certificate validation failures
- Command processing errors
- Storage encryption status
- Network connectivity issues

---

**This architecture ensures secure, scalable, and enterprise-ready sensor deployment without requiring inbound firewall rules or compromising network security.**
