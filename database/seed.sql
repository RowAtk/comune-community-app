BEGIN;

SET LOCAL row_security = off;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE OR REPLACE FUNCTION pg_temp.seed_uuid(seed TEXT)
RETURNS UUID AS $$
    SELECT (
        SUBSTRING(md5(seed), 1, 8) || '-' ||
        SUBSTRING(md5(seed), 9, 4) || '-' ||
        SUBSTRING(md5(seed), 13, 4) || '-' ||
        SUBSTRING(md5(seed), 17, 4) || '-' ||
        SUBSTRING(md5(seed), 21, 12)
    )::uuid;
$$ LANGUAGE SQL IMMUTABLE;

TRUNCATE TABLE
    outbox_events,
    audit_logs,
    notifications,
    announcements,
    payments,
    invoices,
    invoice_plan_unit_overrides,
    invoice_plans,
    maintenance_requests,
    security_logs,
    visitors,
    vendors,
    resident_invitations,
    residents,
    households,
    units,
    community_invitations,
    organization_invitations,
    community_users,
    organization_users,
    auth_accounts,
    users,
    communities,
    subscriptions,
    organizations
RESTART IDENTITY CASCADE;

CREATE TEMP TABLE seed_orgs AS
SELECT
    seq,
    pg_temp.seed_uuid('org-' || seq) AS organization_id,
    pg_temp.seed_uuid('community-' || seq) AS community_id,
    org_name,
    org_slug,
    legal_name,
    domain,
    community_name,
    community_slug,
    address,
    org_status,
    community_status,
    gates
