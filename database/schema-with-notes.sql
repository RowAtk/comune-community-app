-- =========================================================
-- GATED COMMUNITY SAAS - FOUNDATION SCHEMA
-- =========================================================
--
-- High-level model:
--
--   organizations -> the real SaaS tenant / paying customer
--   communities   -> one or more gated communities owned by an organization
--   users         -> global login identities
--   memberships   -> organization-level and community-level access
--   domain tables -> resident, visitor, security, maintenance, invoicing, etc.
--
-- Why this model:
--
-- 1. The subscription is billed to an organization, not an individual community.
-- 2. One organization can manage multiple communities.
-- 3. A single user can belong to an organization and possibly multiple communities.
-- 4. Operational tables keep both organization_id and community_id for:
--      - easier tenant scoping
--      - simpler authorization logic
--      - better indexing
--      - stronger audit clarity
--
-- =========================================================
-- EXTENSIONS
-- =========================================================

-- pgcrypto gives us gen_random_uuid(), which is the UUID generator used below.
-- This is preferred over older uuid-ossp in many modern PostgreSQL setups.
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- =========================================================
-- SHARED TRIGGER FUNCTION
-- =========================================================

-- This function is used by many tables to keep updated_at current.
-- Every UPDATE automatically refreshes updated_at.
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- =========================================================
-- 1. ORGANIZATIONS
-- =========================================================
--
-- organizations = top-level SaaS tenant / customer account / billing owner
--
-- Example:
--   "Island Property Management Ltd"
--
-- An organization can own multiple communities.
-- Subscription/billing is attached at this level.
--
CREATE TABLE IF NOT EXISTS organizations (
    -- Global primary key for the organization.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Human-readable customer name.
    -- Example: "Palm View Property Services Ltd"
    name VARCHAR(255) NOT NULL,

    -- URL-safe stable identifier for routing, slugs, public admin URLs, etc.
    -- Must be unique across all organizations.
    -- Example: "palm-view-property-services"
    slug VARCHAR(100) NOT NULL,

    -- Optional legal/business registered name if different from display name.
    legal_name VARCHAR(255),

    -- Billing contact email for invoices/subscription notices.
    billing_email VARCHAR(255),

    -- General contact phone for the organization.
    phone VARCHAR(20),

    -- Two-letter country code.
    -- Defaults to JM because your primary market is Jamaica.
    country_code VARCHAR(2) DEFAULT 'JM',

    -- Default timezone for org-level operations and default creation behavior.
    timezone VARCHAR(50) NOT NULL DEFAULT 'America/Jamaica',

    -- Current organization lifecycle status.
    -- ACTIVE     = usable customer account
    -- SUSPENDED  = temporarily disabled, e.g. billing issue
    -- INACTIVE   = no longer active / closed
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',

    -- Flexible org-level configuration.
    -- Good for settings that are not worth first-class columns yet.
    -- Keep this for settings, not core operational data.
    settings JSONB NOT NULL DEFAULT '{}',

    -- Audit timestamps.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Soft delete timestamp.
    -- Useful if an organization is retired but historical records should remain.
    deleted_at TIMESTAMPTZ,

    -- Slug must be globally unique.
    CONSTRAINT uq_organizations_slug UNIQUE (slug),

    -- Restricts allowed organization statuses.
    CONSTRAINT chk_organizations_status CHECK (
        status IN ('ACTIVE', 'SUSPENDED', 'INACTIVE')
    )
);

