// Package handlers provides HTTP handlers for the report generator service.
// This file contains the core report generation and management functionality.
package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GenerateReport handles report generation requests from the web UI.
// It accepts report parameters, creates a new report record, and starts
// asynchronous report generation. Returns a 202 Accepted status with
// the report ID for tracking.
func (h *Handler) GenerateReport(c *gin.Context) {
	var req struct {
		Type       string                 `json:"type" binding:"required"` // Report type (crypto_summary, compliance_status, etc.)
		Title      string                 `json:"title"`                   // Optional custom title
		Parameters map[string]interface{} `json:"parameters"`              // Report-specific parameters
		Format     string                 `json:"format"`                  // Output format (pdf, excel, json)
	}

	// Parse and validate the request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Generate unique report ID for tracking
	reportID := uuid.New().String()

	// Set default title if not provided by user
	if req.Title == "" {
		req.Title = fmt.Sprintf("%s Report", req.Type)
	}

	// Set default output format to PDF if not specified
	if req.Format == "" {
		req.Format = "pdf"
	}

	// Create report record with initial status
	report := &Report{
		ID:        reportID,
		Title:     req.Title,
		Type:      req.Type,
		Status:    "generating",
		CreatedAt: time.Now(),
	}

	// Store report in memory (in production, this would be a database insert)
	h.reports[reportID] = report

	// Start asynchronous report generation
	// In production, this would queue the job in a message queue system
	go h.generateReportAsync(reportID, req.Type, req.Format)

	// Return 202 Accepted with report ID for tracking
	c.JSON(202, gin.H{
		"report_id": reportID,
		"status":    "generating",
		"message":   "Report generation started",
	})
}

// generateReportAsync simulates asynchronous report generation.
// In production, this would be handled by a background job processor
// that queries the database for crypto inventory data and generates
// the appropriate report content.
func (h *Handler) generateReportAsync(reportID, reportType, format string) {
	// Simulate processing time for realistic demo experience
	// In production, this would be the actual report generation time
	time.Sleep(2 * time.Second)

	// Retrieve the report record
	report := h.reports[reportID]
	if report == nil {
		return // Report not found, skip processing
	}

	// Generate report data based on the requested type
	// In production, this would query the actual database for real data
	var data interface{}
	switch reportType {
	case "crypto_summary":
		data = h.generateCryptoSummaryData()
	case "compliance_status":
		data = h.generateComplianceStatusData()
	case "network_topology":
		data = h.generateNetworkTopologyData()
	case "risk_assessment":
		data = h.generateRiskAssessmentData()
	default:
		// Fallback for unknown report types
		data = map[string]interface{}{
			"message": "Report data generated",
			"type":    reportType,
		}
	}

	// Update report status to completed
	now := time.Now()
	report.Status = "completed"
	report.CompletedAt = &now
	report.Data = data
	report.DownloadURL = fmt.Sprintf("/api/v1/reports/%s/download", reportID)

	// Store updated report (in production, this would be a database update)
	h.reports[reportID] = report
}

// GetReport retrieves a specific report by its ID.
// Returns the complete report data including status, content, and metadata.
func (h *Handler) GetReport(c *gin.Context) {
	reportID := c.Param("id")

	// Look up report in storage
	report, exists := h.reports[reportID]
	if !exists {
		c.JSON(404, gin.H{"error": "Report not found"})
		return
	}

	// Return the complete report data
	c.JSON(200, report)
}

// ListReports returns a list of all reports in the system.
// This endpoint is used by the web UI to display the reports dashboard.
func (h *Handler) ListReports(c *gin.Context) {
	var reports []*Report
	for _, report := range h.reports {
		reports = append(reports, report)
	}

	c.JSON(200, gin.H{"reports": reports})
}

// DeleteReport removes a report from the system.
// This permanently deletes the report and its associated data.
func (h *Handler) DeleteReport(c *gin.Context) {
	reportID := c.Param("id")

	// Check if report exists before attempting deletion
	if _, exists := h.reports[reportID]; !exists {
		c.JSON(404, gin.H{"error": "Report not found"})
		return
	}

	// Remove report from storage
	delete(h.reports, reportID)
	c.JSON(200, gin.H{"message": "Report deleted successfully"})
}

