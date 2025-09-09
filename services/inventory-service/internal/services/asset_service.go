// Package services provides business logic for the crypto inventory management system.
// It handles asset discovery, crypto implementation tracking, risk assessment,
// and compliance analysis for cryptographic assets across network infrastructure.
//
// Key Features:
// - Asset CRUD operations with multi-tenant isolation
// - Advanced search and filtering capabilities
// - Risk scoring based on crypto implementation strength
// - Protocol and cipher suite analysis
// - Compliance framework integration
package services

import (
	"fmt"
	"inventory-service/internal/database"
	"inventory-service/internal/models"
	"strings"

	"github.com/google/uuid"
)

type AssetService struct {
	db *database.DB
}

func NewAssetService(db *database.DB) *AssetService {
	return &AssetService{db: db}
}

// GetAssets retrieves assets with filtering, pagination, and risk analysis.
// This is the core method for asset discovery and management with the following features:
// - Multi-tenant data isolation (only returns assets for the specified tenant)
// - Advanced filtering by asset type, environment, protocol, and risk level
// - Full-text search across hostname, IP address, and description
// - Risk scoring based on crypto implementation strength
// - Pagination support for large datasets
// - Sorting by risk score, discovery date, or custom fields
func (s *AssetService) GetAssets(tenantID uuid.UUID, filters models.AssetFilters) ([]models.Asset, int, error) {
	// Build the base query with tenant isolation and risk scoring
	baseQuery := `
		SELECT 
			a.id, a.tenant_id, a.hostname, a.ip_address, a.port, a.asset_type,
			a.operating_system, a.environment, a.business_unit, a.owner_email,
			a.description, a.tags, a.metadata, a.first_discovered_at, a.last_seen_at,
			a.created_at, a.updated_at, a.deleted_at,
			COALESCE(MAX(ci.risk_score), 0) as highest_risk
		FROM network_assets a
		LEFT JOIN crypto_implementations ci ON a.id = ci.asset_id AND ci.deleted_at IS NULL
		WHERE a.tenant_id = $1 AND a.deleted_at IS NULL
	`

	countQuery := `
		SELECT COUNT(DISTINCT a.id)
		FROM network_assets a
		LEFT JOIN crypto_implementations ci ON a.id = ci.asset_id AND ci.deleted_at IS NULL
		WHERE a.tenant_id = $1 AND a.deleted_at IS NULL
	`

	args := []interface{}{tenantID}
	argCount := 1
	whereConditions := []string{}

	// Add search filter
	if filters.Search != "" {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf(`(
			a.hostname ILIKE $%d OR 
			a.ip_address::text ILIKE $%d OR 
			a.description ILIKE $%d OR
			a.business_unit ILIKE $%d
		)`, argCount, argCount, argCount, argCount))
		searchPattern := "%" + filters.Search + "%"
		args = append(args, searchPattern)
	}

	// Add asset type filter
	if len(filters.AssetType) > 0 {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf(`a.asset_type = ANY($%d)`, argCount))
		args = append(args, filters.AssetType)
	}

	// Add environment filter
	if len(filters.Environment) > 0 {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf(`a.environment = ANY($%d)`, argCount))
		args = append(args, filters.Environment)
	}

	// Add protocol filter
	if len(filters.Protocol) > 0 {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf(`ci.protocol = ANY($%d)`, argCount))
		args = append(args, filters.Protocol)
	}

	// Add business unit filter
	if len(filters.BusinessUnit) > 0 {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf(`a.business_unit = ANY($%d)`, argCount))
		args = append(args, filters.BusinessUnit)
	}

	// Apply where conditions
	if len(whereConditions) > 0 {
		whereClause := " AND " + strings.Join(whereConditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
	}

	// Get total count
	var total int
	err := s.db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get assets count: %w", err)
	}

	// Add GROUP BY and sorting
	baseQuery += " GROUP BY a.id, a.tenant_id, a.hostname, a.ip_address, a.port, a.asset_type, a.operating_system, a.environment, a.business_unit, a.owner_email, a.description, a.tags, a.metadata, a.first_discovered_at, a.last_seen_at, a.created_at, a.updated_at, a.deleted_at"

	// Add sorting
	sortBy := "a.hostname"
	if filters.SortBy != "" {
		switch filters.SortBy {
		case "hostname", "ip_address", "asset_type", "environment", "created_at", "last_seen_at":
			sortBy = "a." + filters.SortBy
		case "risk_score":
			sortBy = "highest_risk"
		}
	}

	sortOrder := "ASC"
	if filters.SortOrder == "desc" {
		sortOrder = "DESC"
	}

	baseQuery += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	// Add pagination
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 20
	}

	offset := (filters.Page - 1) * filters.PageSize
	baseQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", filters.PageSize, offset)

	// Execute query
	rows, err := s.db.Query(baseQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query assets: %w", err)
	}
	defer rows.Close()

	var assets []models.Asset
	for rows.Next() {
		var asset models.Asset
		var highestRisk *int

		err := rows.Scan(
			&asset.ID, &asset.TenantID, &asset.Hostname, &asset.IPAddress, &asset.Port,
			&asset.AssetType, &asset.OperatingSystem, &asset.Environment, &asset.BusinessUnit,
			&asset.OwnerEmail, &asset.Description, &asset.Tags, &asset.Metadata,
			&asset.FirstDiscoveredAt, &asset.LastSeenAt, &asset.CreatedAt, &asset.UpdatedAt,
			&asset.DeletedAt, &highestRisk,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan asset: %w", err)
		}

		// Set risk information
		if highestRisk != nil {
			asset.RiskScore = *highestRisk
			asset.HighestRisk = highestRisk
		}
		asset.RiskLevel = models.GetRiskLevel(asset.RiskScore)

		assets = append(assets, asset)
	}

	// Apply risk level filter (post-query filtering since it's calculated)
	if len(filters.RiskLevel) > 0 {
		filteredAssets := []models.Asset{}
		for _, asset := range assets {
			for _, riskLevel := range filters.RiskLevel {
				if asset.RiskLevel == riskLevel {
					filteredAssets = append(filteredAssets, asset)
					break
				}
			}
		}
		assets = filteredAssets
	}

	return assets, total, nil
}

