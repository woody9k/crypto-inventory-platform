# 📊 Reporting System Guide

## Overview

The Crypto Inventory Platform includes a comprehensive reporting system that generates various compliance and security reports based on the collected cryptographic inventory data. The system provides both real-time data access and scheduled report generation capabilities.

## 🏗️ Architecture

### Components

1. **Report Generator Service** (`services/report-generator/`)
   - Go-based microservice for report generation
   - RESTful API for report management
   - Asynchronous report processing
   - Multiple output formats (PDF, Excel, JSON)

2. **Reports UI** (`web-ui/src/pages/ReportsPage.tsx`)
   - React-based user interface
   - Report template selection
   - Real-time status tracking
   - Download and management capabilities

3. **Report Templates**
   - Predefined report structures
   - Configurable parameters
   - Category-based organization

## 📋 Available Report Types

### 1. Crypto Summary Report
**Purpose**: Overview of all cryptographic implementations across the network

**Key Metrics**:
- Total implementations: 1,247
- TLS connections: 892
- SSH servers: 156
- Certificates: 234
- Weak algorithms: 23
- Expired certificates: 12

**Data Includes**:
- Protocol distribution (TLS 1.3, TLS 1.2, SSH 2.0, etc.)
- Algorithm usage (AES-256-GCM, RSA-4096, ECDSA, etc.)
- Risk level analysis (Critical: 5, High: 18, Medium: 45, etc.)
- Implementation trends

### 2. Compliance Status Report
**Purpose**: Current compliance status against various frameworks

**Key Metrics**:
- Overall compliance score: 78%
- PCI-DSS compliance: 85%
- NIST Cybersecurity Framework: 72%
- FIPS 140-2 Level 2: 90%

**Data Includes**:
- Framework-specific scores
- Requirements met vs. total
- Critical and high issues count
- Actionable recommendations

### 3. Network Topology Report
**Purpose**: Network topology and sensor coverage analysis

**Key Metrics**:
- Total networks: 12
- Monitored networks: 10
- Coverage percentage: 83.3%
- Active sensors: 12

**Data Includes**:
- Sensor locations and status
- Network segment coverage
- Discovery statistics by network
- Sensor performance metrics

### 4. Risk Assessment Report
**Purpose**: Security risk assessment and recommendations

**Key Metrics**:
- Overall risk score: 6.2 (Medium)
- Critical findings: 2
- High-risk findings: 3
- Recommendations: 4

**Data Includes**:
- Risk level distribution
- Critical findings with affected assets
- Prioritized recommendations
- Effort vs. impact analysis

### 5. Certificate Audit Report
**Purpose**: SSL/TLS certificate inventory and expiration analysis

**Data Includes**:
- Certificate inventory
- Expiration tracking
- Weak certificate identification
- Renewal recommendations

## 🚀 API Endpoints

### Report Generation
```http
POST /api/v1/reports/generate
Content-Type: application/json

{
  "type": "crypto_summary",
  "title": "Monthly Crypto Summary",
  "parameters": {
    "date_range": "2024-12-01 to 2024-12-31",
    "include_trends": true
  },
  "format": "pdf"
}
```

### Report Management
```http
GET /api/v1/reports                    # List all reports
GET /api/v1/reports/{id}              # Get specific report
DELETE /api/v1/reports/{id}           # Delete report
GET /api/v1/reports/templates         # Get available templates
```

### Demo Data Access
```http
GET /api/v1/reports/demo/crypto-summary     # Crypto summary data
GET /api/v1/reports/demo/compliance-status  # Compliance status data
GET /api/v1/reports/demo/network-topology   # Network topology data
```

## 🎯 Usage Guide

### Generating Reports

1. **Via Web UI**:
   - Navigate to `/reports` in the web interface
   - Click "Generate Report" button
   - Select report template from the modal
   - Configure parameters as needed
   - Monitor generation status in real-time

2. **Via API**:
   ```bash
   curl -X POST http://localhost:8083/api/v1/reports/generate \
     -H "Content-Type: application/json" \
     -d '{
       "type": "crypto_summary",
       "title": "Custom Report",
       "format": "pdf"
     }'
   ```