FROM (
    VALUES
        (1, 'Palm View Property Management', 'palm-view-property-management', 'Palm View Property Management Ltd.', 'palmview.local', 'Palm View Estate', 'palm-view-estate', '12 Seaview Road, Kingston, Jamaica', 'ACTIVE', 'ACTIVE', 2),
        (2, 'Harbour Crest Management', 'harbour-crest-management', 'Harbour Crest Management Limited', 'harbourcrest.local', 'Harbour Crest Residences', 'harbour-crest-residences', '4 Ocean Drive, Montego Bay, Jamaica', 'ACTIVE', 'ACTIVE', 3),
        (3, 'Coral Gardens Property Services', 'coral-gardens-property-services', 'Coral Gardens Property Services Ltd.', 'coralgardens.local', 'Coral Gardens', 'coral-gardens', '27 Coral Way, Ocho Rios, Jamaica', 'ACTIVE', 'ACTIVE', 1),
        (4, 'Blue Mountain Communities', 'blue-mountain-communities', 'Blue Mountain Communities Ltd.', 'bluemountain.local', 'Blue Mountain Villas', 'blue-mountain-villas', '9 Hillcrest Avenue, Kingston, Jamaica', 'ACTIVE', 'ACTIVE', 2),
        (5, 'Seaview Heights Management', 'seaview-heights-management', 'Seaview Heights Management Ltd.', 'seaviewheights.local', 'Seaview Heights', 'seaview-heights', '18 Fortlands Road, Discovery Bay, Jamaica', 'ACTIVE', 'ACTIVE', 2),
        (6, 'Cedar Grove Operations', 'cedar-grove-operations', 'Cedar Grove Operations Limited', 'cedargrove.local', 'Cedar Grove Court', 'cedar-grove-court', '31 Fairfield Lane, Mandeville, Jamaica', 'SUSPENDED', 'INACTIVE', 1),
        (7, 'Mango Walk Estates', 'mango-walk-estates', 'Mango Walk Estates Limited', 'mangowalk.local', 'Mango Walk Townhomes', 'mango-walk-townhomes', '7 Mango Walk Boulevard, St. James, Jamaica', 'ACTIVE', 'ACTIVE', 2),
        (8, 'Kingsway Terraces Group', 'kingsway-terraces-group', 'Kingsway Terraces Group Ltd.', 'kingswayterraces.local', 'Kingsway Terraces', 'kingsway-terraces', '55 Kingsway Terrace, Kingston, Jamaica', 'ACTIVE', 'ACTIVE', 3),
        (9, 'Liguanea Court Management', 'liguanea-court-management', 'Liguanea Court Management Ltd.', 'liguaneacourt.local', 'Liguanea Court', 'liguanea-court', '22 Trafalgar Crescent, Kingston, Jamaica', 'ACTIVE', 'ACTIVE', 2),
        (10, 'Portside Villas Services', 'portside-villas-services', 'Portside Villas Services Ltd.', 'portsidevillas.local', 'Portside Villas', 'portside-villas', '3 Harbour Street, Portmore, Jamaica', 'INACTIVE', 'ARCHIVED', 1),
        (11, 'Fern Ridge Management', 'fern-ridge-management', 'Fern Ridge Management Limited', 'fernridge.local', 'Fern Ridge Apartments', 'fern-ridge-apartments', '14 Orange Grove Road, St. Ann, Jamaica', 'ACTIVE', 'ACTIVE', 2),
        (12, 'Hopewell Gardens Management', 'hopewell-gardens-management', 'Hopewell Gardens Management Ltd.', 'hopewellgardens.local', 'Hopewell Gardens', 'hopewell-gardens', '8 Main Street, Hanover, Jamaica', 'ACTIVE', 'ACTIVE', 1),
        (13, 'Sunrise Meadows Property Group', 'sunrise-meadows-property-group', 'Sunrise Meadows Property Group Ltd.', 'sunrisemeadows.local', 'Sunrise Meadows', 'sunrise-meadows', '44 Cumberland Road, Spanish Town, Jamaica', 'ACTIVE', 'ACTIVE', 2),
        (14, 'Orchid Park Communities', 'orchid-park-communities', 'Orchid Park Communities Ltd.', 'orchidpark.local', 'Orchid Park Residences', 'orchid-park-residences', '61 Old Hope Road, Kingston, Jamaica', 'ACTIVE', 'ACTIVE', 2),
        (15, 'Silver Sands Community Services', 'silver-sands-community-services', 'Silver Sands Community Services Ltd.', 'silversands.local', 'Silver Sands', 'silver-sands', '6 Beachfront Avenue, Trelawny, Jamaica', 'ACTIVE', 'ACTIVE', 4),
        (16, 'Fairfield Commons Management', 'fairfield-commons-management', 'Fairfield Commons Management Ltd.', 'fairfieldcommons.local', 'Fairfield Commons', 'fairfield-commons', '28 Fairfield Road, Montego Bay, Jamaica', 'ACTIVE', 'ACTIVE', 2),
        (17, 'Tamarind Estate Management', 'tamarind-estate-management', 'Tamarind Estate Management Ltd.', 'tamarindestate.local', 'Tamarind Estate', 'tamarind-estate', '77 Brunswick Avenue, Clarendon, Jamaica', 'ACTIVE', 'ACTIVE', 2),
        (18, 'West Bay Residences Group', 'west-bay-residences-group', 'West Bay Residences Group Ltd.', 'westbay.local', 'West Bay Residences', 'west-bay-residences', '19 Beach Road, Negril, Jamaica', 'SUSPENDED', 'INACTIVE', 1),
        (19, 'Green Acres Management', 'green-acres-management', 'Green Acres Management Ltd.', 'greenacres.local', 'Green Acres Village', 'green-acres-village', '11 Green Acres Drive, May Pen, Jamaica', 'ACTIVE', 'ACTIVE', 2),
        (20, 'Horizon Pointe Services', 'horizon-pointe-services', 'Horizon Pointe Services Ltd.', 'horizonpointe.local', 'Horizon Pointe', 'horizon-pointe', '39 Horizon Close, Port Antonio, Jamaica', 'ACTIVE', 'ACTIVE', 3)
) AS t(
    seq,
    org_name,
    org_slug,
    legal_name,
    domain,
    community_name,
    community_slug,
    address,
    org_status,
    community_status,
    gates
);