-- Keep updated_at in sync for organizations.
CREATE TRIGGER update_organizations_modtime
BEFORE UPDATE ON organizations
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-- =========================================================
-- 2. SUBSCRIPTIONS
-- =========================================================
--
-- subscriptions = SaaS billing/subscription ownership at organization level
--
-- This is NOT the community operational billing table.
-- This is for YOUR platform subscription billing.
--
-- Example:
--   Organization is on plan "pro"
--   current billing cycle monthly
--   status ACTIVE
--
CREATE TABLE IF NOT EXISTS subscriptions (
    -- Unique subscription record identifier.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Owning organization.
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Internal plan code.
    -- Example: STARTER, PRO, ENTERPRISE
    plan_code VARCHAR(100) NOT NULL,

    -- Subscription status.
    -- TRIAL      = active trial period
    -- ACTIVE     = paying active subscription
    -- PAST_DUE   = failed renewal / payment issue
    -- CANCELLED  = manually cancelled
    -- EXPIRED    = trial or subscription ended
    status VARCHAR(50) NOT NULL DEFAULT 'TRIAL',

    -- Billing interval.
    billing_cycle VARCHAR(20) NOT NULL DEFAULT 'MONTHLY',

    -- Current paid or trial period start.
    current_period_start TIMESTAMPTZ,

    -- Current paid or trial period end.
    current_period_end TIMESTAMPTZ,

    -- If true, do not renew after current period ends.
    cancel_at_period_end BOOLEAN NOT NULL DEFAULT FALSE,

    -- External billing provider customer ID.
    -- Example: Stripe customer ID.
    external_customer_id VARCHAR(255),

    -- External billing provider subscription ID.
    external_subscription_id VARCHAR(255),

    -- Flexible metadata from billing systems if needed.
    metadata JSONB NOT NULL DEFAULT '{}',

    -- Audit timestamps.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Restrict subscription lifecycle values.
    CONSTRAINT chk_subscriptions_status CHECK (
        status IN ('TRIAL', 'ACTIVE', 'PAST_DUE', 'CANCELLED', 'EXPIRED')
    ),

    -- Restrict billing interval values.
    CONSTRAINT chk_subscriptions_billing_cycle CHECK (
        billing_cycle IN ('MONTHLY', 'QUARTERLY', 'YEARLY')
    )
);

CREATE TRIGGER update_subscriptions_modtime
BEFORE UPDATE ON subscriptions
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-- Useful for fetching all subscriptions for an org quickly.
CREATE INDEX IF NOT EXISTS idx_subscriptions_organization_id
ON subscriptions(organization_id);

