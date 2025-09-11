# Reports System Documentation

## Overview

The Crypto Inventory Management System includes a comprehensive reports system that allows users to generate, view, and download various types of reports in multiple formats. The system is designed to provide detailed insights into cryptographic implementations, compliance status, network topology, and security assessments.

## Features

### Report Types

1. **Crypto Summary Report**
   - Overview of all cryptographic implementations across the network
   - Protocol breakdown (TLS, SSH versions)
   - Algorithm analysis (AES, RSA, ECDSA, etc.)
   - Risk level distribution
   - Implementation trends

2. **Compliance Status Report**
   - Current compliance status against various frameworks
   - Framework-specific scores (PCI-DSS, NIST, FIPS 140-2)
   - Requirements met vs. total requirements
   - Critical and high-priority issues
   - Actionable recommendations

3. **Network Topology Report**
   - Network coverage and sensor deployment status
   - Sensor details and discovery statistics
   - Network segmentation analysis
   - Coverage percentage metrics

4. **Risk Assessment Report**
   - Overall risk scoring and level assessment
   - Critical findings and security issues
   - Prioritized recommendations
   - Effort vs. impact analysis

5. **Certificate Audit Report**
   - SSL/TLS certificate inventory
   - Expiration analysis and alerts
   - Certificate authority distribution
   - Weak certificate identification

### Download Formats

All reports can be downloaded in three formats:

1. **PDF Format**
   - Formatted text reports with proper headers
   - Structured sections and subsections
   - Professional layout suitable for presentations
   - MIME type: `application/pdf`

2. **Excel Format (CSV)**
   - Structured data tables
   - Suitable for further analysis in spreadsheet applications
   - MIME type: `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`

3. **JSON Format**
   - Complete report data in structured JSON
   - Suitable for API integration and programmatic processing
   - MIME type: `application/json`

## API Endpoints

### Report Generation
```http
POST /api/v1/reports/generate
Content-Type: application/json

{
  "type": "crypto_summary",
  "title": "Monthly Crypto Report",
  "format": "pdf"
}
```

### Report Management
```http
GET    /api/v1/reports                    # List all reports
GET    /api/v1/reports/{id}               # Get specific report
DELETE /api/v1/reports/{id}               # Delete report
GET    /api/v1/reports/templates          # Get available templates
```

### Report Downloads
```http
GET /api/v1/reports/{id}/download?format=pdf     # Download as PDF
GET /api/v1/reports/{id}/download?format=excel   # Download as Excel
GET /api/v1/reports/{id}/download?format=json    # Download as JSON
```

## Frontend Components

### ReportsPage Component

The main reports interface (`/web-ui/src/pages/ReportsPage.tsx`) provides:

- **Report List**: Displays all generated reports with status indicators
- **Generate Button**: Opens modal to select report type and generate new reports
- **Download Buttons**: Individual buttons for PDF, Excel, and JSON formats
- **View Button**: Opens interactive report viewer modal
- **Delete Button**: Removes reports from the system
- **Real-time Updates**: Automatic polling for status changes

### Report Viewer Modal

The interactive report viewer provides:

- **Header Information**: Report type, status, and generation date
- **Data Visualization**: 
  - Summary metrics in card format
  - Protocol breakdown grids
  - Algorithm analysis tables
  - Risk level distributions
- **Download Actions**: All three formats available within the viewer
- **Responsive Design**: Works on desktop and mobile devices

## Backend Implementation

### Report Generator Service

Located in `/services/report-generator/`, the service includes:

- **Handlers**: HTTP request handlers for all report operations
- **Data Generation**: Realistic sample data for demonstration
- **Format Conversion**: Functions to convert data to PDF, Excel, and JSON
- **Async Processing**: Background report generation with status tracking

### Key Files

- `cmd/main.go`: Service entry point and route registration
- `internal/handlers/reports.go`: Core report generation and management logic
- `internal/handlers/data.go`: Sample data generation functions
- `internal/models/reports.go`: Data structures and types

### Data Formatting Functions

The system includes specialized formatting functions for each report type:

- `formatCryptoSummaryAsText()`: Formats crypto data for PDF output
- `formatCryptoSummaryAsCSV()`: Formats crypto data for Excel output
- `formatComplianceStatusAsText()`: Formats compliance data for PDF
- `formatComplianceStatusAsCSV()`: Formats compliance data for Excel
- And similar functions for all report types

## Usage Examples

### Generate a Report
```javascript
// Frontend example
const response = await fetch('/api/v1/reports/generate', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    type: 'crypto_summary',
    title: 'Q4 Security Report',
    format: 'pdf'
  })
});
```

### Download a Report
```javascript
// Frontend example
const downloadReport = async (reportId, format) => {
  const response = await fetch(`/api/v1/reports/${reportId}/download?format=${format}`);
  const blob = await response.blob();
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `report-${reportId}.${format === 'excel' ? 'xlsx' : format}`;
  a.click();
};
```

### Test with cURL
```bash
# Generate a report
curl -X POST http://localhost:8080/api/v1/reports/generate \
  -H "Content-Type: application/json" \
  -d '{"type": "crypto_summary", "title": "Test Report", "format": "pdf"}'

# Download as PDF
curl "http://localhost:8080/api/v1/reports/{report-id}/download?format=pdf" \
  --output report.pdf

# Download as Excel
curl "http://localhost:8080/api/v1/reports/{report-id}/download?format=excel" \
  --output report.xlsx

# Download as JSON
curl "http://localhost:8080/api/v1/reports/{report-id}/download?format=json" \
  --output report.json
```

## Configuration

### Service Configuration

The report generator service runs on port 8083 and is accessible through the API gateway on port 8080.

### Nginx Configuration

The API gateway includes proper routing for report endpoints:

```nginx
location /api/v1/reports {
    limit_req zone=api burst=10 nodelay;
    proxy_pass http://report_service/api/v1/reports;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_connect_timeout 5s;
    proxy_send_timeout 300s;  # Reports may take longer
    proxy_read_timeout 300s;
}
```

## Troubleshooting

### Common Issues

1. **Download Not Working**
   - Check if report status is "completed"
   - Verify API gateway routing
   - Check report service logs

2. **Report Generation Fails**
   - Ensure report service is running
   - Check database connectivity
   - Review service logs for errors

3. **Frontend Not Loading Reports**
   - Verify API gateway is accessible
   - Check CORS configuration
   - Ensure all services are healthy

### Debug Commands

```bash
# Check service status
docker-compose ps report-generator

# View service logs
docker-compose logs report-generator

# Test API endpoints
curl http://localhost:8080/api/v1/reports/templates
curl http://localhost:8080/api/v1/reports/

# Check API gateway logs
docker-compose logs api-gateway
```

## Future Enhancements

- **Real PDF Generation**: Integration with libraries like wkhtmltopdf
- **Advanced Excel Formatting**: Using libraries like excelize for better formatting
- **Report Scheduling**: Automated report generation on schedules
- **Email Integration**: Send reports via email
- **Custom Templates**: User-defined report templates
- **Data Export**: Export raw data for external analysis
- **Report Sharing**: Share reports with team members
- **Version History**: Track report generation history