CREATE TEMP TABLE seed_users AS
WITH roles AS (
    SELECT *
    FROM (
        VALUES
            ('owner', 'OWNER', 'BOARD_MEMBER', 'owner', 1),
            ('admin', 'ORG_ADMIN', 'COMMUNITY_ADMIN', 'admin', 2),
            ('resident', NULL, 'RESIDENT', 'resident', 3),
            ('security', NULL, 'SECURITY', 'security', 4),
            ('operations', 'OPERATIONS', 'MANAGER', 'ops', 5)
    ) AS r(kind, organization_role, community_role, email_localpart, offset_index)
)
SELECT
    o.seq AS organization_seq,
    o.organization_id,
    o.community_id,
    r.kind,
    r.organization_role,
    r.community_role,
    r.offset_index,
    pg_temp.seed_uuid('user-' || o.seq || '-' || r.kind) AS user_id,
    r.email_localpart || '@' || o.domain AS email,
    (ARRAY['Olivia', 'Andre', 'Rhea', 'Samuel', 'Nadia', 'Dwayne', 'Talia', 'Kemar', 'Janelle', 'Marcus', 'Alicia', 'Damian', 'Shanice', 'Trevor', 'Monique', 'Ricardo', 'Kayla', 'Jerome', 'Arielle', 'Nathan'])[((o.seq + r.offset_index - 2) % 20) + 1] AS first_name,
    (ARRAY['Bennett', 'Campbell', 'Morgan', 'Brown', 'Williams', 'Johnson', 'Clarke', 'Reid', 'Thompson', 'Edwards', 'Grant', 'Walker', 'McKenzie', 'Powell', 'Allen', 'Robinson', 'Stewart', 'Davis', 'Bailey', 'Lewis'])[((o.seq + r.offset_index * 2 - 2) % 20) + 1] AS last_name,
    '+1-876-555-' || LPAD((1000 + (o.seq * 10) + r.offset_index)::TEXT, 4, '0') AS phone,
    NOT (r.kind = 'security' AND o.seq IN (6, 18)) AS is_active
FROM seed_orgs o
CROSS JOIN roles r;

CREATE TEMP TABLE seed_units AS
SELECT
    o.seq AS organization_seq,
    o.organization_id,
    o.community_id,
    pos AS unit_position,
    pg_temp.seed_uuid('unit-' || o.seq || '-' || pos) AS unit_id,
    unit_number,
    block_floor,
    unit_type,
    CASE
        WHEN pos = 2 AND o.seq IN (6, 10, 18) THEN 'INACTIVE'
        ELSE 'ACTIVE'
    END AS status
FROM seed_orgs o
CROSS JOIN (
    VALUES
        (1, 'A-101', 'Block A / Ground', 'TOWNHOUSE'),
        (2, 'B-202', 'Block B / Second Floor', 'APARTMENT')
) AS u(pos, unit_number, block_floor, unit_type);

CREATE TEMP TABLE seed_households AS
SELECT
    u.organization_seq,
    u.organization_id,
    u.community_id,
    u.unit_id,
    u.unit_position,
    pg_temp.seed_uuid('household-' || u.organization_seq || '-' || u.unit_position) AS household_id,
    'The ' ||
    (ARRAY['Bennett', 'Campbell', 'Morgan', 'Brown', 'Williams', 'Johnson', 'Clarke', 'Reid', 'Thompson', 'Edwards', 'Grant', 'Walker', 'McKenzie', 'Powell', 'Allen', 'Robinson', 'Stewart', 'Davis', 'Bailey', 'Lewis'])[((u.organization_seq + u.unit_position - 2) % 20) + 1] ||
    ' Household' AS household_name
FROM seed_units u;

