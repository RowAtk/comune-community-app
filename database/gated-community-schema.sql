-- =========================================================
-- GATED COMMUNITY SAAS - FOUNDATION SCHEMA
-- Lightly commented version for implementation use
-- =========================================================

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- =========================================================
-- 1. Organizations
-- =========================================================
CREATE TABLE IF NOT EXISTS organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    legal_name VARCHAR(255),
    billing_email VARCHAR(255),
    phone VARCHAR(20),
    country_code VARCHAR(2) DEFAULT 'JM',
    timezone VARCHAR(50) NOT NULL DEFAULT 'America/Jamaica',
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    settings JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_organizations_slug UNIQUE (slug),
    CONSTRAINT chk_organizations_status CHECK (
        status IN ('ACTIVE', 'SUSPENDED', 'INACTIVE')
    )
);

CREATE TRIGGER update_organizations_modtime
BEFORE UPDATE ON organizations
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-- =========================================================
-- 2. Subscriptions
-- =========================================================
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    plan_code VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'TRIAL',
    billing_cycle VARCHAR(20) NOT NULL DEFAULT 'MONTHLY',
    current_period_start TIMESTAMPTZ,
    current_period_end TIMESTAMPTZ,
    cancel_at_period_end BOOLEAN NOT NULL DEFAULT FALSE,
    external_customer_id VARCHAR(255),
    external_subscription_id VARCHAR(255),
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_subscriptions_status CHECK (
        status IN ('TRIAL', 'ACTIVE', 'PAST_DUE', 'CANCELLED', 'EXPIRED')
    ),
    CONSTRAINT chk_subscriptions_billing_cycle CHECK (
        billing_cycle IN ('MONTHLY', 'QUARTERLY', 'YEARLY')
    )
);

CREATE TRIGGER update_subscriptions_modtime
BEFORE UPDATE ON subscriptions
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE INDEX IF NOT EXISTS idx_subscriptions_organization_id
ON subscriptions(organization_id);

-- =========================================================
-- 3. Communities
-- =========================================================
CREATE TABLE IF NOT EXISTS communities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    address TEXT,
    timezone VARCHAR(50) NOT NULL DEFAULT 'America/Jamaica',
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    settings JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_communities_org_slug UNIQUE (organization_id, slug),
    CONSTRAINT uq_communities_org_name UNIQUE (organization_id, name),
    CONSTRAINT chk_communities_status CHECK (
        status IN ('ACTIVE', 'INACTIVE', 'ARCHIVED')
    )
);

CREATE TRIGGER update_communities_modtime
BEFORE UPDATE ON communities
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE INDEX IF NOT EXISTS idx_communities_organization_id
ON communities(organization_id);

-- =========================================================
-- 4. Users
-- =========================================================
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL,
    --password_hash TEXT NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(20),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_users_email UNIQUE (email)
);

CREATE TABLE auth_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
    email_at_provider VARCHAR(255),
    password_hash TEXT,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT auth_accounts_provider_user_unique
        UNIQUE (provider, provider_user_id)
);

CREATE TRIGGER update_users_modtime
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE INDEX IF NOT EXISTS idx_users_active
ON users(id)
WHERE deleted_at IS NULL;

-- =========================================================
-- 5. Organization memberships
-- =========================================================
CREATE TABLE IF NOT EXISTS organization_users (
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    invited_by UUID REFERENCES users(id) ON DELETE SET NULL,
    joined_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (organization_id, user_id),
    CONSTRAINT chk_organization_users_role CHECK (
        role IN ('OWNER', 'ORG_ADMIN', 'BILLING_ADMIN', 'OPERATIONS', 'SUPPORT')
    ),
    CONSTRAINT chk_organization_users_status CHECK (
        status IN ('INVITED', 'ACTIVE', 'SUSPENDED', 'REMOVED')
    )
);

CREATE TRIGGER update_organization_users_modtime
BEFORE UPDATE ON organization_users
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE INDEX IF NOT EXISTS idx_organization_users_user_id
ON organization_users(user_id);

