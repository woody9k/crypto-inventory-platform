package billing

import (
	"database/sql"
	"net/http"
)

// UpdateRequest models a normalized billing update request
type UpdateRequest struct {
	Action            string  `json:"action"`   // "change_plan", "cancel", "resume"
	PlanKey           *string `json:"plan_key"` // new plan when changing
	CancelAtPeriodEnd *bool   `json:"cancel_at_period_end"`
}

// BillingProvider is a provider-agnostic interface
type BillingProvider interface {
	Key() string
	VerifyWebhook(headers http.Header, body []byte) error
	ApplyUpdate(db *sql.DB, tenantID string, req UpdateRequest) error
}
