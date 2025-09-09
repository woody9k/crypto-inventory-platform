# 🚀 Agent Deployment Guide

## 📋 **Overview**

This guide covers the complete deployment process for the Crypto Inventory Network Agent (sensor). The agent is designed for secure, outbound-only communication with the control plane, eliminating the need for inbound firewall rules.

## 🏗️ **Architecture**

### **Agent Components**
- **Network Sensor**: Core packet capture and analysis engine
- **Configuration Manager**: Handles settings and profile management
- **Communication Client**: Outbound-only API client with mTLS
- **Storage Engine**: Encrypted local storage for discoveries
- **Command Processor**: Handles remote commands from control plane

### **Communication Pattern**
```
Agent → Control Plane (Outbound Only)
├── Registration (HTTPS POST)
├── Heartbeat (HTTPS POST)
├── Discovery Submission (HTTPS POST)
├── Command Polling (HTTPS GET)
└── Health Reporting (HTTPS POST)
```

## 🔧 **Installation Methods**

### **Method 1: Interactive Installation (Recommended)**

```bash
# Download and run in interactive mode
curl -sSL https://crypto-inventory.company.com/scripts/install-sensor.sh | sudo bash -s -- --interactive
```

**Interactive prompts:**
1. **Registration Key**: Enter the key from the web UI
2. **IP Address**: Specify the expected IP address
3. **Sensor Name**: Optional (auto-detected if not provided)
4. **Control Plane URL**: Optional (defaults to production)
5. **Profile Selection**: Choose deployment profile
6. **Network Interfaces**: Select interfaces to monitor
7. **Installation Directory**: Optional (defaults to `/opt/crypto-sensor`)

### **Method 2: One-Line Installation**

```bash
# Copy-paste command from web UI
curl -sSL https://crypto-inventory.company.com/scripts/install-sensor.sh | sudo bash -s -- \
  --key REG-tenant-123-20241215-A7B3C9 \
  --ip 192.168.1.100 \
  --name sensor-dc01 \
  --profile datacenter_host \
  --interfaces "eth0,eth1" \
  --url https://crypto-inventory.company.com
```

### **Method 3: Manual Installation**

```bash
# Download script manually
curl -O https://crypto-inventory.company.com/scripts/install-sensor.sh
chmod +x install-sensor.sh

# Run with parameters
sudo ./install-sensor.sh \
  --key REG-tenant-123-20241215-A7B3C9 \
  --ip 192.168.1.100 \
  --name sensor-dc01 \
  --profile datacenter_host
```

## 🎛️ **Deployment Profiles**

### **Datacenter Host Profile**
```bash
--profile datacenter_host
```
- **Features**: Full network discovery, active probing, multiple interfaces
- **Use Case**: Server rooms, data centers, high-traffic environments
- **Interfaces**: All available network interfaces
- **Reporting**: 30-second intervals

### **Cloud Instance Profile**
```bash
--profile cloud_instance
```
- **Features**: Cloud-optimized, periodic reporting
- **Use Case**: AWS, Azure, GCP instances
- **Interfaces**: Primary cloud interface (ens3, eth0)
- **Reporting**: 1-minute intervals

### **End User Machine Profile**
```bash
--profile end_user_machine
```
- **Features**: Minimal footprint, single interface
- **Use Case**: Laptops, workstations, office machines
- **Interfaces**: Primary user interface (wlan0, eth0)
- **Reporting**: 5-minute intervals

### **Air-Gapped Profile**
```bash
--profile air_gapped
```
- **Features**: Offline mode, export files
- **Use Case**: Isolated networks, secure environments
- **Interfaces**: Single interface for local network
- **Reporting**: 1-hour intervals

## 🔒 **Security Features**

### **IP Address Validation**
- Registration key is bound to specific IP address
- Agent validates IP during installation
- Prevents unauthorized key usage

### **mTLS Authentication**
- Mutual TLS between agent and control plane
- Certificates generated during registration
- Automatic certificate rotation

### **Outbound-Only Communication**
- No inbound connections required
- No firewall rules needed
- Works through NAT and proxies

### **Encrypted Storage**
- Local discoveries encrypted with AES-256-GCM
- Keys derived from registration credentials
- Secure file rotation and cleanup

## 📊 **Configuration**

### **Configuration File Location**
```
/opt/crypto-sensor/config.yaml
```

### **Key Settings**
```yaml
sensor:
  name: "sensor-dc01"
  platform: "linux"
  version: "1.0.0"
  profile: "datacenter_host"

control_plane:
  url: "https://crypto-inventory.company.com"
  registration_key: "REG-tenant-123-20241215-A7B3C9"

network:
  interfaces: ["eth0", "eth1"]
  active_probing: true
  network_discovery: true

storage:
  data_path: "/opt/crypto-sensor/data"
  max_size: 100MB
  rotation_size: 10MB
  retention_days: 7

security:
  use_tls: true
  client_cert: "/opt/crypto-sensor/certs/client.crt"
  client_key: "/opt/crypto-sensor/certs/client.key"
  server_ca_cert: "/opt/crypto-sensor/certs/server-ca.crt"

features:
  tls_analysis: true
  ssh_analysis: true
  certificate_analysis: true
  active_probing: true
  network_discovery: true
  air_gapped_export: false
```

## 🚀 **Service Management**

