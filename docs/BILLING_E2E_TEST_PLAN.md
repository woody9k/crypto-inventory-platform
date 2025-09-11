# Billing E2E Test Plan

## Scenarios

- Trial → Paid upgrade
  - Create tenant (trial) and verify default subscription.
  - Change plan to `professional` via PUT `/admin/tenants/:id/billing`.
  - Simulate webhook `customer.subscription.updated` with active status.
  - Verify `billing_subscriptions` and `tenants.subscription_tier_id` aligned.

- Downgrade / Upgrade with proration
  - Issue change_plan to `enterprise` then back to `professional`.
  - Verify subscription status remains `active` and periods update.

- Cancel at period end
  - Toggle `cancel_at_period_end` true.
  - Verify flag in `billing_subscriptions` and UI reflects state.

- Resume subscription
  - Toggle `cancel_at_period_end` false.
  - Verify updated flag.

- Invoices listing
  - Insert invoices for tenant; GET `/admin/billing/invoices?tenantId=...`.
  - Verify amounts, currency, and statuses render in UI.

- Webhook security
  - Send webhook without signature → 401.
  - Send with invalid signature → 401.
  - Send with valid signature → 200 and event persisted.

## Validation Checklist

- DB:
  - `billing_events` row inserted per webhook (unique by provider/event id).
  - `billing_subscriptions` upserted with correct status/periods/plan.
  - `tenants.subscription_tier_id` updated when plan maps.

- API:
  - Billing read/update endpoints enforce role guards.
  - Invoices endpoint supports global and per-tenant.

- UI:
  - Tenant Billing page shows current plan and cancel flag.
  - Invoices page lists invoices and supports tenant filter.