CREATE TEMP TABLE seed_residents AS
SELECT
    h.organization_seq,
    h.organization_id,
    h.community_id,
    h.unit_id,
    h.household_id,
    h.unit_position,
    resident_position,
    pg_temp.seed_uuid('resident-' || h.organization_seq || '-' || h.unit_position || '-' || resident_position) AS resident_id,
    CASE
        WHEN h.unit_position = 1 AND resident_position = 1
            THEN pg_temp.seed_uuid('user-' || h.organization_seq || '-resident')
        ELSE NULL
    END AS user_id,
    (ARRAY['Adrian', 'Bianca', 'Christopher', 'Danielle', 'Ethan', 'Farah', 'Gavin', 'Hannah', 'Isaac', 'Jada', 'Kyle', 'Latoya', 'Malik', 'Naomi', 'Owen', 'Priya', 'Quinton', 'Rochelle', 'Stefan', 'Tiana'])[((h.organization_seq + h.unit_position + resident_position - 3) % 20) + 1] AS first_name,
    (ARRAY['Bennett', 'Campbell', 'Morgan', 'Brown', 'Williams', 'Johnson', 'Clarke', 'Reid', 'Thompson', 'Edwards', 'Grant', 'Walker', 'McKenzie', 'Powell', 'Allen', 'Robinson', 'Stewart', 'Davis', 'Bailey', 'Lewis'])[((h.organization_seq + h.unit_position - 2) % 20) + 1] AS last_name,
    LOWER(
        (ARRAY['adrian', 'bianca', 'chris', 'danielle', 'ethan', 'farah', 'gavin', 'hannah', 'isaac', 'jada', 'kyle', 'latoya', 'malik', 'naomi', 'owen', 'priya', 'quinton', 'rochelle', 'stefan', 'tiana'])[((h.organization_seq + h.unit_position + resident_position - 3) % 20) + 1] ||
        '.' ||
        (ARRAY['bennett', 'campbell', 'morgan', 'brown', 'williams', 'johnson', 'clarke', 'reid', 'thompson', 'edwards', 'grant', 'walker', 'mckenzie', 'powell', 'allen', 'robinson', 'stewart', 'davis', 'bailey', 'lewis'])[((h.organization_seq + h.unit_position - 2) % 20) + 1] ||
        '+' || h.organization_seq || h.unit_position || resident_position || '@' ||
        (SELECT domain FROM seed_orgs o WHERE o.seq = h.organization_seq)
    ) AS email,
    '+1-876-555-' || LPAD((3000 + (h.organization_seq * 10) + (h.unit_position * 2) + resident_position)::TEXT, 4, '0') AS phone,
    CASE
        WHEN h.unit_position = 1 AND resident_position = 1 THEN 'OWNER'
        WHEN h.unit_position = 1 AND resident_position = 2 THEN 'DEPENDENT'
        WHEN h.unit_position = 2 AND resident_position = 1 THEN 'TENANT'
        ELSE 'OCCUPANT'
    END AS resident_type,
    CASE
        WHEN resident_position = 1 THEN 'HOUSEHOLD_ADMIN'
        WHEN h.unit_position = 1 THEN 'HOUSEHOLD_MEMBER'
        ELSE 'HOUSEHOLD_VIEWER'
    END AS household_role,
    CASE
        WHEN h.unit_position = 2 AND h.organization_seq IN (5, 10, 15, 20) AND resident_position = 1 THEN 'MOVED_OUT'
        WHEN resident_position = 2 AND h.organization_seq IN (3, 8, 13, 18) THEN 'PENDING'
        WHEN h.unit_position = 2 AND h.organization_seq IN (6, 12, 18) THEN 'INACTIVE'
        ELSE 'ACTIVE'
    END AS status,
    resident_position = 1 AS is_primary_contact,
    DATE '2023-01-01' + ((h.organization_seq * 11) + (h.unit_position * 3) + resident_position) AS move_in_date,
    CASE
        WHEN h.unit_position = 2 AND h.organization_seq IN (5, 10, 15, 20) AND resident_position = 1
            THEN DATE '2025-12-01' + h.organization_seq
        ELSE NULL
    END AS move_out_date
FROM seed_households h
CROSS JOIN (VALUES (1), (2)) AS rp(resident_position);

CREATE TEMP TABLE seed_invoice_plans AS
SELECT
    o.seq AS organization_seq,
    o.organization_id,
    o.community_id,
    pg_temp.seed_uuid('invoice-plan-' || o.seq) AS invoice_plan_id,
    'Monthly Maintenance ' || o.community_name AS plan_name,
    CASE
        WHEN o.seq IN (6, 18) THEN 'PAUSED'
        WHEN o.seq IN (10, 15) THEN 'ENDED'
        ELSE 'ACTIVE'
    END AS status,
    ((o.seq - 1) % 4) + 1 AS issue_day_of_month,
    8 + ((o.seq - 1) % 7) AS due_day_of_month,
    (7800 + (o.seq * 225))::NUMERIC(12,2) AS default_amount,
    DATE '2026-01-01' + (o.seq - 1) AS starts_on,
    CASE
        WHEN o.seq IN (10, 15) THEN DATE '2026-04-30'
        ELSE NULL
    END AS ends_on,
    'Recurring maintenance dues for ' || o.community_name AS description,
    pg_temp.seed_uuid('user-' || o.seq || '-admin') AS created_by
FROM seed_orgs o;

CREATE TEMP TABLE seed_invoice_plan_unit_overrides AS
SELECT
    p.organization_seq,
    p.organization_id,
    p.community_id,
    p.invoice_plan_id,
    u.unit_id,
    (p.default_amount + CASE WHEN p.organization_seq % 2 = 0 THEN 1500 ELSE 900 END)::NUMERIC(12,2) AS amount
FROM seed_invoice_plans p
JOIN seed_units u
  ON u.organization_seq = p.organization_seq
 AND u.unit_position = 2;

