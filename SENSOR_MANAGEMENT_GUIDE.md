# 🎛️ Complete Sensor Management Guide

## 📋 **Sensor Management Overview**

The crypto inventory platform provides comprehensive sensor management capabilities with outbound-only communication, eliminating the need for inbound firewall rules.

## 🔑 **Step 1: Registration Key Generation**

### **Generate Registration Key**
```bash
# Generate a registration key for a tenant
cd scripts
go run generate-registration-key.go tenant-123 "Datacenter Sensors" 10

# Output:
# 🔑 Sensor Registration Key Generated
# =====================================
# Key: REG-tenant-123-20241215-A7B3C9
# Tenant ID: tenant-123
# Description: Datacenter Sensors
# Expires: 2024-12-15 10:30:00
# Max Sensors: 10
```

### **Registration Key Format**
```
Format: REG-{tenant_id}-{timestamp}-{checksum}
Example: REG-tenant-123-20241215-A7B3C9

Components:
- tenant_id: 8-character tenant identifier
- timestamp: YYYYMMDD format
- checksum: 6-character verification code
```

## 🚀 **Step 2: Sensor Installation**

### **Automatic Installation**
```bash
# Basic installation with auto-detection
sudo ./scripts/install-sensor.sh --key REG-tenant-123-20241215-A7B3C9

# Custom configuration
sudo ./scripts/install-sensor.sh \
  --key REG-tenant-123-20241215-A7B3C9 \
  --name sensor-dc01 \
  --interfaces "eth0,eth1" \
  --profile datacenter_host \
  --url https://crypto-inventory.company.com
```

### **Installation Profiles**

#### **Datacenter Host Profile**
```bash
sudo ./scripts/install-sensor.sh \
  --key REG-tenant-123-20241215-A7B3C9 \
  --profile datacenter_host \
  --interfaces "eth0,eth1,eth2"
```
- **Features**: Full network discovery, active probing, multiple interfaces
- **Use Case**: Server rooms, data centers
- **Interfaces**: All available network interfaces

#### **Cloud Instance Profile**
```bash
sudo ./scripts/install-sensor.sh \
  --key REG-tenant-123-20241215-A7B3C9 \
  --profile cloud_instance \
  --interfaces "ens3"
```
- **Features**: Cloud-optimized, periodic reporting
- **Use Case**: AWS, Azure, GCP instances
- **Interfaces**: Primary cloud interface

#### **End User Machine Profile**
```bash
sudo ./scripts/install-sensor.sh \
  --key REG-tenant-123-20241215-A7B3C9 \
  --profile end_user_machine \
  --interfaces "wlan0"
```
- **Features**: Minimal footprint, single interface
- **Use Case**: Laptops, workstations
- **Interfaces**: Primary user interface

#### **Air-Gapped Profile**
```bash
sudo ./scripts/install-sensor.sh \
  --key REG-tenant-123-20241215-A7B3C9 \
  --profile air_gapped \
  --interfaces "eth0"
```
- **Features**: Offline mode, export files
- **Use Case**: Isolated networks, secure environments
- **Interfaces**: Single interface for local network

## ⚙️ **Step 3: Sensor Configuration**

### **Configuration File Location**
```
/opt/crypto-sensor/config.yaml
```

### **Configuration Options**
```yaml
# Sensor Identity
sensor:
  name: "sensor-dc01"
  platform: "linux"
  version: "1.0.0"
  profile: "datacenter_host"

# Control Plane Connection
control_plane:
  url: "https://crypto-inventory.company.com"
  registration_key: "REG-tenant-123-20241215-A7B3C9"

# Network Configuration
network:
  interfaces: ["eth0", "eth1"]
  active_probing: true
  network_discovery: true

# Storage Configuration
storage:
  data_path: "/opt/crypto-sensor/data"
  max_size: 100MB
  rotation_size: 10MB
  retention_days: 7

# Security Configuration
security:
  use_tls: true
  client_cert: "/opt/crypto-sensor/certs/client.crt"
  client_key: "/opt/crypto-sensor/certs/client.key"
  server_ca_cert: "/opt/crypto-sensor/certs/server-ca.crt"

# Feature Flags
features:
  tls_analysis: true
  ssh_analysis: true
  certificate_analysis: true
  active_probing: true
  network_discovery: true
  air_gapped_export: false
```

