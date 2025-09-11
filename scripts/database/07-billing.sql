-- =================================================================
-- Billing Schema (Provider-agnostic with Stripe initial adapter)
-- =================================================================

-- Billing providers (e.g., stripe, paddle, chargebee)
CREATE TABLE IF NOT EXISTS billing_providers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    key VARCHAR(50) UNIQUE NOT NULL, -- 'stripe', 'paddle', 'chargebee'
    display_name VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Billing customers (per tenant per provider)
CREATE TABLE IF NOT EXISTS billing_customers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    provider_id UUID NOT NULL REFERENCES billing_providers(id) ON DELETE RESTRICT,
    external_customer_id VARCHAR(255) NOT NULL, -- e.g., Stripe customer ID
    email VARCHAR(255),
    default_payment_method VARCHAR(255),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(tenant_id, provider_id)
);

-- Billing subscriptions
CREATE TABLE IF NOT EXISTS billing_subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    provider_id UUID NOT NULL REFERENCES billing_providers(id) ON DELETE RESTRICT,
    external_subscription_id VARCHAR(255) NOT NULL,
    plan_key VARCHAR(100) NOT NULL, -- maps to subscription_tiers.name or provider product/price
    status VARCHAR(50) NOT NULL, -- trialing, active, past_due, canceled, paused
    current_period_start TIMESTAMP WITH TIME ZONE,
    current_period_end TIMESTAMP WITH TIME ZONE,
    cancel_at_period_end BOOLEAN DEFAULT false,
    quantity INTEGER DEFAULT 1,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Billing invoices (high-level reference to provider invoices)
CREATE TABLE IF NOT EXISTS billing_invoices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    provider_id UUID NOT NULL REFERENCES billing_providers(id) ON DELETE RESTRICT,
    external_invoice_id VARCHAR(255) NOT NULL,
    amount_cents INTEGER NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'USD',
    status VARCHAR(50) NOT NULL, -- draft, open, paid, uncollectible, void
    issued_at TIMESTAMP WITH TIME ZONE,
    due_at TIMESTAMP WITH TIME ZONE,
    paid_at TIMESTAMP WITH TIME ZONE,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Billing events (webhook ingestions)
CREATE TABLE IF NOT EXISTS billing_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    provider_id UUID NOT NULL REFERENCES billing_providers(id) ON DELETE RESTRICT,
    event_type VARCHAR(100) NOT NULL,
    external_event_id VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    received_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    processed_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(provider_id, external_event_id)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_billing_customers_tenant ON billing_customers(tenant_id);
CREATE INDEX IF NOT EXISTS idx_billing_subscriptions_tenant ON billing_subscriptions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_billing_invoices_tenant ON billing_invoices(tenant_id);

-- Seed Stripe provider
INSERT INTO billing_providers (key, display_name, is_active)
VALUES ('stripe', 'Stripe', true)
ON CONFLICT (key) DO NOTHING;