CREATE TEMP TABLE seed_invoices AS
SELECT
    pg_temp.seed_uuid('invoice-apr-' || u.organization_seq || '-' || u.unit_position) AS invoice_id,
    u.organization_seq,
    u.organization_id,
    u.community_id,
    u.unit_id,
    p.invoice_plan_id,
    'SCHEDULED' AS source,
    COALESCE(o.amount, p.default_amount) AS amount,
    CASE
        WHEN u.unit_position = 1 AND u.organization_seq IN (1, 5, 10, 15, 20) THEN COALESCE(o.amount, p.default_amount)
        WHEN u.organization_seq % 3 = 0 THEN ROUND(COALESCE(o.amount, p.default_amount) * 0.4, 2)
        ELSE 0.00
    END AS paid_amount,
    DATE '2026-04-01' + (p.issue_day_of_month - 1) AS issued_on,
    DATE '2026-04-01' + (p.due_day_of_month - 1) AS due_date,
    CASE
        WHEN u.unit_position = 1 AND u.organization_seq IN (1, 5, 10, 15, 20) THEN 'PAID'
        WHEN u.organization_seq % 3 = 0 THEN 'PARTIAL'
        ELSE 'OVERDUE'
    END AS status,
    '2026-04' AS billing_period,
    'April 2026 maintenance dues' AS description
FROM seed_units u
JOIN seed_invoice_plans p
  ON p.organization_seq = u.organization_seq
LEFT JOIN seed_invoice_plan_unit_overrides o
  ON o.invoice_plan_id = p.invoice_plan_id
 AND o.unit_id = u.unit_id
UNION ALL
SELECT
    pg_temp.seed_uuid('invoice-may-' || u.organization_seq || '-' || u.unit_position) AS invoice_id,
    u.organization_seq,
    u.organization_id,
    u.community_id,
    u.unit_id,
    p.invoice_plan_id,
    'SCHEDULED' AS source,
    COALESCE(o.amount, p.default_amount) AS amount,
    CASE
        WHEN u.organization_seq IN (4, 8, 12, 16, 20) THEN ROUND(COALESCE(o.amount, p.default_amount) * 0.35, 2)
        ELSE 0.00
    END AS paid_amount,
    DATE '2026-05-01' + (p.issue_day_of_month - 1) AS issued_on,
    DATE '2026-05-01' + (p.due_day_of_month - 1) AS due_date,
    CASE
        WHEN u.organization_seq IN (4, 8, 12, 16, 20) THEN 'PARTIAL'
        ELSE 'UNPAID'
    END AS status,
    '2026-05' AS billing_period,
    'May 2026 maintenance dues' AS description
FROM seed_units u
JOIN seed_invoice_plans p
  ON p.organization_seq = u.organization_seq
LEFT JOIN seed_invoice_plan_unit_overrides o
  ON o.invoice_plan_id = p.invoice_plan_id
 AND o.unit_id = u.unit_id
WHERE u.unit_position = 1
  AND p.status = 'ACTIVE'
UNION ALL
SELECT
    pg_temp.seed_uuid('invoice-manual-' || u.organization_seq) AS invoice_id,
    u.organization_seq,
    u.organization_id,
    u.community_id,
    u.unit_id,
    NULL::UUID AS invoice_plan_id,
    'MANUAL' AS source,
    (2500 + (u.organization_seq * 120))::NUMERIC(12,2) AS amount,
    CASE
        WHEN u.organization_seq IN (8, 16) THEN ROUND((2500 + (u.organization_seq * 120)) * 0.5, 2)
        ELSE 0.00
    END AS paid_amount,
    DATE '2026-05-05' AS issued_on,
    CASE
        WHEN u.organization_seq IN (4, 12, 20) THEN DATE '2026-04-28'
        ELSE DATE '2026-05-20'
    END AS due_date,
    CASE
        WHEN u.organization_seq IN (8, 16) THEN 'PARTIAL'
        WHEN u.organization_seq IN (4, 12, 20) THEN 'OVERDUE'
        ELSE 'UNPAID'
    END AS status,
    CASE
        WHEN u.organization_seq IN (4, 12, 20) THEN '2026-04'
        ELSE '2026-05'
    END AS billing_period,
    'One-time grounds recovery levy' AS description
FROM seed_units u
WHERE u.unit_position = 2
  AND u.organization_seq % 2 = 0;