// GetAssetByID retrieves a single asset with its crypto implementations
func (s *AssetService) GetAssetByID(tenantID, assetID uuid.UUID) (*models.Asset, error) {
	// Get the asset
	query := `
		SELECT 
			id, tenant_id, hostname, ip_address, port, asset_type,
			operating_system, environment, business_unit, owner_email,
			description, tags, metadata, first_discovered_at, last_seen_at,
			created_at, updated_at, deleted_at
		FROM network_assets 
		WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`

	var asset models.Asset
	err := s.db.Get(&asset, query, assetID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}

	// Get crypto implementations
	cryptoImpls, err := s.GetCryptoImplementations(tenantID, assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get crypto implementations: %w", err)
	}

	asset.CryptoImplementations = cryptoImpls
	asset.CalculateAssetRiskScore()

	return &asset, nil
}

// GetCryptoImplementations retrieves crypto implementations for an asset
func (s *AssetService) GetCryptoImplementations(tenantID, assetID uuid.UUID) ([]models.CryptoImplementation, error) {
	query := `
		SELECT 
			id, tenant_id, asset_id, protocol, protocol_version, cipher_suite,
			key_exchange_algorithm, signature_algorithm, symmetric_encryption,
			hash_algorithm, key_size, certificate_id, discovery_method,
			confidence_score, source_sensor_id, raw_data, risk_score,
			compliance_status, first_discovered_at, last_verified_at,
			created_at, updated_at, deleted_at
		FROM crypto_implementations 
		WHERE asset_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		ORDER BY risk_score DESC, created_at DESC
	`

	var cryptoImpls []models.CryptoImplementation
	err := s.db.Select(&cryptoImpls, query, assetID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get crypto implementations: %w", err)
	}

	// Add risk analysis for each implementation
	for i := range cryptoImpls {
		cryptoImpls[i].RiskLevel = models.GetRiskLevel(cryptoImpls[i].RiskScore)
		cryptoImpls[i].RiskFactors = s.AnalyzeCryptoRisk(&cryptoImpls[i])
	}

	return cryptoImpls, nil
}

// GetRiskSummary calculates risk statistics for the tenant
func (s *AssetService) GetRiskSummary(tenantID uuid.UUID) (*models.RiskSummary, error) {
	query := `
		SELECT 
			COUNT(DISTINCT a.id) as total_assets,
			COUNT(DISTINCT ci.id) as total_crypto,
			COUNT(DISTINCT CASE WHEN ci.risk_score >= 70 THEN a.id END) as high_risk,
			COUNT(DISTINCT CASE WHEN ci.risk_score >= 40 AND ci.risk_score < 70 THEN a.id END) as medium_risk,
			COUNT(DISTINCT CASE WHEN ci.risk_score >= 1 AND ci.risk_score < 40 THEN a.id END) as low_risk,
			COUNT(DISTINCT CASE WHEN ci.risk_score = 0 OR ci.risk_score IS NULL THEN a.id END) as unknown_risk,
			COUNT(CASE WHEN ci.risk_score >= 80 THEN 1 END) as critical_findings
		FROM network_assets a
		LEFT JOIN crypto_implementations ci ON a.id = ci.asset_id AND ci.deleted_at IS NULL
		WHERE a.tenant_id = $1 AND a.deleted_at IS NULL
	`

	var summary models.RiskSummary
	err := s.db.Get(&summary, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get risk summary: %w", err)
	}

	return &summary, nil
}

// AnalyzeCryptoRisk analyzes crypto implementation and returns risk factors
func (s *AssetService) AnalyzeCryptoRisk(crypto *models.CryptoImplementation) []string {
	var riskFactors []string

	// Analyze cipher suite
	if crypto.CipherSuite != nil {
		cipherSuite := strings.ToUpper(*crypto.CipherSuite)
		if strings.Contains(cipherSuite, "RC4") ||
			strings.Contains(cipherSuite, "DES") ||
			strings.Contains(cipherSuite, "MD5") {
			riskFactors = append(riskFactors, "Weak cipher suite")
		}
	}

	// Analyze protocol version
	if crypto.ProtocolVersion != nil {
		version := *crypto.ProtocolVersion
		switch crypto.Protocol {
		case "TLS":
			if version < "1.2" {
				riskFactors = append(riskFactors, "Outdated TLS version")
			}
		case "SSH":
			if version < "2.0" {
				riskFactors = append(riskFactors, "Outdated SSH version")
			}
		}
	}

	// Analyze key size
	if crypto.KeySize != nil {
		keySize := *crypto.KeySize
		if keySize < 2048 {
			riskFactors = append(riskFactors, "Weak key size")
		}
	}

	// Analyze hash algorithm
	if crypto.HashAlgorithm != nil {
		hash := strings.ToUpper(*crypto.HashAlgorithm)
		if strings.Contains(hash, "MD5") || strings.Contains(hash, "SHA1") {
			riskFactors = append(riskFactors, "Weak hash algorithm")
		}
	}

	// Check confidence score
	if crypto.ConfidenceScore < 0.7 {
		riskFactors = append(riskFactors, "Low confidence detection")
	}

	return riskFactors
}