## 🎛️ **Step 4: Sensor Management**

### **Service Management**
```bash
# Check sensor status
systemctl status crypto-sensor

# View sensor logs
journalctl -u crypto-sensor -f

# Restart sensor
systemctl restart crypto-sensor

# Stop sensor
systemctl stop crypto-sensor

# Start sensor
systemctl start crypto-sensor
```

### **Sensor Health Monitoring**
```bash
# Check sensor health via API
curl https://crypto-inventory.company.com/api/v1/sensors/sensor-dc01/health

# View recent discoveries
curl https://crypto-inventory.company.com/api/v1/sensors/sensor-dc01/discoveries

# Get sensor configuration
curl https://crypto-inventory.company.com/api/v1/sensors/sensor-dc01/config
```

### **Web UI Management**
Access the sensor management interface at:
```
https://crypto-inventory.company.com/sensors
```

**Features:**
- **Sensor Dashboard**: Real-time status and metrics
- **Add/Remove Sensors**: Easy sensor lifecycle management
- **Configuration**: Update sensor settings
- **Monitoring**: View discoveries and health data
- **Commands**: Send remote commands to sensors

## 🔧 **Step 5: Remote Sensor Commands**

### **Command Types**
```json
{
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
    },
    {
      "id": "cmd-002",
      "type": "restart",
      "priority": 8,
      "payload": {},
      "requires_ack": true
    },
    {
      "id": "cmd-003",
      "type": "stop_capture",
      "priority": 3,
      "payload": {
        "reason": "maintenance"
      },
      "requires_ack": false
    }
  ]
}
```

### **Available Commands**
- **`update_config`**: Update sensor configuration
- **`restart`**: Restart sensor service
- **`stop`**: Stop sensor service
- **`start_capture`**: Start packet capture
- **`stop_capture`**: Stop packet capture
- **`export_data`**: Export discoveries for air-gapped transfer

## 🌐 **Step 6: Network Interface Configuration**

### **Interface Detection**
The installer automatically detects available network interfaces:
```bash
# Auto-detected interfaces
eth0, eth1, ens3, wlan0

# Manual specification
--interfaces "eth0,eth1,ens3"
```

### **Interface Types**
- **`eth*`**: Ethernet interfaces
- **`ens*`**: Cloud provider interfaces (AWS, Azure)
- **`enp*`**: PCIe interfaces
- **`wlan*`**: Wireless interfaces

### **Interface Filtering**
```bash
# Filter specific interfaces
--interfaces "eth0,eth1"  # Only monitor these interfaces

# Exclude interfaces
--interfaces "eth0"       # Exclude eth1, eth2, etc.
```

## 🔒 **Step 7: Security Configuration**

### **mTLS Certificates**
Sensors automatically receive certificates during registration:
```
/opt/crypto-sensor/certs/
├── client.crt      # Sensor client certificate
├── client.key      # Sensor private key
└── server-ca.crt   # Control plane CA certificate
```

### **Certificate Management**
```bash
# View certificate details
openssl x509 -in /opt/crypto-sensor/certs/client.crt -text -noout

# Check certificate expiration
openssl x509 -in /opt/crypto-sensor/certs/client.crt -noout -dates

# Renew certificates (automatic)
systemctl restart crypto-sensor
```

## 📊 **Step 8: Monitoring and Alerting**