INSERT INTO organizations (
    id,
    name,
    slug,
    legal_name,
    billing_email,
    phone,
    country_code,
    timezone,
    status,
    settings
)
SELECT
    organization_id,
    org_name,
    org_slug,
    legal_name,
    'billing@' || domain,
    '+1-876-555-' || LPAD((7000 + seq)::TEXT, 4, '0'),
    'JM',
    'America/Jamaica',
    org_status,
    jsonb_build_object('seeded', true, 'org_index', seq, 'portfolio_gates', gates)
FROM seed_orgs;

INSERT INTO communities (
    id,
    organization_id,
    name,
    slug,
    address,
    timezone,
    status,
    settings
)
SELECT
    community_id,
    organization_id,
    community_name,
    community_slug,
    address,
    'America/Jamaica',
    community_status,
    jsonb_build_object('seeded', true, 'gates', gates, 'parking_zones', CASE WHEN seq % 2 = 0 THEN 2 ELSE 1 END)
FROM seed_orgs;

INSERT INTO users (
    id,
    email,
    first_name,
    last_name,
    phone,
    is_active
)
SELECT
    user_id,
    email,
    first_name,
    last_name,
    phone,
    is_active
FROM seed_users;

INSERT INTO auth_accounts (
    id,
    user_id,
    provider,
    provider_user_id,
    email_at_provider,
    password_hash,
    last_login_at
)
SELECT
    pg_temp.seed_uuid('auth-account-' || organization_seq || '-' || kind),
    user_id,
    'password',
    email,
    email,
    '636f6d756e652d6f776e65722d303031:17d3ede1b40479a8046fa817609d4e12fea374b101bd599369371a84808e36ad',
    TIMESTAMPTZ '2026-05-01 09:00:00-05' - MAKE_INTERVAL(days => organization_seq + offset_index)
FROM seed_users;

INSERT INTO organization_users (
    organization_id,
    user_id,
    role,
    status,
    invited_by,
    joined_at
)
SELECT
    organization_id,
    user_id,
    organization_role,
    CASE
        WHEN organization_seq IN (6, 18) AND kind = 'operations' THEN 'SUSPENDED'
        ELSE 'ACTIVE'
    END,
    CASE
        WHEN kind = 'owner' THEN NULL
        ELSE pg_temp.seed_uuid('user-' || organization_seq || '-owner')
    END,
    TIMESTAMPTZ '2024-01-15 10:00:00-05' + MAKE_INTERVAL(days => organization_seq * 5 + offset_index)
FROM seed_users
WHERE organization_role IS NOT NULL;

INSERT INTO community_users (
    organization_id,
    community_id,
    user_id,
    role,
    status,
    invited_by,
    joined_at
)
SELECT
    organization_id,
    community_id,
    user_id,
    community_role,
    CASE
        WHEN kind = 'resident' AND organization_seq IN (3, 8, 13, 18) THEN 'INVITED'
        WHEN kind = 'security' AND organization_seq IN (6, 18) THEN 'SUSPENDED'
        ELSE 'ACTIVE'
    END,
    CASE
        WHEN kind = 'owner' THEN NULL
        ELSE pg_temp.seed_uuid('user-' || organization_seq || '-admin')
    END,
    TIMESTAMPTZ '2024-02-01 09:00:00-05' + MAKE_INTERVAL(days => organization_seq * 4 + offset_index)
FROM seed_users;

INSERT INTO organization_invitations (
    id,
    organization_id,
    email,
    role,
    token,
    expires_at,
    accepted_at,
    invited_by
)
SELECT
    pg_temp.seed_uuid('org-invite-' || seq),
    organization_id,
    'billing+' || LPAD(seq::TEXT, 2, '0') || '@' || domain,
    CASE WHEN seq % 2 = 0 THEN 'BILLING_ADMIN' ELSE 'SUPPORT' END,
    pg_temp.seed_uuid('org-invite-token-' || seq),
    TIMESTAMPTZ '2026-06-15 12:00:00-05' + MAKE_INTERVAL(days => seq),
    CASE WHEN seq IN (1, 4, 7, 10, 13, 16, 19) THEN TIMESTAMPTZ '2026-05-20 09:00:00-05' + MAKE_INTERVAL(days => seq) ELSE NULL END,
    pg_temp.seed_uuid('user-' || seq || '-owner')
FROM seed_orgs;