// DownloadReport handles report downloads in various formats (PDF, JSON, Excel).
// This endpoint serves the actual report files for download by the frontend.
func (h *Handler) DownloadReport(c *gin.Context) {
	reportID := c.Param("id")
	format := c.Query("format")

	// Default to JSON if no format specified
	if format == "" {
		format = "json"
	}

	// Look up report in storage
	report, exists := h.reports[reportID]
	if !exists {
		c.JSON(404, gin.H{"error": "Report not found"})
		return
	}

	// Check if report is completed
	if report.Status != "completed" {
		c.JSON(400, gin.H{"error": "Report not ready for download"})
		return
	}

	// Generate file based on requested format
	switch format {
	case "pdf":
		h.generatePDFReport(c, report)
	case "excel":
		h.generateExcelReport(c, report)
	case "json":
		h.generateJSONReport(c, report)
	default:
		c.JSON(400, gin.H{"error": "Unsupported format. Use: pdf, excel, or json"})
	}
}

// generatePDFReport creates a PDF version of the report.
// In production, this would use a proper PDF library like wkhtmltopdf or similar.
func (h *Handler) generatePDFReport(c *gin.Context, report *Report) {
	// For demo purposes, we'll return a simple text representation
	// In production, you'd use a proper PDF generation library
	content := h.formatReportAsText(report)

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.pdf\"", report.Title))
	c.String(200, content)
}

// generateExcelReport creates an Excel version of the report.
// In production, this would use a proper Excel library like excelize.
func (h *Handler) generateExcelReport(c *gin.Context, report *Report) {
	// For demo purposes, we'll return a CSV representation
	// In production, you'd use a proper Excel generation library
	content := h.formatReportAsCSV(report)

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.xlsx\"", report.Title))
	c.String(200, content)
}

// generateJSONReport creates a JSON version of the report.
func (h *Handler) generateJSONReport(c *gin.Context, report *Report) {
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.json\"", report.Title))
	c.JSON(200, report)
}

// formatReportAsText formats the report data as plain text for PDF generation.
func (h *Handler) formatReportAsText(report *Report) string {
	content := fmt.Sprintf("Report: %s\n", report.Title)
	content += fmt.Sprintf("Type: %s\n", report.Type)
	content += fmt.Sprintf("Generated: %s\n", report.CreatedAt.Format("2006-01-02 15:04:05"))
	content += "=" + strings.Repeat("=", 50) + "\n\n"

	// Format data based on report type
	switch report.Type {
	case "crypto_summary":
		content += h.formatCryptoSummaryAsText(report.Data)
	case "compliance_status":
		content += h.formatComplianceStatusAsText(report.Data)
	case "network_topology":
		content += h.formatNetworkTopologyAsText(report.Data)
	case "risk_assessment":
		content += h.formatRiskAssessmentAsText(report.Data)
	default:
		content += "Report data:\n"
		content += fmt.Sprintf("%+v\n", report.Data)
	}

	return content
}

// formatReportAsCSV formats the report data as CSV for Excel generation.
func (h *Handler) formatReportAsCSV(report *Report) string {
	var content strings.Builder
	content.WriteString(fmt.Sprintf("Report,Type,Generated\n"))
	content.WriteString(fmt.Sprintf("%s,%s,%s\n", report.Title, report.Type, report.CreatedAt.Format("2006-01-02 15:04:05")))
	content.WriteString("\n")

	// Format data based on report type
	switch report.Type {
	case "crypto_summary":
		content.WriteString(h.formatCryptoSummaryAsCSV(report.Data))
	case "compliance_status":
		content.WriteString(h.formatComplianceStatusAsCSV(report.Data))
	case "network_topology":
		content.WriteString(h.formatNetworkTopologyAsCSV(report.Data))
	case "risk_assessment":
		content.WriteString(h.formatRiskAssessmentAsCSV(report.Data))
	default:
		content.WriteString("Data\n")
		content.WriteString(fmt.Sprintf("%+v\n", report.Data))
	}

	return content.String()
}

// GetReportTemplates returns the list of available report templates.
// Templates define the structure and parameters for different report types
// that users can generate through the web UI.
func (h *Handler) GetReportTemplates(c *gin.Context) {
	templates := []ReportTemplate{
		{
			ID:          "crypto_summary",
			Name:        "Crypto Summary Report",
			Description: "Overview of all cryptographic implementations across the network",
			Type:        "summary",
			Category:    "crypto",
		},
		{
			ID:          "compliance_status",
			Name:        "Compliance Status Report",
			Description: "Current compliance status against various frameworks",
			Type:        "compliance",
			Category:    "compliance",
		},
		{
			ID:          "network_topology",
			Name:        "Network Topology Report",
			Description: "Network topology and sensor coverage map",
			Type:        "topology",
			Category:    "network",
		},
		{
			ID:          "risk_assessment",
			Name:        "Risk Assessment Report",
			Description: "Security risk assessment and recommendations",
			Type:        "risk",
			Category:    "security",
		},
		{
			ID:          "certificate_audit",
			Name:        "Certificate Audit Report",
			Description: "SSL/TLS certificate inventory and expiration analysis",
			Type:        "audit",
			Category:    "crypto",
		},
	}

	c.JSON(200, gin.H{"templates": templates})
}