-- =========================================================
-- 6. Community memberships
-- =========================================================
CREATE TABLE IF NOT EXISTS community_users (
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    invited_by UUID REFERENCES users(id) ON DELETE SET NULL,
    joined_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (community_id, user_id),
    CONSTRAINT chk_community_users_role CHECK (
        role IN (
            'COMMUNITY_ADMIN',
            'RESIDENT',
            'SECURITY',
            'MANAGER',
            'BOARD_MEMBER',
            'MAINTENANCE_STAFF'
        )
    ),
    CONSTRAINT chk_community_users_status CHECK (
        status IN ('INVITED', 'ACTIVE', 'SUSPENDED', 'REMOVED')
    )
);

CREATE TRIGGER update_community_users_modtime
BEFORE UPDATE ON community_users
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE INDEX IF NOT EXISTS idx_community_users_user_id
ON community_users(user_id);

CREATE INDEX IF NOT EXISTS idx_community_users_org_id
ON community_users(organization_id);

-- =========================================================
-- 7. Invitations
-- =========================================================
CREATE TABLE IF NOT EXISTS organization_invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    token UUID NOT NULL DEFAULT gen_random_uuid(),
    expires_at TIMESTAMPTZ NOT NULL,
    accepted_at TIMESTAMPTZ,
    invited_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_organization_invitations_token UNIQUE (token)
);

CREATE TABLE IF NOT EXISTS community_invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    token UUID NOT NULL DEFAULT gen_random_uuid(),
    expires_at TIMESTAMPTZ NOT NULL,
    accepted_at TIMESTAMPTZ,
    invited_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_community_invitations_token UNIQUE (token)
);

-- =========================================================
-- 8. Units
-- =========================================================
CREATE TABLE IF NOT EXISTS units (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    unit_number VARCHAR(50) NOT NULL,
    block_floor VARCHAR(50),
    unit_type VARCHAR(50),
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_units_community_unit_number UNIQUE (community_id, unit_number),
    CONSTRAINT chk_units_status CHECK (status IN ('ACTIVE', 'INACTIVE'))
);

CREATE TRIGGER update_units_modtime
BEFORE UPDATE ON units
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE INDEX IF NOT EXISTS idx_units_organization_id
ON units(organization_id);

CREATE INDEX IF NOT EXISTS idx_units_community_id
ON units(community_id);

-- =========================================================
-- 9. Households
-- =========================================================
CREATE TABLE IF NOT EXISTS households (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    unit_id UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    name VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER update_households_modtime
BEFORE UPDATE ON households
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE INDEX IF NOT EXISTS idx_households_unit_id
ON households(unit_id);

-- =========================================================
-- 10. Residents
-- =========================================================
CREATE TABLE IF NOT EXISTS residents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    unit_id UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    household_id UUID REFERENCES households(id) ON DELETE SET NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(20),
    resident_type VARCHAR(50) NOT NULL,
    household_role VARCHAR(50),
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    is_primary_contact BOOLEAN NOT NULL DEFAULT FALSE,
    move_in_date DATE,
    move_out_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_residents_type CHECK (
        resident_type IN ('OWNER', 'TENANT', 'DEPENDENT', 'OCCUPANT')
    ),
    CONSTRAINT chk_residents_household_role CHECK (
        household_role IS NULL
        OR household_role IN ('HOUSEHOLD_ADMIN', 'HOUSEHOLD_MEMBER', 'HOUSEHOLD_VIEWER')
    ),
    CONSTRAINT chk_residents_status CHECK (
        status IN ('ACTIVE', 'INACTIVE', 'PENDING', 'MOVED_OUT')
    ),
    CONSTRAINT chk_residents_move_dates CHECK (
        move_out_date IS NULL
        OR move_in_date IS NULL
        OR move_out_date >= move_in_date
    )
);