INSERT INTO community_invitations (
    id,
    organization_id,
    community_id,
    email,
    role,
    token,
    expires_at,
    accepted_at,
    invited_by
)
SELECT
    pg_temp.seed_uuid('community-invite-' || seq),
    organization_id,
    community_id,
    'board+' || LPAD(seq::TEXT, 2, '0') || '@' || domain,
    CASE WHEN seq % 2 = 0 THEN 'BOARD_MEMBER' ELSE 'MAINTENANCE_STAFF' END,
    pg_temp.seed_uuid('community-invite-token-' || seq),
    TIMESTAMPTZ '2026-06-20 12:00:00-05' + MAKE_INTERVAL(days => seq),
    CASE WHEN seq IN (2, 5, 8, 11, 14, 17, 20) THEN TIMESTAMPTZ '2026-05-22 15:00:00-05' + MAKE_INTERVAL(days => seq) ELSE NULL END,
    pg_temp.seed_uuid('user-' || seq || '-admin')
FROM seed_orgs;

INSERT INTO units (
    id,
    organization_id,
    community_id,
    unit_number,
    block_floor,
    unit_type,
    status
)
SELECT
    unit_id,
    organization_id,
    community_id,
    unit_number,
    block_floor,
    unit_type,
    status
FROM seed_units;

INSERT INTO households (
    id,
    organization_id,
    community_id,
    unit_id,
    name
)
SELECT
    household_id,
    organization_id,
    community_id,
    unit_id,
    household_name
FROM seed_households;

INSERT INTO residents (
    id,
    organization_id,
    community_id,
    unit_id,
    user_id,
    household_id,
    first_name,
    last_name,
    email,
    phone,
    resident_type,
    household_role,
    status,
    is_primary_contact,
    move_in_date,
    move_out_date
)
SELECT
    resident_id,
    organization_id,
    community_id,
    unit_id,
    user_id,
    household_id,
    first_name,
    last_name,
    email,
    phone,
    resident_type,
    household_role,
    status,
    is_primary_contact,
    move_in_date,
    move_out_date
FROM seed_residents;

INSERT INTO resident_invitations (
    id,
    organization_id,
    community_id,
    resident_id,
    token,
    expires_at,
    accepted_at,
    invited_by,
    accepted_by_user_id
)
SELECT
    pg_temp.seed_uuid('resident-invite-' || organization_seq),
    organization_id,
    community_id,
    resident_id,
    pg_temp.seed_uuid('resident-invite-token-' || organization_seq),
    TIMESTAMPTZ '2026-06-30 18:00:00-05' + MAKE_INTERVAL(days => organization_seq),
    CASE WHEN organization_seq IN (1, 3, 5, 7, 9, 11, 13, 15, 17, 19) THEN TIMESTAMPTZ '2026-05-18 14:00:00-05' + MAKE_INTERVAL(days => organization_seq) ELSE NULL END,
    pg_temp.seed_uuid('user-' || organization_seq || '-admin'),
    CASE WHEN organization_seq IN (1, 3, 5, 7, 9, 11, 13, 15, 17, 19) THEN user_id ELSE NULL END
FROM seed_residents
WHERE unit_position = 1
  AND resident_position = 1
  AND user_id IS NOT NULL;

INSERT INTO invoice_plans (
    id,
    organization_id,
    community_id,
    name,
    plan_type,
    status,
    issue_day_of_month,
    due_day_of_month,
    default_amount,
    starts_on,
    ends_on,
    description,
    created_by
)
SELECT
    invoice_plan_id,
    organization_id,
    community_id,
    plan_name,
    'MAINTENANCE',
    status,
    issue_day_of_month,
    due_day_of_month,
    default_amount,
    starts_on,
    ends_on,
    description,
    created_by
FROM seed_invoice_plans;

INSERT INTO invoice_plan_unit_overrides (
    invoice_plan_id,
    organization_id,
    community_id,
    unit_id,
    amount
)
SELECT
    invoice_plan_id,
    organization_id,
    community_id,
    unit_id,
    amount
FROM seed_invoice_plan_unit_overrides;

INSERT INTO invoices (
    id,
    organization_id,
    community_id,
    unit_id,
    invoice_plan_id,
    source,
    amount,
    paid_amount,
    issued_on,
    due_date,
    status,
    billing_period,
    description
)
SELECT
    invoice_id,
    organization_id,
    community_id,
    unit_id,
    invoice_plan_id,
    source,
    amount,
    paid_amount,
    issued_on,
    due_date,
    status,
    billing_period,
    description
FROM seed_invoices;

COMMIT;
