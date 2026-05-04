CREATE TABLE IF NOT EXISTS invoice_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    plan_type VARCHAR(50) NOT NULL DEFAULT 'MAINTENANCE',
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    issue_day_of_month INTEGER NOT NULL,
    due_day_of_month INTEGER NOT NULL,
    default_amount NUMERIC(12,2) NOT NULL,
    starts_on DATE NOT NULL,
    ends_on DATE,
    description TEXT,
    settings JSONB NOT NULL DEFAULT '{}',
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

ALTER TABLE invoice_plans
DROP CONSTRAINT IF EXISTS chk_invoice_plans_type;
ALTER TABLE invoice_plans
ADD CONSTRAINT chk_invoice_plans_type CHECK (
    plan_type IN ('MAINTENANCE')
);

ALTER TABLE invoice_plans
DROP CONSTRAINT IF EXISTS chk_invoice_plans_status;
ALTER TABLE invoice_plans
ADD CONSTRAINT chk_invoice_plans_status CHECK (
    status IN ('ACTIVE', 'PAUSED', 'ENDED')
);

ALTER TABLE invoice_plans
DROP CONSTRAINT IF EXISTS chk_invoice_plans_issue_day;
ALTER TABLE invoice_plans
ADD CONSTRAINT chk_invoice_plans_issue_day CHECK (
    issue_day_of_month BETWEEN 1 AND 28
);

ALTER TABLE invoice_plans
DROP CONSTRAINT IF EXISTS chk_invoice_plans_due_day;
ALTER TABLE invoice_plans
ADD CONSTRAINT chk_invoice_plans_due_day CHECK (
    due_day_of_month BETWEEN 1 AND 28
);

ALTER TABLE invoice_plans
DROP CONSTRAINT IF EXISTS chk_invoice_plans_default_amount;
ALTER TABLE invoice_plans
ADD CONSTRAINT chk_invoice_plans_default_amount CHECK (
    default_amount >= 0
);

ALTER TABLE invoice_plans
DROP CONSTRAINT IF EXISTS chk_invoice_plans_date_range;
ALTER TABLE invoice_plans
ADD CONSTRAINT chk_invoice_plans_date_range CHECK (
    ends_on IS NULL OR ends_on >= starts_on
);

DROP TRIGGER IF EXISTS update_invoice_plans_modtime ON invoice_plans;
CREATE TRIGGER update_invoice_plans_modtime
BEFORE UPDATE ON invoice_plans
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE INDEX IF NOT EXISTS idx_invoice_plans_community_status
ON invoice_plans(community_id, status);

CREATE TABLE IF NOT EXISTS invoice_plan_unit_overrides (
    invoice_plan_id UUID NOT NULL REFERENCES invoice_plans(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    unit_id UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    amount NUMERIC(12,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (invoice_plan_id, unit_id)
);

ALTER TABLE invoice_plan_unit_overrides
DROP CONSTRAINT IF EXISTS chk_invoice_plan_unit_overrides_amount;
ALTER TABLE invoice_plan_unit_overrides
ADD CONSTRAINT chk_invoice_plan_unit_overrides_amount CHECK (
    amount >= 0
);

DROP TRIGGER IF EXISTS update_invoice_plan_unit_overrides_modtime ON invoice_plan_unit_overrides;
CREATE TRIGGER update_invoice_plan_unit_overrides_modtime
BEFORE UPDATE ON invoice_plan_unit_overrides
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE INDEX IF NOT EXISTS idx_invoice_plan_unit_overrides_community_plan
ON invoice_plan_unit_overrides(community_id, invoice_plan_id);

ALTER TABLE invoices
ADD COLUMN IF NOT EXISTS invoice_plan_id UUID REFERENCES invoice_plans(id) ON DELETE SET NULL;

ALTER TABLE invoices
ADD COLUMN IF NOT EXISTS source VARCHAR(50) NOT NULL DEFAULT 'MANUAL';

ALTER TABLE invoices
ADD COLUMN IF NOT EXISTS issued_on DATE;

ALTER TABLE invoices
DROP CONSTRAINT IF EXISTS chk_invoices_source;
ALTER TABLE invoices
ADD CONSTRAINT chk_invoices_source CHECK (
    source IN ('MANUAL', 'SCHEDULED')
);

ALTER TABLE invoices
DROP CONSTRAINT IF EXISTS chk_invoices_plan_source_consistency;
ALTER TABLE invoices
ADD CONSTRAINT chk_invoices_plan_source_consistency CHECK (
    (source = 'MANUAL' AND invoice_plan_id IS NULL)
    OR (source = 'SCHEDULED' AND invoice_plan_id IS NOT NULL)
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_invoices_scheduled_plan_unit_period
ON invoices(invoice_plan_id, unit_id, billing_period)
WHERE source = 'SCHEDULED';

ALTER TABLE payments
ADD COLUMN IF NOT EXISTS status VARCHAR(50) NOT NULL DEFAULT 'PENDING';

ALTER TABLE payments
ADD COLUMN IF NOT EXISTS submitted_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE payments
ADD COLUMN IF NOT EXISTS recorded_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE payments
ADD COLUMN IF NOT EXISTS notes TEXT;

ALTER TABLE payments
ADD COLUMN IF NOT EXISTS evidence JSONB NOT NULL DEFAULT '{}';

ALTER TABLE payments
ADD COLUMN IF NOT EXISTS reviewed_at TIMESTAMPTZ;

ALTER TABLE payments
DROP CONSTRAINT IF EXISTS chk_payments_status;
ALTER TABLE payments
ADD CONSTRAINT chk_payments_status CHECK (
    status IN ('PENDING', 'APPROVED', 'REJECTED', 'RECORDED')
);

CREATE INDEX IF NOT EXISTS idx_payments_community_status
ON payments(community_id, status);