// Demo endpoints for quick data access
// These endpoints provide immediate access to sample data for demonstration purposes.
// In production, these would be replaced with actual database queries.

// GetCryptoSummary returns crypto summary data for demo purposes.
// This endpoint provides sample data showing cryptographic implementations
// across the network without requiring report generation.
func (h *Handler) GetCryptoSummary(c *gin.Context) {
	data := h.generateCryptoSummaryData()
	c.JSON(200, data)
}

// GetComplianceStatus returns compliance status data for demo purposes.
// This endpoint provides sample compliance scores and framework status
// for immediate demonstration of the reporting capabilities.
func (h *Handler) GetComplianceStatus(c *gin.Context) {
	data := h.generateComplianceStatusData()
	c.JSON(200, data)
}

// GetNetworkTopology returns network topology data for demo purposes.
// This endpoint provides sample sensor and network coverage data
// for demonstrating the network monitoring capabilities.
func (h *Handler) GetNetworkTopology(c *gin.Context) {
	data := h.generateNetworkTopologyData()
	c.JSON(200, data)
}

// Data generation functions
// These functions generate realistic sample data for demonstration purposes.
// In production, these would query the actual database for real inventory data.

// generateCryptoSummaryData creates sample data for the crypto summary report.
// This includes protocol distributions, algorithm usage, risk levels, and trends
// that would typically be found in a real crypto inventory system.
func (h *Handler) generateCryptoSummaryData() map[string]interface{} {
	return map[string]interface{}{
		"summary": map[string]interface{}{
			"total_implementations": 1247,
			"tls_connections":       892,
			"ssh_servers":           156,
			"certificates":          234,
			"weak_algorithms":       23,
			"expired_certificates":  12,
		},
		"by_protocol": map[string]interface{}{
			"TLS 1.3":  456,
			"TLS 1.2":  321,
			"TLS 1.1":  89,
			"TLS 1.0":  26,
			"SSH 2.0":  156,
			"SSH 1.99": 2,
		},
		"by_algorithm": map[string]interface{}{
			"AES-256-GCM":       234,
			"AES-128-GCM":       189,
			"ChaCha20-Poly1305": 67,
			"RSA-4096":          123,
			"RSA-2048":          89,
			"ECDSA P-256":       156,
			"ECDSA P-384":       78,
		},
		"risk_levels": map[string]interface{}{
			"critical": 5,
			"high":     18,
			"medium":   45,
			"low":      89,
			"info":     1090,
		},
		"trends": map[string]interface{}{
			"new_implementations":     23,
			"updated_implementations": 45,
			"removed_implementations": 8,
		},
	}
}

// generateComplianceStatusData creates sample compliance status data.
// This includes overall compliance scores, framework-specific status,
// requirements met, and actionable recommendations.
func (h *Handler) generateComplianceStatusData() map[string]interface{} {
	return map[string]interface{}{
		"overall_score": 78,
		"frameworks": []map[string]interface{}{
			{
				"name":               "PCI-DSS",
				"version":            "4.0",
				"score":              85,
				"status":             "compliant",
				"requirements_met":   8,
				"requirements_total": 10,
				"critical_issues":    0,
				"high_issues":        2,
			},
			{
				"name":               "NIST Cybersecurity Framework",
				"version":            "1.1",
				"score":              72,
				"status":             "partially_compliant",
				"requirements_met":   15,
				"requirements_total": 20,
				"critical_issues":    1,
				"high_issues":        4,
			},
			{
				"name":               "FIPS 140-2",
				"version":            "Level 2",
				"score":              90,
				"status":             "compliant",
				"requirements_met":   18,
				"requirements_total": 20,
				"critical_issues":    0,
				"high_issues":        2,
			},
		},
		"recommendations": []string{
			"Upgrade TLS 1.0/1.1 implementations to TLS 1.2 or higher",
			"Replace weak RSA-1024 certificates with RSA-2048 or higher",
			"Implement certificate lifecycle management",
			"Enable HSTS headers on all HTTPS endpoints",
		},
	}
}