### Report Status Tracking

Reports progress through the following states:
- **Generating**: Report is being created
- **Completed**: Report is ready for download
- **Failed**: Report generation encountered an error

### Downloading Reports

Completed reports can be downloaded via:
- Web UI download button
- Direct API access to download URL
- Programmatic access using report ID

## 🔧 Configuration

### Environment Variables

```bash
# Report Generator Service
PORT=8083                    # Service port
ENV=development              # Environment mode
DATABASE_URL=postgres://...  # Database connection
INFLUXDB_URL=http://...      # InfluxDB connection
NATS_URL=nats://...          # Message queue connection
LOG_LEVEL=debug              # Logging level
```

### Docker Configuration

The report generator is included in the Docker Compose setup:

```yaml
report-generator:
  build:
    context: ./services/report-generator
    dockerfile: Dockerfile.dev
  ports:
    - "8083:8083"
  environment:
    - PORT=8083
    - ENV=development
  depends_on:
    - postgres
    - redis
    - influxdb
    - nats
```

## 📊 Demo Data

The system includes comprehensive demo data for demonstration purposes:

- **1,247 Network Assets**: Realistic enterprise-scale inventory
- **892 Crypto Implementations**: Various protocols and algorithms
- **78% Compliance Score**: Mixed compliance status
- **12 Active Sensors**: Distributed network monitoring
- **Multiple Report Types**: All report categories represented

## 🔒 Security Considerations

### Data Protection
- Reports contain sensitive cryptographic information
- Access control through authentication system
- Secure download URLs with expiration
- Audit logging for report access

### Production Deployment
- Replace in-memory storage with database
- Implement proper authentication/authorization
- Add rate limiting for report generation
- Set up monitoring and alerting

## 🚀 Quick Start

1. **Start the Platform**:
   ```bash
   docker-compose up -d
   ```

2. **Access Reports**:
   - Web UI: http://localhost:3000/reports
   - API: http://localhost:8083/api/v1/reports

3. **Generate Demo Report**:
   ```bash
   curl http://localhost:8083/api/v1/reports/demo/crypto-summary
   ```

## 📈 Future Enhancements

### Planned Features
- **Scheduled Reports**: Automated report generation
- **Custom Templates**: User-defined report structures
- **Advanced Filtering**: Complex data filtering options
- **Export Formats**: Additional output formats (CSV, XML)
- **Email Delivery**: Automated report distribution
- **Report Scheduling**: Cron-based report generation
- **Advanced Analytics**: Machine learning insights
- **Compliance Automation**: Automated compliance checking

### Integration Points
- **SIEM Integration**: Security information and event management
- **Ticketing Systems**: Automated issue creation
- **Compliance Tools**: Direct integration with compliance platforms
- **Notification Systems**: Alert and notification integration

## 🛠️ Development

### Adding New Report Types

1. **Define Data Structure**:
   ```go
   func (h *Handler) generateNewReportData() map[string]interface{} {
     return map[string]interface{}{
       "metric1": value1,
       "metric2": value2,
     }
   }
   ```

2. **Add Template**:
   ```go
   {
     ID:          "new_report_type",
     Name:        "New Report Type",
     Description: "Description of the new report",
     Type:        "custom",
     Category:    "security",
   }
   ```

3. **Update UI**:
   - Add icon mapping in `getTypeIcon()`
   - Update template list
   - Add any specific UI components

### Testing

```bash
# Test report generator
cd services/report-generator
go test ./...

# Test API endpoints
curl http://localhost:8083/health
curl http://localhost:8083/api/v1/reports/templates
```

## 📞 Support

For questions or issues with the reporting system:

1. Check the API health endpoint: `GET /health`
2. Review service logs: `docker-compose logs report-generator`
3. Verify database connectivity
4. Check network connectivity between services

---

**The reporting system provides comprehensive insights into your cryptographic inventory, enabling informed security decisions and compliance management.**
