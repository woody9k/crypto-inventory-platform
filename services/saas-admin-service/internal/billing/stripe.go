package billing

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"net/http"
)

type StripeProvider struct {
	secret string // webhook signing secret
}

func NewStripeProvider(secret string) *StripeProvider {
	return &StripeProvider{secret: secret}
}

func (s *StripeProvider) Key() string { return "stripe" }

// VerifyWebhook performs a minimal signature check compatible with Stripe-style signatures.
func (s *StripeProvider) VerifyWebhook(headers http.Header, body []byte) error {
	// NOTE: This is a placeholder. In production, use official Stripe SDK and verify 'Stripe-Signature'
	sig := headers.Get("Stripe-Signature")
	if sig == "" {
		return errors.New("missing Stripe-Signature header")
	}

	mac := hmac.New(sha256.New, []byte(s.secret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	if subtleConstantTimeCompare(expected, sig) == false {
		return errors.New("invalid webhook signature")
	}
	return nil
}

func (s *StripeProvider) ApplyUpdate(db *sql.DB, tenantID string, req UpdateRequest) error {
	// Placeholder: persist intent and let an async worker call Stripe API
	// For now, we simulate success.
	return nil
}

func subtleConstantTimeCompare(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	var v byte
	for i := 0; i < len(a); i++ {
		v |= a[i] ^ b[i]
	}
	return v == 0
}