// generateNetworkTopologyData creates sample network topology data.
// This includes sensor information, network coverage, and discovery statistics
// that would be collected from deployed network sensors.
func (h *Handler) generateNetworkTopologyData() map[string]interface{} {
	return map[string]interface{}{
		"sensors": []map[string]interface{}{
			{
				"id":          "sensor-dc01",
				"name":        "Datacenter Core",
				"status":      "active",
				"location":    "Primary DC - Rack 1",
				"interfaces":  []string{"eth0", "eth1"},
				"discoveries": 456,
				"last_seen":   "2024-12-15T10:30:00Z",
			},
			{
				"id":          "sensor-dc02",
				"name":        "Datacenter Edge",
				"status":      "active",
				"location":    "Primary DC - Rack 2",
				"interfaces":  []string{"eth0"},
				"discoveries": 234,
				"last_seen":   "2024-12-15T10:29:45Z",
			},
			{
				"id":          "sensor-cloud01",
				"name":        "Cloud Instance",
				"status":      "active",
				"location":    "AWS us-east-1",
				"interfaces":  []string{"ens3"},
				"discoveries": 123,
				"last_seen":   "2024-12-15T10:28:30Z",
			},
		},
		"coverage": map[string]interface{}{
			"total_networks":      12,
			"monitored_networks":  10,
			"coverage_percentage": 83.3,
			"uncovered_networks":  []string{"192.168.100.0/24", "10.0.50.0/24"},
		},
		"discoveries_by_network": map[string]interface{}{
			"192.168.1.0/24": 456,
			"192.168.2.0/24": 234,
			"10.0.1.0/24":    123,
			"172.16.1.0/24":  89,
		},
	}
}

// generateRiskAssessmentData creates sample risk assessment data.
// This includes overall risk scores, critical findings, and prioritized
// recommendations based on security analysis of the crypto inventory.
func (h *Handler) generateRiskAssessmentData() map[string]interface{} {
	return map[string]interface{}{
		"overall_risk_score": 6.2,
		"risk_level":         "medium",
		"critical_findings": []map[string]interface{}{
			{
				"type":            "weak_cipher",
				"description":     "TLS 1.0 with RC4 cipher detected",
				"severity":        "critical",
				"count":           5,
				"affected_assets": []string{"web-server-01", "api-gateway-02"},
			},
			{
				"type":            "expired_certificate",
				"description":     "SSL certificate expired 30 days ago",
				"severity":        "high",
				"count":           3,
				"affected_assets": []string{"legacy-app-01", "internal-tool-02"},
			},
		},
		"recommendations": []map[string]interface{}{
			{
				"priority": "high",
				"action":   "Disable TLS 1.0 and RC4 cipher",
				"effort":   "low",
				"impact":   "high",
			},
			{
				"priority": "high",
				"action":   "Renew expired certificates",
				"effort":   "medium",
				"impact":   "high",
			},
		},
	}
}

// Formatting functions for different report types

// formatCryptoSummaryAsText formats crypto summary data as text
func (h *Handler) formatCryptoSummaryAsText(data interface{}) string {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return "Invalid data format\n"
	}

	content := "CRYPTO SUMMARY REPORT\n"
	content += strings.Repeat("-", 30) + "\n\n"

	// Summary section
	if summary, ok := dataMap["summary"].(map[string]interface{}); ok {
		content += "SUMMARY:\n"
		content += fmt.Sprintf("  Total Implementations: %v\n", summary["total_implementations"])
		content += fmt.Sprintf("  TLS Connections: %v\n", summary["tls_connections"])
		content += fmt.Sprintf("  SSH Servers: %v\n", summary["ssh_servers"])
		content += fmt.Sprintf("  Certificates: %v\n", summary["certificates"])
		content += fmt.Sprintf("  Weak Algorithms: %v\n", summary["weak_algorithms"])
		content += fmt.Sprintf("  Expired Certificates: %v\n", summary["expired_certificates"])
		content += "\n"
	}

	// Protocol breakdown
	if protocols, ok := dataMap["by_protocol"].(map[string]interface{}); ok {
		content += "PROTOCOL BREAKDOWN:\n"
		for protocol, count := range protocols {
			content += fmt.Sprintf("  %s: %v\n", protocol, count)
		}
		content += "\n"
	}

	// Algorithm breakdown
	if algorithms, ok := dataMap["by_algorithm"].(map[string]interface{}); ok {
		content += "ALGORITHM BREAKDOWN:\n"
		for algorithm, count := range algorithms {
			content += fmt.Sprintf("  %s: %v\n", algorithm, count)
		}
		content += "\n"
	}

	// Risk levels
	if risks, ok := dataMap["risk_levels"].(map[string]interface{}); ok {
		content += "RISK LEVELS:\n"
		for level, count := range risks {
			content += fmt.Sprintf("  %s: %v\n", strings.Title(level), count)
		}
	}

	return content
}