CREATE TRIGGER update_residents_modtime
BEFORE UPDATE ON residents
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE INDEX IF NOT EXISTS idx_residents_household_id
ON residents(household_id);

CREATE INDEX IF NOT EXISTS idx_residents_community_id
ON residents(community_id);

CREATE INDEX IF NOT EXISTS idx_residents_user_id
ON residents(user_id);

CREATE UNIQUE INDEX IF NOT EXISTS uq_residents_primary_contact_per_household
ON residents(household_id)
WHERE is_primary_contact = TRUE AND deleted_at IS NULL;

-- =========================================================
-- 11. Vendors
-- =========================================================
CREATE TABLE IF NOT EXISTS vendors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    company_name VARCHAR(255) NOT NULL,
    category VARCHAR(100),
    contact_name VARCHAR(100),
    email VARCHAR(255),
    phone VARCHAR(20),
    address TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    rating NUMERIC(3,2) NOT NULL DEFAULT 0.00,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_vendors_status CHECK (
        status IN ('ACTIVE', 'INACTIVE', 'PROBATION')
    ),
    CONSTRAINT chk_vendors_rating CHECK (
        rating >= 0.00 AND rating <= 5.00
    )
);

CREATE TRIGGER update_vendors_modtime
BEFORE UPDATE ON vendors
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-- =========================================================
-- 12. Visitors
-- =========================================================
CREATE TABLE IF NOT EXISTS visitors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    unit_id UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    vehicle_number VARCHAR(50),
    expected_arrival TIMESTAMPTZ,
    invite_code VARCHAR(20) NOT NULL,
    purpose VARCHAR(255),
    status VARCHAR(50) NOT NULL DEFAULT 'SCHEDULED',
    qr_token UUID NOT NULL DEFAULT gen_random_uuid(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ,
    CONSTRAINT uq_visitors_invite_code UNIQUE (invite_code),
    CONSTRAINT uq_visitors_qr_token UNIQUE (qr_token),
    CONSTRAINT chk_visitors_status CHECK (
        status IN ('SCHEDULED', 'ARRIVED', 'DEPARTED', 'EXPIRED', 'CANCELLED')
    )
);

CREATE TRIGGER update_visitors_modtime
BEFORE UPDATE ON visitors
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE INDEX IF NOT EXISTS idx_visitors_invite_code
ON visitors(invite_code);

CREATE INDEX IF NOT EXISTS idx_visitors_qr_token
ON visitors(qr_token);

CREATE INDEX IF NOT EXISTS idx_visitors_community_status
ON visitors(community_id, status);

-- =========================================================
-- 13. Security logs
-- =========================================================
CREATE TABLE IF NOT EXISTS security_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    person_type VARCHAR(50) NOT NULL,
    visitor_id UUID REFERENCES visitors(id) ON DELETE SET NULL,
    resident_id UUID REFERENCES residents(id) ON DELETE SET NULL,
    vendor_id UUID REFERENCES vendors(id) ON DELETE SET NULL,
    unit_id UUID REFERENCES units(id) ON DELETE SET NULL,
    time_in TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    time_out TIMESTAMPTZ,
    guard_id UUID REFERENCES users(id) ON DELETE SET NULL,
    remarks TEXT,
    CONSTRAINT chk_security_logs_person_type CHECK (
        person_type IN ('RESIDENT', 'VISITOR', 'VENDOR')
    ),
    CONSTRAINT chk_security_logs_time_range CHECK (
        time_out IS NULL OR time_out >= time_in
    )
);

CREATE INDEX IF NOT EXISTS idx_security_logs_community_id
ON security_logs(community_id);

CREATE INDEX IF NOT EXISTS idx_security_logs_time_in
ON security_logs(time_in DESC);