### **Systemd Service**
```bash
# Check status
sudo systemctl status crypto-sensor

# View logs
sudo journalctl -u crypto-sensor -f

# Restart service
sudo systemctl restart crypto-sensor

# Stop service
sudo systemctl stop crypto-sensor

# Start service
sudo systemctl start crypto-sensor
```

### **Docker Deployment**
```yaml
version: '3.8'
services:
  crypto-sensor:
    image: crypto-inventory/sensor:latest
    container_name: crypto-sensor
    restart: unless-stopped
    environment:
      - REGISTRATION_KEY=REG-tenant-123-20241215-A7B3C9
      - CONTROL_PLANE_URL=https://crypto-inventory.company.com
      - SENSOR_NAME=sensor-dc01
      - SENSOR_IP=192.168.1.100
      - SENSOR_PROFILE=datacenter_host
    volumes:
      - /opt/crypto-sensor/data:/app/data
      - /opt/crypto-sensor/logs:/app/logs
    network_mode: host
    cap_add:
      - NET_RAW
    privileged: true
```

### **Kubernetes Deployment**
```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: crypto-sensor
spec:
  selector:
    matchLabels:
      app: crypto-sensor
  template:
    metadata:
      labels:
        app: crypto-sensor
    spec:
      containers:
      - name: crypto-sensor
        image: crypto-inventory/sensor:latest
        env:
        - name: REGISTRATION_KEY
          value: "REG-tenant-123-20241215-A7B3C9"
        - name: CONTROL_PLANE_URL
          value: "https://crypto-inventory.company.com"
        - name: SENSOR_NAME
          value: "sensor-dc01"
        - name: SENSOR_IP
          value: "192.168.1.100"
        - name: SENSOR_PROFILE
          value: "datacenter_host"
        volumeMounts:
        - name: data
          mountPath: /app/data
        - name: logs
          mountPath: /app/logs
        securityContext:
          capabilities:
            add:
            - NET_RAW
          privileged: true
      volumes:
      - name: data
        hostPath:
          path: /opt/crypto-sensor/data
      - name: logs
        hostPath:
          path: /opt/crypto-sensor/logs
```

## 🔍 **Monitoring and Troubleshooting**

### **Health Checks**
```bash
# Check agent health via API
curl https://crypto-inventory.company.com/api/v1/sensors/sensor-dc01/health

# View recent discoveries
curl https://crypto-inventory.company.com/api/v1/sensors/sensor-dc01/discoveries

# Check agent logs
sudo journalctl -u crypto-sensor -f
```

### **Common Issues**

#### **Registration Failed**
```bash
# Check registration key validity
curl -X POST https://crypto-inventory.company.com/api/v1/sensors/register \
  -H "Content-Type: application/json" \
  -d '{"registration_key": "REG-tenant-123-20241215-A7B3C9", ...}'

# Check network connectivity
ping crypto-inventory.company.com
telnet crypto-inventory.company.com 443
```

#### **IP Validation Failed**
```bash
# Check available IP addresses
ip addr show

# Verify IP is assigned to an interface
ip addr show | grep 192.168.1.100
```

#### **Permission Issues**
```bash
# Check capabilities
sudo setcap cap_net_raw+ep /opt/crypto-sensor/crypto-sensor

# Verify service user
id crypto-sensor
```

### **Log Analysis**
```bash
# View all logs
sudo journalctl -u crypto-sensor

# Filter error logs
sudo journalctl -u crypto-sensor --priority=err

# View specific time range
sudo journalctl -u crypto-sensor --since "2024-12-15 10:00:00"

# Follow logs in real-time
sudo journalctl -u crypto-sensor -f
```

## 📋 **Uninstallation**

### **Complete Removal**
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

### **Docker Removal**
```bash
# Stop and remove container
docker stop crypto-sensor
docker rm crypto-sensor

# Remove image (optional)
docker rmi crypto-inventory/sensor:latest
```

### **Kubernetes Removal**
```bash
# Delete DaemonSet
kubectl delete daemonset crypto-sensor

# Remove persistent volumes (optional)
kubectl delete pv crypto-sensor-data
kubectl delete pv crypto-sensor-logs
```

## 🔧 **Advanced Configuration**

### **Custom Network Interfaces**
```bash
# Specify custom interfaces
--interfaces "eth0,eth1,ens3"

# Monitor specific interface types
--interfaces "eth*"  # All ethernet interfaces
--interfaces "wlan*" # All wireless interfaces
```

### **Custom Installation Directory**
```bash
# Install to custom location
--dir /custom/path/crypto-sensor
```

### **Verbose Logging**
```bash
# Enable verbose output
--verbose
```

### **Environment Variables**
```bash
# Set via environment
export REGISTRATION_KEY="REG-tenant-123-20241215-A7B3C9"
export SENSOR_IP="192.168.1.100"
export SENSOR_NAME="sensor-dc01"
export SENSOR_PROFILE="datacenter_host"
```

## 📚 **Additional Resources**

- [Sensor Management Guide](./SENSOR_MANAGEMENT_GUIDE.md)
- [Security Architecture](./SECURITY_ARCHITECTURE.md)
- [API Documentation](./architecture_docs/05_api_specifications.md)
- [Troubleshooting Guide](./docs/troubleshooting.md)

---

**For support and questions, contact the platform team or check the documentation.**
