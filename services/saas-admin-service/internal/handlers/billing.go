package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"saas-admin-service/internal/billing"
	"saas-admin-service/internal/config"

	"github.com/gin-gonic/gin"
)

// GetTenantBilling returns normalized billing info for a tenant across providers
func GetTenantBilling(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.Param("id")

		var payload = gin.H{}

		// Basic customer reference (if exists)
		queryCustomer := `
			SELECT bc.external_customer_id, bp.key as provider_key, bp.display_name
			FROM billing_customers bc
			JOIN billing_providers bp ON bc.provider_id = bp.id
			WHERE bc.tenant_id = $1
		`

		rows, err := db.Query(queryCustomer, tenantID)
		if err == nil {
			defer rows.Close()
			customers := []gin.H{}
			for rows.Next() {
				var externalID, providerKey, providerName string
				if err := rows.Scan(&externalID, &providerKey, &providerName); err == nil {
					customers = append(customers, gin.H{
						"provider":      providerKey,
						"provider_name": providerName,
						"customer_id":   externalID,
					})
				}
			}
			payload["customers"] = customers
		}

		// Subscription snapshot
		querySub := `
			SELECT bs.external_subscription_id, bs.plan_key, bs.status, bs.current_period_start, bs.current_period_end,
			       bp.key as provider_key
			FROM billing_subscriptions bs
			JOIN billing_providers bp ON bs.provider_id = bp.id
			WHERE bs.tenant_id = $1
		`
		subRows, err := db.Query(querySub, tenantID)
		if err == nil {
			defer subRows.Close()
			subs := []gin.H{}
			for subRows.Next() {
				var extID, planKey, status, providerKey string
				var start, end sql.NullTime
				if err := subRows.Scan(&extID, &planKey, &status, &start, &end, &providerKey); err == nil {
					subs = append(subs, gin.H{
						"provider":             providerKey,
						"subscription_id":      extID,
						"plan":                 planKey,
						"status":               status,
						"current_period_start": start.Time,
						"current_period_end":   end.Time,
					})
				}
			}
			payload["subscriptions"] = subs
		}

		c.JSON(http.StatusOK, gin.H{"billing": payload})
	}
}

// UpdateTenantBilling allows plan changes or cancellation flags (normalized)
func UpdateTenantBilling(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.Param("id")
		var req billing.UpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Determine provider for tenant (simplified: default to stripe)
		registry := getBillingRegistry(c)
		prov, err := registry.Get("stripe")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Billing provider not configured"})
			return
		}

		if err := prov.ApplyUpdate(db, tenantID, req); err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		// Best-effort reflect intent in billing_subscriptions (provider-agnostic upsert)
		if req.PlanKey != nil || req.CancelAtPeriodEnd != nil {
			var providerID string
			_ = db.QueryRow("SELECT id FROM billing_providers WHERE key = $1", "stripe").Scan(&providerID)
			now := time.Now()
			if req.PlanKey != nil {
				_, _ = db.Exec(`
					INSERT INTO billing_subscriptions (tenant_id, provider_id, external_subscription_id, plan_key, status, current_period_start, current_period_end, cancel_at_period_end)
					VALUES ($1, $2, 'pending', $3, 'active', $4, $5, COALESCE($6, false))
					ON CONFLICT DO NOTHING
				`, tenantID, providerID, *req.PlanKey, now, now.Add(30*24*time.Hour), req.CancelAtPeriodEnd)
				_, _ = db.Exec(`
					UPDATE billing_subscriptions
					SET plan_key = $1, updated_at = NOW(), cancel_at_period_end = COALESCE($2, cancel_at_period_end)
					WHERE tenant_id = $3 AND provider_id = $4
				`, *req.PlanKey, req.CancelAtPeriodEnd, tenantID, providerID)
			} else if req.CancelAtPeriodEnd != nil {
				_, _ = db.Exec(`
					UPDATE billing_subscriptions
					SET cancel_at_period_end = $1, updated_at = NOW()
					WHERE tenant_id = $2 AND provider_id = $3
				`, *req.CancelAtPeriodEnd, tenantID, providerID)
			}
		}

		c.JSON(http.StatusAccepted, gin.H{"message": "Billing update scheduled"})
	}
}

// ListInvoices returns a list of invoices for a tenant or globally
func ListInvoices(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.Query("tenantId")

		var rows *sql.Rows
		var err error
		if tenantID != "" {
			rows, err = db.Query(
				"SELECT external_invoice_id, amount_cents, currency, status, issued_at, due_at, paid_at FROM billing_invoices WHERE tenant_id = $1 ORDER BY issued_at DESC",
				tenantID,
			)
		} else {
			rows, err = db.Query(
				"SELECT external_invoice_id, amount_cents, currency, status, issued_at, due_at, paid_at FROM billing_invoices ORDER BY issued_at DESC",
			)
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch invoices"})
			return
		}
		defer rows.Close()

		invoices := []gin.H{}
		for rows.Next() {
			var externalID, currency, status string
			var amount int
			var issuedAt, dueAt, paidAt sql.NullTime
			if err := rows.Scan(&externalID, &amount, &currency, &status, &issuedAt, &dueAt, &paidAt); err == nil {
				invoices = append(invoices, gin.H{
					"invoice_id":   externalID,
					"amount_cents": amount,
					"currency":     currency,
					"status":       status,
					"issued_at":    issuedAt.Time,
					"due_at":       dueAt.Time,
					"paid_at":      paidAt.Time,
				})
			}
		}

		c.JSON(http.StatusOK, gin.H{"invoices": invoices})
	}
}

// Webhook endpoint is provider-agnostic; verification done per provider key
func BillingWebhook(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		providerKey := c.Param("provider")
		body, _ := io.ReadAll(c.Request.Body)

		registry := getBillingRegistry(c)
		prov, err := registry.Get(providerKey)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Unknown provider"})
			return
		}

		if err := prov.VerifyWebhook(c.Request.Header, body); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
			return
		}

		// Persist minimal event record (optional: enqueue async processing)
		var providerID string
		if err := db.QueryRow("SELECT id FROM billing_providers WHERE key = $1", providerKey).Scan(&providerID); err == nil {
			var payload map[string]any
			_ = json.Unmarshal(body, &payload)
			eventType := "unknown"
			externalEventID := ""
			if t, ok := payload["type"].(string); ok {
				eventType = t
			}
			if eid, ok := payload["id"].(string); ok {
				externalEventID = eid
			}
			_, _ = db.Exec(`
				INSERT INTO billing_events (provider_id, event_type, external_event_id, payload, received_at)
				VALUES ($1, $2, $3, $4::jsonb, NOW())
				ON CONFLICT (provider_id, external_event_id) DO NOTHING
			`, providerID, eventType, externalEventID, string(body))
		}

		c.JSON(http.StatusOK, gin.H{"status": "received"})
	}
}

// getBillingRegistry builds a registry with configured providers.
func getBillingRegistry(c *gin.Context) *billing.Registry {
	// Pull config from server context if available; otherwise, re-load (fallback)
	cfg := config.Load()
	reg := billing.NewRegistry()
	reg.Register(billing.NewStripeProvider(cfg.StripeWebhookSecret))
	return reg
}