-- =========================================================
-- 3. COMMUNITIES
-- =========================================================
--
-- communities = gated communities / properties / schemes under an organization
--
-- Example:
--   Organization: Island Property Management
--   Communities:
--     - Palm View Estate
--     - Kingston Heights
--
CREATE TABLE IF NOT EXISTS communities (
    -- Community primary key.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Parent organization that owns/manages this community.
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Display name of the community.
    name VARCHAR(255) NOT NULL,

    -- Slug unique only inside the organization.
    -- Two different organizations can both have "main-estate";
    -- one organization cannot have two of them.
    slug VARCHAR(100) NOT NULL,

    -- Optional physical address.
    address TEXT,

    -- Community-specific timezone.
    -- Often inherited from organization, but stored explicitly here for clarity.
    timezone VARCHAR(50) NOT NULL DEFAULT 'America/Jamaica',

    -- Community lifecycle status.
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',

    -- Community-specific settings/configuration.
    settings JSONB NOT NULL DEFAULT '{}',

    -- Audit columns.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    -- Enforce unique slug per organization.
    CONSTRAINT uq_communities_org_slug UNIQUE (organization_id, slug),

    -- Enforce unique name per organization.
    CONSTRAINT uq_communities_org_name UNIQUE (organization_id, name),

    -- Restrict lifecycle values.
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
-- 4. USERS
-- =========================================================
--
-- users = platform login identities
--
-- Important:
-- This table is NOT only for residents.
-- A user may be:
--   - org admin
--   - billing admin
--   - community admin
--   - resident
--   - security officer
--   - maintenance staff
--
-- Residents are modeled separately in the residents table because:
--   - not every resident needs a login
--   - not every user is a resident
--   - resident is a business/domain concept
--   - user is an auth/identity concept
--
CREATE TABLE IF NOT EXISTS users (
    -- User primary key.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Login email.
    -- Globally unique in this model.
    email VARCHAR(255) NOT NULL,

    -- Secure password hash. Never store plain text passwords.
    password_hash TEXT NOT NULL,

    -- Personal profile fields.
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(20),

    -- Simple active flag for auth-level disable/enable.
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    -- Last successful authentication timestamp.
    last_login_at TIMESTAMPTZ,

    -- Audit columns.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    -- Prevent duplicate login accounts by email.
    CONSTRAINT uq_users_email UNIQUE (email)
);

CREATE TRIGGER update_users_modtime
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-- Partial index for active/non-deleted users.
CREATE INDEX IF NOT EXISTS idx_users_active
ON users(id)
WHERE deleted_at IS NULL;

-- =========================================================
-- 5. ORGANIZATION MEMBERSHIPS
-- =========================================================
--
-- organization_users = users granted access at organization/customer level
--
-- This is where org-wide roles live.
-- Example roles:
--   OWNER
--   ORG_ADMIN
--   BILLING_ADMIN
--   OPERATIONS
--
-- Why separate this from community_users:
--   - billing admin might not belong to every community
--   - org owner should exist above community-specific access
--   - organization-level reporting/settings need a separate scope
--
CREATE TABLE IF NOT EXISTS organization_users (
    -- Parent organization.
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- User in that organization.
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Organization-level role.
    role VARCHAR(50) NOT NULL,

    -- Membership lifecycle state.
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',

    -- Which user sent the invitation or added this membership.
    invited_by UUID REFERENCES users(id) ON DELETE SET NULL,

    -- When the person effectively joined.
    joined_at TIMESTAMPTZ,

    -- Audit timestamps.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- A user can only appear once per organization.
    PRIMARY KEY (organization_id, user_id),

    -- Restrict allowed organization roles.
    CONSTRAINT chk_organization_users_role CHECK (
        role IN ('OWNER', 'ORG_ADMIN', 'BILLING_ADMIN', 'OPERATIONS', 'SUPPORT')
    ),

    -- Restrict allowed membership states.
    CONSTRAINT chk_organization_users_status CHECK (
        status IN ('INVITED', 'ACTIVE', 'SUSPENDED', 'REMOVED')
    )
);

CREATE TRIGGER update_organization_users_modtime
BEFORE UPDATE ON organization_users
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-- Useful when looking up a user's organizations.
CREATE INDEX IF NOT EXISTS idx_organization_users_user_id
ON organization_users(user_id);

-- =========================================================
-- 6. COMMUNITY MEMBERSHIPS
-- =========================================================
--
-- community_users = users granted access inside a specific community
--
-- Examples:
--   COMMUNITY_ADMIN
--   RESIDENT
--   SECURITY
--   BOARD_MEMBER
--
-- organization_id is stored explicitly here even though the community already
-- points to an organization. This is a deliberate denormalization for:
--   - easier tenant filtering
--   - easier indexing
--   - simpler authorization logic
--
CREATE TABLE IF NOT EXISTS community_users (
    -- Owning organization for the membership row.
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Target community.
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,

    -- User receiving the community membership.
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Community-scoped role.
    role VARCHAR(50) NOT NULL,

    -- Membership lifecycle.
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',

    -- Which user invited/added this community member.
    invited_by UUID REFERENCES users(id) ON DELETE SET NULL,

    -- When the member joined the community.
    joined_at TIMESTAMPTZ,

    -- Audit timestamps.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- One user cannot have duplicate membership rows in the same community.
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
-- 7. INVITATIONS
-- =========================================================
--
-- These tables support onboarding and invite flows.
--
-- organization_invitations:
--   invites someone into the organization/customer account
--
-- community_invitations:
--   invites someone into a specific community
--
-- Token is UUID-based and should be used in invite acceptance flows.
--
CREATE TABLE IF NOT EXISTS organization_invitations (
    -- Invitation primary key.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Organization receiving the member.
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Email invited.
    email VARCHAR(255) NOT NULL,

    -- Role to be granted when accepted.
    role VARCHAR(50) NOT NULL,

    -- Secure token used in acceptance URL/workflow.
    token UUID NOT NULL DEFAULT gen_random_uuid(),

    -- Invite expiry time.
    expires_at TIMESTAMPTZ NOT NULL,

    -- When accepted.
    accepted_at TIMESTAMPTZ,

    -- Inviter.
    invited_by UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Creation timestamp.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_organization_invitations_token UNIQUE (token)
);

CREATE TABLE IF NOT EXISTS community_invitations (
    -- Invitation primary key.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Organization owning the community.
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Community being joined.
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,

    -- Email invited.
    email VARCHAR(255) NOT NULL,

    -- Community role to grant.
    role VARCHAR(50) NOT NULL,

    -- Secure invite token.
    token UUID NOT NULL DEFAULT gen_random_uuid(),

    -- Invite expiry time.
    expires_at TIMESTAMPTZ NOT NULL,

    -- Acceptance timestamp.
    accepted_at TIMESTAMPTZ,

    -- Inviter.
    invited_by UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Creation timestamp.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_community_invitations_token UNIQUE (token)
);

-- =========================================================
-- 8. UNITS
-- =========================================================
--
-- units = physical residential units inside a community
--
-- Example:
--   Unit 12A
--   Villa 5
--   Block B / Floor 2 / Apt 6
--
CREATE TABLE IF NOT EXISTS units (
    -- Unit primary key.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Tenant scope.
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Community containing the unit.
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,

    -- Human-readable unit identifier.
    unit_number VARCHAR(50) NOT NULL,

    -- Optional building/block/floor notation.
    block_floor VARCHAR(50),

    -- Type/category of unit.
    -- Example: APARTMENT, VILLA, TOWNHOUSE
    unit_type VARCHAR(50),

    -- Lifecycle status for the unit record.
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',

    -- Audit columns.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    -- Unit number should be unique within a community.
    CONSTRAINT uq_units_community_unit_number UNIQUE (community_id, unit_number),

    CONSTRAINT chk_units_status CHECK (
        status IN ('ACTIVE', 'INACTIVE')
    )
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
-- 9. HOUSEHOLDS
-- =========================================================
--
-- households = grouping of people occupying/living together in a unit
--
-- This allows:
--   - multiple residents in one unit
--   - one primary contact per household
--   - family grouping logic
--
CREATE TABLE IF NOT EXISTS households (
    -- Household primary key.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Tenant scope.
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Community scope.
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,

    -- Unit for the household.
    unit_id UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,

    -- Optional display name.
    -- Example: "The Smith Household"
    name VARCHAR(255),

    -- Audit columns.
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
-- 10. RESIDENTS
-- =========================================================
--
-- residents = business/domain representation of people associated with a unit
--
-- Important:
-- A resident may or may not have a linked user account.
--
-- Examples:
--   - owner with login account
--   - tenant with login account
--   - dependent without login account
--
CREATE TABLE IF NOT EXISTS residents (
    -- Resident primary key.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Tenant scope.
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Community scope.
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,

    -- Unit occupied/owned by the resident.
    unit_id UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,

    -- Optional linked login identity.
    -- Nullable because not all residents need user accounts.
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Optional household grouping.
    household_id UUID REFERENCES households(id) ON DELETE SET NULL,

    -- Resident profile data.
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(20),

    -- Resident role in the property.
    -- OWNER, TENANT, DEPENDENT, OCCUPANT, etc.
    resident_type VARCHAR(50) NOT NULL,

    -- Optional application access level within the household.
    -- Null means the resident has the lowest/default household privileges.
    household_role VARCHAR(50),

    -- Current lifecycle state.
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',

    -- Whether this resident is the primary contact for the household.
    is_primary_contact BOOLEAN NOT NULL DEFAULT FALSE,

    -- Occupancy dates.
    move_in_date DATE,
    move_out_date DATE,

    -- Audit fields.
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

    -- If both dates exist, move_out cannot be before move_in.
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

-- Only one active primary contact per household.
-- Soft-deleted rows do not participate in the uniqueness rule.
CREATE UNIQUE INDEX IF NOT EXISTS uq_residents_primary_contact_per_household
ON residents(household_id)
WHERE is_primary_contact = TRUE AND deleted_at IS NULL;

-- =========================================================
-- 11. RESIDENT INVITATIONS
-- =========================================================
--
-- resident_invitations = token-based flows for linking a login account to an
-- existing resident record.
--
-- Why this exists:
--
-- 1. Community admins often create the resident record first.
-- 2. The resident may sign up later using a shared token or invite link.
-- 3. Accepting the invite should both link the user to the resident record and
--    grant active community membership.
--
-- This is intentionally separate from community_invitations:
--
-- - community_invitations grant general community membership
-- - resident_invitations bind the accepted account to a specific resident
--
CREATE TABLE IF NOT EXISTS resident_invitations (
    -- Invitation primary key.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Tenant scope.
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Community scope.
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,

    -- Resident that will be claimed/linked.
    resident_id UUID NOT NULL REFERENCES residents(id) ON DELETE CASCADE,

    -- Secure invite token.
    token UUID NOT NULL DEFAULT gen_random_uuid(),

    -- Invite expiry time.
    expires_at TIMESTAMPTZ NOT NULL,

    -- When the invite was accepted.
    accepted_at TIMESTAMPTZ,

    -- Admin/member who created the invite.
    invited_by UUID REFERENCES users(id) ON DELETE SET NULL,

    -- User who accepted the invite.
    accepted_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Creation timestamp.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_resident_invitations_token UNIQUE (token)
);

CREATE INDEX IF NOT EXISTS idx_resident_invitations_resident_id
ON resident_invitations(resident_id);

CREATE INDEX IF NOT EXISTS idx_resident_invitations_community_id
ON resident_invitations(organization_id, community_id);

-- =========================================================
-- 12. VENDORS
-- =========================================================
--
-- vendors = service providers known to the community
--
-- Examples:
--   - plumber
--   - electrician
--   - landscaping provider
--   - pest control
--
CREATE TABLE IF NOT EXISTS vendors (
    -- Vendor primary key.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Tenant scope.
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Community scope.
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,

    -- Company/business name.
    company_name VARCHAR(255) NOT NULL,

    -- High-level category.
    category VARCHAR(100),

    -- Primary contact person.
    contact_name VARCHAR(100),

    -- Contact details.
    email VARCHAR(255),
    phone VARCHAR(20),
    address TEXT,

    -- Vendor lifecycle state.
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',

    -- Rating score. Restrained to 0.00 - 5.00.
    rating NUMERIC(3,2) NOT NULL DEFAULT 0.00,

    -- Free-form notes.
    notes TEXT,

    -- Audit fields.
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
-- 12. VISITORS
-- =========================================================
--
-- visitors = invited visitor entries with pass/access data
--
-- Current design choice:
-- This table combines the visitor identity and the invitation/pass record.
-- That is acceptable for MVP.
--
-- Later you may split this into:
--   visitors
--   visitor_passes
--   visitor_access_events
--
CREATE TABLE IF NOT EXISTS visitors (
    -- Visitor/invitation primary key.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Tenant scope.
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Community scope.
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,

    -- Unit expecting the visitor.
    unit_id UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,

    -- Visitor profile fields.
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    vehicle_number VARCHAR(50),

    -- Optional expected arrival time set by resident/admin.
    expected_arrival TIMESTAMPTZ,

    -- Shorter human-friendly invite code.
    -- Useful for manual entry by guards.
    invite_code VARCHAR(20) NOT NULL,

    -- Why the visitor is coming.
    purpose VARCHAR(255),

    -- Visit/access lifecycle state.
    status VARCHAR(50) NOT NULL DEFAULT 'SCHEDULED',

    -- Unique QR-related token for scan-based access.
    qr_token UUID NOT NULL DEFAULT gen_random_uuid(),

    -- Which user created the visitor entry.
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Audit fields.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Optional expiry timestamp.
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
-- 13. SECURITY LOGS
-- =========================================================
--
-- security_logs = record of entry/exit activity or logged security interaction
--
-- person_type tells you which foreign key is expected to be relevant:
--   VISITOR  -> visitor_id
--   RESIDENT -> resident_id
--   VENDOR   -> vendor_id
--
-- Later, you may split this into:
--   - access event logs
--   - incident/security reports
--
CREATE TABLE IF NOT EXISTS security_logs (
    -- Security log primary key.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Tenant scope.
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Community scope.
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,

    -- Entity type entering/leaving.
    person_type VARCHAR(50) NOT NULL,

    -- Optional pointers to the specific person/entity involved.
    visitor_id UUID REFERENCES visitors(id) ON DELETE SET NULL,
    resident_id UUID REFERENCES residents(id) ON DELETE SET NULL,
    vendor_id UUID REFERENCES vendors(id) ON DELETE SET NULL,

    -- Associated unit if relevant.
    unit_id UUID REFERENCES units(id) ON DELETE SET NULL,

    -- Access timestamps.
    time_in TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    time_out TIMESTAMPTZ,

    -- Security guard / user who logged the event.
    guard_id UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Free-form remarks.
    remarks TEXT,

    CONSTRAINT chk_security_logs_person_type CHECK (
        person_type IN ('RESIDENT', 'VISITOR', 'VENDOR')
    ),

    -- time_out cannot be earlier than time_in.
    CONSTRAINT chk_security_logs_time_range CHECK (
        time_out IS NULL OR time_out >= time_in
    )
);

CREATE INDEX IF NOT EXISTS idx_security_logs_community_id
ON security_logs(community_id);

CREATE INDEX IF NOT EXISTS idx_security_logs_time_in
ON security_logs(time_in DESC);

-- =========================================================
-- 14. MAINTENANCE REQUESTS
-- =========================================================
--
-- maintenance_requests = resident or admin-requested maintenance work
--
-- Includes:
--   - resident requester
--   - related unit
--   - optional vendor assignment
--   - optional internal assignee
--   - category / priority / status
--
CREATE TABLE IF NOT EXISTS maintenance_requests (
    -- Maintenance request primary key.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Tenant scope.
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Community scope.
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,

    -- Unit affected.
    unit_id UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,

    -- Requesting/associated resident.
    resident_id UUID NOT NULL REFERENCES residents(id) ON DELETE RESTRICT,

    -- External vendor assigned, if any.
    vendor_id UUID REFERENCES vendors(id) ON DELETE SET NULL,

    -- Internal user assigned, if any.
    assigned_to_user_id UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Short summary/title of the issue.
    title VARCHAR(255) NOT NULL,

    -- Detailed description.
    description TEXT,

    -- Category label.
    category VARCHAR(100),

    -- Lifecycle status.
    status VARCHAR(50) NOT NULL DEFAULT 'OPEN',

    -- Priority/severity.
    priority VARCHAR(50) NOT NULL DEFAULT 'MEDIUM',

    -- Optional scheduled work date/time.
    scheduled_date TIMESTAMPTZ,

    -- Simple attachment metadata bucket for MVP.
    -- Later this may become a normalized attachments table.
    attachments JSONB NOT NULL DEFAULT '[]',

    -- Audit columns.
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
-- 15. COMMUNITY OPERATIONAL INVOICING
-- =========================================================
--
-- IMPORTANT:
-- These tables represent invoicing INSIDE the product for the community.
-- Example:
--   - maintenance fee invoice
--   - community dues invoice
--   - resident payment
--
-- This is separate from the SaaS subscription tables above.
-- The intended module boundary is `invoicing`, not maintenance requests.
--
-- Current design is MVP-friendly:
--   recurring plan rules stay separate from generated invoices
--   payments directly reference invoices
--
-- Later, for more flexible accounting, you may split into:
--   invoice_plans
--   invoice_plan_unit_overrides
--   invoices
--   invoice_items
--   payments
--   payment_allocations
--
CREATE TABLE IF NOT EXISTS invoice_plans (
    -- Invoice plan primary key.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Tenant scope.
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Community scope.
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,

    -- Human-friendly plan name.
    name VARCHAR(255) NOT NULL,

    -- Invoicing plan category for future expansion.
    plan_type VARCHAR(50) NOT NULL DEFAULT 'MAINTENANCE',

    -- Plan lifecycle.
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',

    -- Day of month when invoices should be issued.
    -- Constrained to 1..28 in v1 to avoid month-length ambiguity.
    issue_day_of_month INTEGER NOT NULL,

    -- Day of month when payment becomes due.
    -- Also constrained to 1..28 for predictable scheduling.
    due_day_of_month INTEGER NOT NULL,

    -- Default amount charged to units without an override.
    default_amount NUMERIC(12,2) NOT NULL,

    -- First date the plan is effective.
    starts_on DATE NOT NULL,

    -- Optional end date for the plan.
    ends_on DATE,

    -- Optional operational notes/description.
    description TEXT,

    -- Flexible non-core settings.
    settings JSONB NOT NULL DEFAULT '{}',

    -- User who created the plan.
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Audit columns.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT chk_invoice_plans_type CHECK (
        plan_type IN ('MAINTENANCE')
    ),

    CONSTRAINT chk_invoice_plans_status CHECK (
        status IN ('ACTIVE', 'PAUSED', 'ENDED')
    ),

    CONSTRAINT chk_invoice_plans_issue_day CHECK (
        issue_day_of_month BETWEEN 1 AND 28
    ),

    CONSTRAINT chk_invoice_plans_due_day CHECK (
        due_day_of_month BETWEEN 1 AND 28
    ),

    CONSTRAINT chk_invoice_plans_default_amount CHECK (
        default_amount >= 0
    ),

    CONSTRAINT chk_invoice_plans_date_range CHECK (
        ends_on IS NULL OR ends_on >= starts_on
    )
);

CREATE TRIGGER update_invoice_plans_modtime
BEFORE UPDATE ON invoice_plans
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE INDEX IF NOT EXISTS idx_invoice_plans_community_status
ON invoice_plans(community_id, status);

CREATE TABLE IF NOT EXISTS invoice_plan_unit_overrides (
    -- One plan-specific override row per unit.
    invoice_plan_id UUID NOT NULL REFERENCES invoice_plans(id) ON DELETE CASCADE,

    -- Tenant scope repeated for operational clarity.
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Community scope repeated for operational clarity.
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,

    -- Unit receiving custom pricing.
    unit_id UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,

    -- Override amount for the unit.
    amount NUMERIC(12,2) NOT NULL,

    -- Audit columns.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (invoice_plan_id, unit_id),

    CONSTRAINT chk_invoice_plan_unit_overrides_amount CHECK (
        amount >= 0
    )
);

CREATE TRIGGER update_invoice_plan_unit_overrides_modtime
BEFORE UPDATE ON invoice_plan_unit_overrides
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE INDEX IF NOT EXISTS idx_invoice_plan_unit_overrides_community_plan
ON invoice_plan_unit_overrides(community_id, invoice_plan_id);

CREATE TABLE IF NOT EXISTS invoices (
    -- Invoice primary key.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Tenant scope.
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Community scope.
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,

    -- Unit being billed.
    unit_id UUID NOT NULL REFERENCES units(id) ON DELETE RESTRICT,

    -- Recurring plan source, if this invoice was generated automatically.
    invoice_plan_id UUID REFERENCES invoice_plans(id) ON DELETE SET NULL,

    -- Whether the invoice was created manually or generated from a plan.
    source VARCHAR(50) NOT NULL DEFAULT 'MANUAL',

    -- Total invoice amount.
    amount NUMERIC(12,2) NOT NULL,

    -- Amount already paid against this invoice.
    paid_amount NUMERIC(12,2) NOT NULL DEFAULT 0.00,

    -- Business-visible issue date if distinct from row creation time.
    issued_on DATE,

    -- Due date for payment.
    due_date DATE NOT NULL,

    -- Invoice lifecycle.
    status VARCHAR(50) NOT NULL DEFAULT 'UNPAID',

    -- Billing period identifier.
    -- Example: 2026-04
    billing_period VARCHAR(20) NOT NULL,

    -- Optional description/details.
    description TEXT,

    -- Audit timestamps.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Prevent invalid negative invoice totals.
    CONSTRAINT chk_invoices_amount CHECK (amount >= 0),

    -- Paid amount must remain within valid bounds.
    CONSTRAINT chk_invoices_paid_amount CHECK (
        paid_amount >= 0 AND paid_amount <= amount
    ),

    CONSTRAINT chk_invoices_source CHECK (
        source IN ('MANUAL', 'SCHEDULED')
    ),

    CONSTRAINT chk_invoices_status CHECK (
        status IN ('UNPAID', 'PARTIAL', 'PAID', 'OVERDUE', 'VOID')
    ),

    CONSTRAINT chk_invoices_plan_source_consistency CHECK (
        (source = 'MANUAL' AND invoice_plan_id IS NULL)
        OR (source = 'SCHEDULED' AND invoice_plan_id IS NOT NULL)
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

CREATE UNIQUE INDEX IF NOT EXISTS uq_invoices_scheduled_plan_unit_period
ON invoices(invoice_plan_id, unit_id, billing_period)
WHERE source = 'SCHEDULED';

CREATE TABLE IF NOT EXISTS payments (
    -- Payment primary key.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Tenant scope.
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Community scope.
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,

    -- Invoice being paid.
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE RESTRICT,

    -- Amount received.
    amount NUMERIC(12,2) NOT NULL,

    -- Review lifecycle for resident-submitted proof and admin-recorded payments.
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',

    -- Payment method label.
    -- Example: CASH, BANK_TRANSFER, CARD, STRIPE
    payment_method VARCHAR(50),

    -- External transaction/reference ID if present.
    transaction_id VARCHAR(255),

    -- User who submitted the proof or payment record.
    submitted_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,

    -- User who directly recorded or finalized the payment.
    recorded_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Optional admin/resident notes.
    notes TEXT,

    -- Flexible proof metadata for uploads or references.
    evidence JSONB NOT NULL DEFAULT '{}',

    -- When the payment was reviewed.
    reviewed_at TIMESTAMPTZ,

    -- When payment was made/recorded.
    payment_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Row creation time.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Payment amount must be positive.
    CONSTRAINT chk_payments_amount CHECK (amount > 0),

    CONSTRAINT chk_payments_status CHECK (
        status IN ('PENDING', 'APPROVED', 'REJECTED', 'RECORDED')
    )
);

CREATE INDEX IF NOT EXISTS idx_payments_invoice_id
ON payments(invoice_id);

CREATE INDEX IF NOT EXISTS idx_payments_community_status
ON payments(community_id, status);

-- =========================================================
-- 16. COMMUNICATION
-- =========================================================
--
-- announcements = broad community notices/messages
-- notifications = user-targeted message records
--
CREATE TABLE IF NOT EXISTS announcements (
    -- Announcement primary key.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Tenant scope.
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Community scope.
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,

    -- User who authored the announcement.
    author_id UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Headline/title.
    title VARCHAR(255) NOT NULL,

    -- Main message body.
    content TEXT NOT NULL,

    -- Intended audience segment.
    audience_type VARCHAR(50) NOT NULL DEFAULT 'ALL',

    -- Priority flag.
    -- Could be used for pinned/emergency/high-visibility notices.
    priority BOOLEAN NOT NULL DEFAULT FALSE,

    -- Optional expiration time after which notice is no longer relevant.
    expires_at TIMESTAMPTZ,

    -- Creation time.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_announcements_audience_type CHECK (
        audience_type IN ('ALL', 'OWNERS', 'TENANTS', 'SECURITY', 'STAFF')
    )
);

CREATE TABLE IF NOT EXISTS notifications (
    -- Notification primary key.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Tenant scope.
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Community scope.
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,

    -- Target user receiving the notification.
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Notification title.
    title VARCHAR(255) NOT NULL,

    -- Notification body/message text.
    message TEXT NOT NULL,

    -- Simple read/unread state.
    is_read BOOLEAN NOT NULL DEFAULT FALSE,

    -- Optional app route or related resource path.
    -- Example: /invoices/<id>
    link_to_resource VARCHAR(255),

    -- Creation time.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Speed up unread-notification lookups for a user.
CREATE INDEX IF NOT EXISTS idx_notifications_user_unread
ON notifications(user_id)
WHERE is_read = FALSE;

-- =========================================================
-- 17. AUDIT / OUTBOX
-- =========================================================
--
-- audit_logs:
--   append-only records of important actions
--
-- outbox_events:
--   supports reliable async/event-driven integration patterns
--
-- Why outbox_events matters:
--   if you create a DB row and also need to send an email/webhook/queue event,
--   the outbox pattern avoids losing the event if the app crashes mid-process.
--
CREATE TABLE IF NOT EXISTS audit_logs (
    -- Audit row primary key.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Tenant scope.
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- User who performed the action, if known.
    actor_user_id UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Type of entity affected.
    -- Example: "invoice", "resident", "visitor"
    entity_type VARCHAR(100) NOT NULL,

    -- ID of the affected entity.
    entity_id UUID NOT NULL,

    -- Action performed.
    -- Example: CREATED, UPDATED, APPROVED, CANCELLED
    action VARCHAR(100) NOT NULL,

    -- Snapshot of relevant values before the change.
    old_values JSONB NOT NULL DEFAULT '{}',

    -- Snapshot of relevant values after the change.
    new_values JSONB NOT NULL DEFAULT '{}',

    -- When the action happened.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS outbox_events (
    -- Outbox row primary key.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Tenant scope.
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Event classification.
    -- Example: notification.created, payment.received
    event_type VARCHAR(100) NOT NULL,

    -- Business aggregate type involved.
    -- Example: invoice, maintenance_request
    aggregate_type VARCHAR(100) NOT NULL,

    -- ID of the business aggregate.
    aggregate_id UUID NOT NULL,

    -- Event payload/body to be processed by workers/integrations.
    payload JSONB NOT NULL,

    -- Processing state.
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',

    -- Earliest time the event is eligible for processing/retry.
    available_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- When successfully processed.
    processed_at TIMESTAMPTZ,

    -- Creation time.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_outbox_status CHECK (
        status IN ('PENDING', 'PROCESSING', 'PROCESSED', 'FAILED')
    )
);

CREATE INDEX IF NOT EXISTS idx_outbox_events_status_available_at
ON outbox_events(status, available_at);

-- =========================================================
-- 18. RLS EXAMPLE
-- =========================================================
--
-- This is an example of Row-Level Security applied to announcements.
--
-- Important note:
-- For MVP, you can enforce tenant/community scoping in Go repositories first.
-- RLS can then be added gradually as defense in depth.
--
-- app.current_organization_id and app.current_community_id are expected to be
-- set by the application per request/session/transaction.
--
ALTER TABLE announcements ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS organization_community_isolation_policy ON announcements;

CREATE POLICY organization_community_isolation_policy
ON announcements
USING (
    organization_id = current_setting('app.current_organization_id')::uuid
    AND community_id = current_setting('app.current_community_id')::uuid
);