### **Sensor Metrics**
- **Status**: active, inactive, error, unknown
- **Uptime**: Sensor running time
- **Memory Usage**: Current memory consumption
- **CPU Usage**: Current CPU utilization
- **Packets Captured**: Total packets analyzed
- **Discoveries Made**: Total crypto discoveries
- **Errors**: Error count and types

### **Health Checks**
```bash
# Sensor health endpoint
GET /api/v1/sensors/{sensor_id}/health

# Response
{
  "sensor_id": "sensor-dc01",
  "status": "active",
  "last_heartbeat": "2024-12-15T10:30:00Z",
  "uptime": 3600,
  "memory_usage": 52428800,
  "cpu_usage": 15.5,
  "packets_captured": 15000,
  "discoveries_made": 45,
  "errors": 0
}
```

### **Alerting Rules**
- **Sensor Offline**: No heartbeat for 5 minutes
- **High Error Rate**: >10 errors per hour
- **High Memory Usage**: >80% of allocated memory
- **Low Discovery Rate**: <1 discovery per hour

## 🚨 **Step 9: Troubleshooting**

### **Common Issues**

#### **Sensor Registration Failed**
```bash
# Check registration key validity
curl -X POST https://crypto-inventory.company.com/api/v1/sensors/register \
  -H "Content-Type: application/json" \
  -d '{"registration_key": "REG-tenant-123-20241215-A7B3C9", ...}'

# Check network connectivity
ping crypto-inventory.company.com
telnet crypto-inventory.company.com 443
```

#### **Sensor Not Capturing Packets**
```bash
# Check interface permissions
sudo setcap cap_net_raw+ep /opt/crypto-sensor/crypto-sensor

# Check interface status
ip link show
ip addr show

# Test packet capture
sudo tcpdump -i eth0 -c 10
```

#### **Certificate Issues**
```bash
# Check certificate validity
openssl x509 -in /opt/crypto-sensor/certs/client.crt -noout -dates

# Regenerate certificates
systemctl stop crypto-sensor
rm /opt/crypto-sensor/certs/*
systemctl start crypto-sensor
```

### **Log Analysis**
```bash
# View sensor logs
journalctl -u crypto-sensor -f

# Filter error logs
journalctl -u crypto-sensor --priority=err

# View specific time range
journalctl -u crypto-sensor --since "2024-12-15 10:00:00"
```

## 📋 **Step 10: Uninstallation**

### **Remove Sensor**
```bash
# Stop and disable service
sudo systemctl stop crypto-sensor
sudo systemctl disable crypto-sensor

# Remove service file
sudo rm /etc/systemd/system/crypto-sensor.service

# Remove installation directory
sudo rm -rf /opt/crypto-sensor

# Remove service user
sudo userdel crypto-sensor

# Reload systemd
sudo systemctl daemon-reload
```

### **Clean Uninstall Script**
```bash
#!/bin/bash
# Uninstall sensor completely
sudo systemctl stop crypto-sensor
sudo systemctl disable crypto-sensor
sudo rm /etc/systemd/system/crypto-sensor.service
sudo rm -rf /opt/crypto-sensor
sudo userdel crypto-sensor
sudo systemctl daemon-reload
echo "Sensor uninstalled successfully"
```

## 🎯 **Best Practices**

### **Deployment**
1. **Test in staging** before production deployment
2. **Use appropriate profiles** for each environment
3. **Monitor sensor health** regularly
4. **Keep certificates updated**
5. **Use descriptive sensor names**

### **Security**
1. **Rotate registration keys** regularly
2. **Monitor certificate expiration**
3. **Use strong network segmentation**
4. **Implement proper access controls**
5. **Audit sensor activities**

### **Performance**
1. **Monitor resource usage** (CPU, memory, disk)
2. **Adjust reporting intervals** based on network size
3. **Use appropriate interface filtering**
4. **Implement proper log rotation**
5. **Scale horizontally** for large networks

---

**This comprehensive sensor management system provides enterprise-grade network monitoring with outbound-only communication, eliminating security concerns while maintaining full functionality.**
