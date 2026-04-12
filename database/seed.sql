BEGIN;

WITH seeded_organization AS (
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
    VALUES (
        '11111111-1111-1111-1111-111111111111',
        'Palm View Property Management',
        'palm-view-property-management',
        'Palm View Property Management Ltd.',
        'billing@palmview.local',
        '+1-876-555-0100',
        'JM',
        'America/Jamaica',
        'ACTIVE',
        '{"seeded": true, "market": "development"}'::jsonb
    )
    ON CONFLICT (slug) DO UPDATE
    SET
        name = EXCLUDED.name,
        legal_name = EXCLUDED.legal_name,
        billing_email = EXCLUDED.billing_email,
        phone = EXCLUDED.phone,
        country_code = EXCLUDED.country_code,
        timezone = EXCLUDED.timezone,
        status = EXCLUDED.status,
        settings = EXCLUDED.settings,
        updated_at = NOW()
    RETURNING id
),
seeded_subscription AS (
    INSERT INTO subscriptions (
        organization_id,
        plan_code,
        status,
        billing_cycle,
        current_period_start,
        current_period_end,
        metadata
    )
    SELECT
        id,
        'PRO',
        'ACTIVE',
        'MONTHLY',
        NOW(),
        NOW() + INTERVAL '30 days',
        '{"seeded": true}'::jsonb
    FROM seeded_organization
    WHERE NOT EXISTS (
        SELECT 1
        FROM subscriptions s
        WHERE s.organization_id = seeded_organization.id
          AND s.plan_code = 'PRO'
    )
    RETURNING organization_id
),
seeded_community AS (
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
        '22222222-2222-2222-2222-222222222222',
        id,
        'Palm View Estate',
        'palm-view-estate',
        '12 Seaview Road, Kingston, Jamaica',
        'America/Jamaica',
        'ACTIVE',
        '{"seeded": true, "gates": 2}'::jsonb
    FROM seeded_organization
    ON CONFLICT (organization_id, slug) DO UPDATE
    SET
        name = EXCLUDED.name,
        address = EXCLUDED.address,
        timezone = EXCLUDED.timezone,
        status = EXCLUDED.status,
        settings = EXCLUDED.settings,
        updated_at = NOW()
    RETURNING id, organization_id
),
owner_user AS (
    INSERT INTO users (id, email, first_name, last_name, phone, is_active)
    VALUES (
        '33333333-3333-3333-3333-333333333333',
        'owner@palmview.local',
        'Olivia',
        'Owner',
        '+1-876-555-0101',
        TRUE
    )
    ON CONFLICT (email) DO UPDATE
    SET
        first_name = EXCLUDED.first_name,
        last_name = EXCLUDED.last_name,
        phone = EXCLUDED.phone,
        is_active = EXCLUDED.is_active,
        deleted_at = NULL,
        updated_at = NOW()
    RETURNING id, email
),
admin_user AS (
    INSERT INTO users (id, email, first_name, last_name, phone, is_active)
    VALUES (
        '44444444-4444-4444-4444-444444444444',
        'admin@palmview.local',
        'Andre',
        'Admin',
        '+1-876-555-0102',
        TRUE
    )
    ON CONFLICT (email) DO UPDATE
    SET
        first_name = EXCLUDED.first_name,
        last_name = EXCLUDED.last_name,
        phone = EXCLUDED.phone,
        is_active = EXCLUDED.is_active,
        deleted_at = NULL,
        updated_at = NOW()
    RETURNING id, email
),
resident_user AS (
    INSERT INTO users (id, email, first_name, last_name, phone, is_active)
    VALUES (
        '55555555-5555-5555-5555-555555555555',
        'resident@palmview.local',
        'Rhea',
        'Resident',
        '+1-876-555-0103',
        TRUE
    )
    ON CONFLICT (email) DO UPDATE
    SET
        first_name = EXCLUDED.first_name,
        last_name = EXCLUDED.last_name,
        phone = EXCLUDED.phone,
        is_active = EXCLUDED.is_active,
        deleted_at = NULL,
        updated_at = NOW()
    RETURNING id, email
),
security_user AS (
    INSERT INTO users (id, email, first_name, last_name, phone, is_active)
    VALUES (
        '66666666-6666-6666-6666-666666666666',
        'security@palmview.local',
        'Sam',
        'Security',
        '+1-876-555-0104',
        TRUE
    )
    ON CONFLICT (email) DO UPDATE
    SET
        first_name = EXCLUDED.first_name,
        last_name = EXCLUDED.last_name,
        phone = EXCLUDED.phone,
        is_active = EXCLUDED.is_active,
        deleted_at = NULL,
        updated_at = NOW()
    RETURNING id, email
),
auth_accounts_seed AS (
    INSERT INTO auth_accounts (
        user_id,
        provider,
        provider_user_id,
        email_at_provider,
        password_hash,
        last_login_at
    )
    VALUES
        (
            (SELECT id FROM owner_user),
            'password',
            'owner@palmview.local',
            'owner@palmview.local',
            '636f6d756e652d6f776e65722d303031:17d3ede1b40479a8046fa817609d4e12fea374b101bd599369371a84808e36ad',
            NOW()
        ),
        (
            (SELECT id FROM admin_user),
            'password',
            'admin@palmview.local',
            'admin@palmview.local',
            '636f6d756e652d61646d696e2d303031:fdad6394c245952f3c88ae2cd8d76e5258270d82584bc683393687edeceb6f43',
            NOW()
        ),
        (
            (SELECT id FROM resident_user),
            'password',
            'resident@palmview.local',
            'resident@palmview.local',
            '636f6d756e652d7265736964656e7431:6f48ac707196226a39c14ecfc440454f20f0c9eb3a656beb34ee5272fb91bc10',
            NOW()
        ),
        (
            (SELECT id FROM security_user),
            'password',
            'security@palmview.local',
            'security@palmview.local',
            '636f6d756e652d736563757269747931:4eed1bf1980d9cc330907d80a683a1a8343257cf99e9d29e8b45b6b4842650ac',
            NOW()
        )
    ON CONFLICT (provider, provider_user_id) DO UPDATE
    SET
        user_id = EXCLUDED.user_id,
        email_at_provider = EXCLUDED.email_at_provider,
        password_hash = EXCLUDED.password_hash,
        last_login_at = EXCLUDED.last_login_at,
        updated_at = NOW()
    RETURNING user_id
)
INSERT INTO organization_users (
    organization_id,
    user_id,
    role,
    status,
    invited_by,
    joined_at
)
VALUES
    (
        (SELECT id FROM seeded_organization),
        (SELECT id FROM owner_user),
        'OWNER',
        'ACTIVE',
        NULL,
        NOW()
    ),
    (
        (SELECT id FROM seeded_organization),
        (SELECT id FROM admin_user),
        'ORG_ADMIN',
        'ACTIVE',
        (SELECT id FROM owner_user),
        NOW()
    )