-- =========================================================
-- 14. Maintenance
-- =========================================================
CREATE TABLE IF NOT EXISTS maintenance_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    unit_id UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    resident_id UUID NOT NULL REFERENCES residents(id) ON DELETE RESTRICT,
    vendor_id UUID REFERENCES vendors(id) ON DELETE SET NULL,
    assigned_to_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    status VARCHAR(50) NOT NULL DEFAULT 'OPEN',
    priority VARCHAR(50) NOT NULL DEFAULT 'MEDIUM',
    scheduled_date TIMESTAMPTZ,
    attachments JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_maintenance_status CHECK (
        status IN ('OPEN', 'IN_PROGRESS', 'COMPLETED', 'CANCELLED')
    ),
    CONSTRAINT chk_maintenance_priority CHECK (
        priority IN ('LOW', 'MEDIUM', 'HIGH', 'URGENT')
    )
);

CREATE TRIGGER update_maintenance_requests_modtime
BEFORE UPDATE ON maintenance_requests
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE INDEX IF NOT EXISTS idx_maintenance_requests_resident_id
ON maintenance_requests(resident_id);

CREATE INDEX IF NOT EXISTS idx_maintenance_requests_status
ON maintenance_requests(community_id, status);

-- =========================================================
-- 15. Operational billing
-- =========================================================
CREATE TABLE IF NOT EXISTS invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    unit_id UUID NOT NULL REFERENCES units(id) ON DELETE RESTRICT,
    amount NUMERIC(12,2) NOT NULL,
    paid_amount NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    due_date DATE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'UNPAID',
    billing_period VARCHAR(20) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_invoices_amount CHECK (amount >= 0),
    CONSTRAINT chk_invoices_paid_amount CHECK (
        paid_amount >= 0 AND paid_amount <= amount
    ),
    CONSTRAINT chk_invoices_status CHECK (
        status IN ('UNPAID', 'PARTIAL', 'PAID', 'OVERDUE', 'VOID')
    )
);

CREATE TRIGGER update_invoices_modtime
BEFORE UPDATE ON invoices
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE INDEX IF NOT EXISTS idx_invoices_unit_id
ON invoices(unit_id);

CREATE INDEX IF NOT EXISTS idx_invoices_status
ON invoices(community_id, status);

CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE RESTRICT,
    amount NUMERIC(12,2) NOT NULL,
    payment_method VARCHAR(50),
    transaction_id VARCHAR(255),
    payment_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_payments_amount CHECK (amount > 0)
);

CREATE INDEX IF NOT EXISTS idx_payments_invoice_id
ON payments(invoice_id);

-- =========================================================
-- 16. Communication
-- =========================================================
CREATE TABLE IF NOT EXISTS announcements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    author_id UUID REFERENCES users(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    audience_type VARCHAR(50) NOT NULL DEFAULT 'ALL',
    priority BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_announcements_audience_type CHECK (
        audience_type IN ('ALL', 'OWNERS', 'TENANTS', 'SECURITY', 'STAFF')
    )
);

CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    link_to_resource VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_unread
ON notifications(user_id)
WHERE is_read = FALSE;

-- =========================================================
-- 17. Audit / outbox
-- =========================================================
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    actor_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    entity_type VARCHAR(100) NOT NULL,
    entity_id UUID NOT NULL,
    action VARCHAR(100) NOT NULL,
    old_values JSONB NOT NULL DEFAULT '{}',
    new_values JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS outbox_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    event_type VARCHAR(100) NOT NULL,
    aggregate_type VARCHAR(100) NOT NULL,
    aggregate_id UUID NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    available_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_outbox_status CHECK (
        status IN ('PENDING', 'PROCESSING', 'PROCESSED', 'FAILED')
    )
);

CREATE INDEX IF NOT EXISTS idx_outbox_events_status_available_at
ON outbox_events(status, available_at);

-- =========================================================
-- 18. Example RLS policy
-- =========================================================
ALTER TABLE announcements ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS organization_community_isolation_policy ON announcements;

CREATE POLICY organization_community_isolation_policy
ON announcements
USING (
    organization_id = current_setting('app.current_organization_id')::uuid
    AND community_id = current_setting('app.current_community_id')::uuid
);