// formatCryptoSummaryAsCSV formats crypto summary data as CSV
func (h *Handler) formatCryptoSummaryAsCSV(data interface{}) string {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return "Invalid data format\n"
	}

	var content strings.Builder
	content.WriteString("Metric,Value\n")

	// Summary section
	if summary, ok := dataMap["summary"].(map[string]interface{}); ok {
		for key, value := range summary {
			content.WriteString(fmt.Sprintf("%s,%v\n", key, value))
		}
	}

	return content.String()
}

// formatComplianceStatusAsText formats compliance status data as text
func (h *Handler) formatComplianceStatusAsText(data interface{}) string {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return "Invalid data format\n"
	}

	content := "COMPLIANCE STATUS REPORT\n"
	content += strings.Repeat("-", 30) + "\n\n"

	// Overall score
	if score, ok := dataMap["overall_score"].(float64); ok {
		content += fmt.Sprintf("Overall Compliance Score: %.0f%%\n\n", score)
	}

	// Framework details
	if frameworks, ok := dataMap["frameworks"].([]interface{}); ok {
		content += "FRAMEWORK DETAILS:\n"
		for _, fw := range frameworks {
			if fwMap, ok := fw.(map[string]interface{}); ok {
				content += fmt.Sprintf("  %s %s: %.0f%% (%s)\n",
					fwMap["name"], fwMap["version"], fwMap["score"], fwMap["status"])
				content += fmt.Sprintf("    Requirements: %.0f/%.0f\n",
					fwMap["requirements_met"], fwMap["requirements_total"])
				content += fmt.Sprintf("    Critical Issues: %.0f, High Issues: %.0f\n",
					fwMap["critical_issues"], fwMap["high_issues"])
				content += "\n"
			}
		}
	}

	// Recommendations
	if recommendations, ok := dataMap["recommendations"].([]interface{}); ok {
		content += "RECOMMENDATIONS:\n"
		for i, rec := range recommendations {
			content += fmt.Sprintf("  %d. %s\n", i+1, rec)
		}
	}

	return content
}

// formatComplianceStatusAsCSV formats compliance status data as CSV
func (h *Handler) formatComplianceStatusAsCSV(data interface{}) string {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return "Invalid data format\n"
	}

	var content strings.Builder
	content.WriteString("Framework,Version,Score,Status,Requirements Met,Total Requirements,Critical Issues,High Issues\n")

	if frameworks, ok := dataMap["frameworks"].([]interface{}); ok {
		for _, fw := range frameworks {
			if fwMap, ok := fw.(map[string]interface{}); ok {
				content.WriteString(fmt.Sprintf("%s,%s,%.0f,%s,%.0f,%.0f,%.0f,%.0f\n",
					fwMap["name"], fwMap["version"], fwMap["score"], fwMap["status"],
					fwMap["requirements_met"], fwMap["requirements_total"],
					fwMap["critical_issues"], fwMap["high_issues"]))
			}
		}
	}

	return content.String()
}

// formatNetworkTopologyAsText formats network topology data as text
func (h *Handler) formatNetworkTopologyAsText(data interface{}) string {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return "Invalid data format\n"
	}

	content := "NETWORK TOPOLOGY REPORT\n"
	content += strings.Repeat("-", 30) + "\n\n"

	// Coverage information
	if coverage, ok := dataMap["coverage"].(map[string]interface{}); ok {
		content += "COVERAGE SUMMARY:\n"
		content += fmt.Sprintf("  Total Networks: %.0f\n", coverage["total_networks"])
		content += fmt.Sprintf("  Monitored Networks: %.0f\n", coverage["monitored_networks"])
		content += fmt.Sprintf("  Coverage Percentage: %.1f%%\n", coverage["coverage_percentage"])
		content += "\n"
	}

	// Sensor details
	if sensors, ok := dataMap["sensors"].([]interface{}); ok {
		content += "SENSOR DETAILS:\n"
		for _, sensor := range sensors {
			if sensorMap, ok := sensor.(map[string]interface{}); ok {
				content += fmt.Sprintf("  %s (%s)\n", sensorMap["name"], sensorMap["id"])
				content += fmt.Sprintf("    Status: %s\n", sensorMap["status"])
				content += fmt.Sprintf("    Location: %s\n", sensorMap["location"])
				content += fmt.Sprintf("    Discoveries: %.0f\n", sensorMap["discoveries"])
				content += fmt.Sprintf("    Last Seen: %s\n", sensorMap["last_seen"])
				content += "\n"
			}
		}
	}

	return content
}