ON CONFLICT (organization_id, user_id) DO UPDATE
SET
    role = EXCLUDED.role,
    status = EXCLUDED.status,
    invited_by = EXCLUDED.invited_by,
    joined_at = COALESCE(organization_users.joined_at, EXCLUDED.joined_at),
    updated_at = NOW();

INSERT INTO community_users (
    organization_id,
    community_id,
    user_id,
    role,
    status,
    invited_by,
    joined_at
)
VALUES
    (
        (
            SELECT o.id
            FROM organizations o
            WHERE o.slug = 'palm-view-property-management'
        ),
        (
            SELECT c.id
            FROM communities c
            JOIN organizations o ON o.id = c.organization_id
            WHERE o.slug = 'palm-view-property-management'
              AND c.slug = 'palm-view-estate'
        ),
        (
            SELECT u.id
            FROM users u
            WHERE u.email = 'admin@palmview.local'
        ),
        'COMMUNITY_ADMIN',
        'ACTIVE',
        (
            SELECT u.id
            FROM users u
            WHERE u.email = 'owner@palmview.local'
        ),
        NOW()
    ),
    (
        (
            SELECT o.id
            FROM organizations o
            WHERE o.slug = 'palm-view-property-management'
        ),
        (
            SELECT c.id
            FROM communities c
            JOIN organizations o ON o.id = c.organization_id
            WHERE o.slug = 'palm-view-property-management'
              AND c.slug = 'palm-view-estate'
        ),
        (
            SELECT u.id
            FROM users u
            WHERE u.email = 'resident@palmview.local'
        ),
        'RESIDENT',
        'ACTIVE',
        (
            SELECT u.id
            FROM users u
            WHERE u.email = 'admin@palmview.local'
        ),
        NOW()
    ),
    (
        (
            SELECT o.id
            FROM organizations o
            WHERE o.slug = 'palm-view-property-management'
        ),
        (
            SELECT c.id
            FROM communities c
            JOIN organizations o ON o.id = c.organization_id
            WHERE o.slug = 'palm-view-property-management'
              AND c.slug = 'palm-view-estate'
        ),
        (
            SELECT u.id
            FROM users u
            WHERE u.email = 'security@palmview.local'
        ),
        'SECURITY',
        'ACTIVE',
        (
            SELECT u.id
            FROM users u
            WHERE u.email = 'admin@palmview.local'
        ),
        NOW()
    )
ON CONFLICT (community_id, user_id) DO UPDATE
SET
    role = EXCLUDED.role,
    status = EXCLUDED.status,
    invited_by = EXCLUDED.invited_by,
    joined_at = COALESCE(community_users.joined_at, EXCLUDED.joined_at),
    updated_at = NOW();

COMMIT;