// formatNetworkTopologyAsCSV formats network topology data as CSV
func (h *Handler) formatNetworkTopologyAsCSV(data interface{}) string {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return "Invalid data format\n"
	}

	var content strings.Builder
	content.WriteString("Sensor ID,Name,Status,Location,Discoveries,Last Seen\n")

	if sensors, ok := dataMap["sensors"].([]interface{}); ok {
		for _, sensor := range sensors {
			if sensorMap, ok := sensor.(map[string]interface{}); ok {
				content.WriteString(fmt.Sprintf("%s,%s,%s,%s,%.0f,%s\n",
					sensorMap["id"], sensorMap["name"], sensorMap["status"],
					sensorMap["location"], sensorMap["discoveries"], sensorMap["last_seen"]))
			}
		}
	}

	return content.String()
}

// formatRiskAssessmentAsText formats risk assessment data as text
func (h *Handler) formatRiskAssessmentAsText(data interface{}) string {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return "Invalid data format\n"
	}

	content := "RISK ASSESSMENT REPORT\n"
	content += strings.Repeat("-", 30) + "\n\n"

	// Overall risk score
	if score, ok := dataMap["overall_risk_score"].(float64); ok {
		content += fmt.Sprintf("Overall Risk Score: %.1f\n", score)
	}
	if level, ok := dataMap["risk_level"].(string); ok {
		content += fmt.Sprintf("Risk Level: %s\n\n", strings.Title(level))
	}

	// Critical findings
	if findings, ok := dataMap["critical_findings"].([]interface{}); ok {
		content += "CRITICAL FINDINGS:\n"
		for i, finding := range findings {
			if findingMap, ok := finding.(map[string]interface{}); ok {
				content += fmt.Sprintf("  %d. %s\n", i+1, findingMap["description"])
				content += fmt.Sprintf("     Type: %s\n", findingMap["type"])
				content += fmt.Sprintf("     Severity: %s\n", findingMap["severity"])
				content += fmt.Sprintf("     Count: %.0f\n", findingMap["count"])
				content += "\n"
			}
		}
	}

	// Recommendations
	if recommendations, ok := dataMap["recommendations"].([]interface{}); ok {
		content += "RECOMMENDATIONS:\n"
		for i, rec := range recommendations {
			if recMap, ok := rec.(map[string]interface{}); ok {
				content += fmt.Sprintf("  %d. %s\n", i+1, recMap["action"])
				content += fmt.Sprintf("     Priority: %s\n", recMap["priority"])
				content += fmt.Sprintf("     Effort: %s, Impact: %s\n", recMap["effort"], recMap["impact"])
				content += "\n"
			}
		}
	}

	return content
}

// formatRiskAssessmentAsCSV formats risk assessment data as CSV
func (h *Handler) formatRiskAssessmentAsCSV(data interface{}) string {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return "Invalid data format\n"
	}

	var content strings.Builder
	content.WriteString("Type,Description,Severity,Count,Priority,Action,Effort,Impact\n")

	// Critical findings
	if findings, ok := dataMap["critical_findings"].([]interface{}); ok {
		for _, finding := range findings {
			if findingMap, ok := finding.(map[string]interface{}); ok {
				content.WriteString(fmt.Sprintf("Finding,%s,%s,%.0f,,,\n",
					findingMap["description"], findingMap["severity"], findingMap["count"]))
			}
		}
	}

	// Recommendations
	if recommendations, ok := dataMap["recommendations"].([]interface{}); ok {
		for _, rec := range recommendations {
			if recMap, ok := rec.(map[string]interface{}); ok {
				content.WriteString(fmt.Sprintf("Recommendation,,,%s,%s,%s,%s\n",
					recMap["priority"], recMap["action"], recMap["effort"], recMap["impact"]))
			}
		}
	}

	return content.String()
}
