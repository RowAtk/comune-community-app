package organizations

import (
	"comune/apps/api/internal/modules/users"
	"comune/apps/api/internal/platform/db"
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func insertOrganization(ctx context.Context, pool *pgxpool.Pool, creatorUserID string, input CreateOrganizationInput) (Organization, error) {
	const query = `
		INSERT INTO organizations (
			name, slug, legal_name, billing_email, phone, country_code, timezone, status, settings
		)
		VALUES ($1, $2, NULLIF($3, ''), NULLIF($4, ''), NULLIF($5, ''), $6, $7, $8, $9::jsonb)
		RETURNING id, name, slug, COALESCE(legal_name, ''), COALESCE(billing_email, ''), COALESCE(phone, ''), country_code, timezone, status, settings, created_at, updated_at, deleted_at
	`

	var org Organization
	err := db.RunInTx(ctx, pool, func(tx pgx.Tx) error {
		if err := tx.QueryRow(
			ctx,
			query,
			input.Name,
			input.Slug,
			input.LegalName,
			input.BillingEmail,
			input.Phone,
			input.CountryCode,
			input.Timezone,
			input.Status,
			jsonbOrEmpty(input.Settings),
		).Scan(
			&org.ID,
			&org.Name,
			&org.Slug,
			&org.LegalName,
			&org.BillingEmail,
			&org.Phone,
			&org.CountryCode,
			&org.Timezone,
			&org.Status,
			&org.Settings,
			&org.CreatedAt,
			&org.UpdatedAt,
			&org.DeletedAt,
		); err != nil {
			return err
		}

		return insertOrganizationMembership(ctx, tx, org.ID, creatorUserID, "OWNER")
	})

	return org, err
}

func listOrganizations(ctx context.Context, db *pgxpool.Pool) ([]Organization, error) {
	const query = `
		SELECT id, name, slug, COALESCE(legal_name, ''), COALESCE(billing_email, ''), COALESCE(phone, ''), country_code, timezone, status, settings, created_at, updated_at, deleted_at
		FROM organizations
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var organizations []Organization
	for rows.Next() {
		var org Organization
		if err := rows.Scan(
			&org.ID,
			&org.Name,
			&org.Slug,
			&org.LegalName,
			&org.BillingEmail,
			&org.Phone,
			&org.CountryCode,
			&org.Timezone,
			&org.Status,
			&org.Settings,
			&org.CreatedAt,
			&org.UpdatedAt,
			&org.DeletedAt,
		); err != nil {
			return nil, err
		}

		organizations = append(organizations, org)
	}

	return organizations, rows.Err()
}

func findOrganizationByID(ctx context.Context, db *pgxpool.Pool, id string) (Organization, error) {
	const query = `
		SELECT id, name, slug, COALESCE(legal_name, ''), COALESCE(billing_email, ''), COALESCE(phone, ''), country_code, timezone, status, settings, created_at, updated_at, deleted_at
		FROM organizations
		WHERE id = $1 AND deleted_at IS NULL
	`

	var org Organization
	err := db.QueryRow(ctx, query, id).Scan(
		&org.ID,
		&org.Name,
		&org.Slug,
		&org.LegalName,
		&org.BillingEmail,
		&org.Phone,
		&org.CountryCode,
		&org.Timezone,
		&org.Status,
		&org.Settings,
		&org.CreatedAt,
		&org.UpdatedAt,
		&org.DeletedAt,
	)

	return org, err
}

func updateOrganization(ctx context.Context, db *pgxpool.Pool, id string, input UpdateOrganizationInput) (Organization, error) {
	const query = `
		UPDATE organizations
		SET
			name = COALESCE($2, name),
			slug = COALESCE($3, slug),
			legal_name = CASE WHEN $4 IS NULL THEN legal_name ELSE NULLIF($4, '') END,
			billing_email = CASE WHEN $5 IS NULL THEN billing_email ELSE NULLIF($5, '') END,
			phone = CASE WHEN $6 IS NULL THEN phone ELSE NULLIF($6, '') END,
			country_code = COALESCE($7, country_code),
			timezone = COALESCE($8, timezone),
			status = COALESCE($9, status),
			settings = COALESCE($10::jsonb, settings)
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id, name, slug, COALESCE(legal_name, ''), COALESCE(billing_email, ''), COALESCE(phone, ''), country_code, timezone, status, settings, created_at, updated_at, deleted_at
	`

	var (
		org      Organization
		settings any
	)

	if input.Settings != nil {
		settings = jsonbOrEmpty(*input.Settings)
	}

	err := db.QueryRow(
		ctx,
		query,
		id,
		input.Name,
		input.Slug,
		input.LegalName,
		input.BillingEmail,
		input.Phone,
		input.CountryCode,
		input.Timezone,
		input.Status,
		settings,
	).Scan(
		&org.ID,
		&org.Name,
		&org.Slug,
		&org.LegalName,
		&org.BillingEmail,
		&org.Phone,
		&org.CountryCode,
		&org.Timezone,
		&org.Status,
		&org.Settings,
		&org.CreatedAt,
		&org.UpdatedAt,
		&org.DeletedAt,
	)

	return org, err
}

func deleteOrganization(ctx context.Context, db *pgxpool.Pool, id string) error {
	const query = `
		UPDATE organizations
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	tag, err := db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func listOrganizationMembers(ctx context.Context, db *pgxpool.Pool, organizationID string) ([]OrganizationMember, error) {
	const query = `
		SELECT
			ou.organization_id,
			ou.user_id,
			ou.role,
			ou.status,
			ou.invited_by,
			ou.joined_at,
			ou.created_at,
			ou.updated_at,
			u.id,
			u.email,
			COALESCE(u.first_name, ''),
			COALESCE(u.last_name, ''),
			COALESCE(u.phone, ''),
			u.is_active,
			u.created_at,
			u.updated_at,
			u.deleted_at
		FROM organization_users ou
		JOIN users u ON u.id = ou.user_id
		WHERE ou.organization_id = $1
		ORDER BY ou.created_at DESC
	`

	rows, err := db.Query(ctx, query, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []OrganizationMember
	for rows.Next() {
		var member OrganizationMember
		if err := rows.Scan(
			&member.OrganizationID,
			&member.UserID,
			&member.Role,
			&member.Status,
			&member.InvitedBy,
			&member.JoinedAt,
			&member.CreatedAt,
			&member.UpdatedAt,
			&member.User.ID,
			&member.User.Email,
			&member.User.FirstName,
			&member.User.LastName,
			&member.User.Phone,
			&member.User.IsActive,
			&member.User.CreatedAt,
			&member.User.UpdatedAt,
			&member.User.DeletedAt,
		); err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	return members, rows.Err()
}

func findOrganizationMemberByID(ctx context.Context, q queryer, organizationID string, userID string) (OrganizationMember, error) {
	const query = `
		SELECT
			ou.organization_id,
			ou.user_id,
			ou.role,
			ou.status,
			ou.invited_by,
			ou.joined_at,
			ou.created_at,
			ou.updated_at,
			u.id,
			u.email,
			COALESCE(u.first_name, ''),
			COALESCE(u.last_name, ''),
			COALESCE(u.phone, ''),
			u.is_active,
			u.created_at,
			u.updated_at,
			u.deleted_at
		FROM organization_users ou
		JOIN users u ON u.id = ou.user_id
		WHERE ou.organization_id = $1 AND ou.user_id = $2
	`

	var member OrganizationMember
	err := q.QueryRow(ctx, query, organizationID, userID).Scan(
		&member.OrganizationID,
		&member.UserID,
		&member.Role,
		&member.Status,
		&member.InvitedBy,
		&member.JoinedAt,
		&member.CreatedAt,
		&member.UpdatedAt,
		&member.User.ID,
		&member.User.Email,
		&member.User.FirstName,
		&member.User.LastName,
		&member.User.Phone,
		&member.User.IsActive,
		&member.User.CreatedAt,
		&member.User.UpdatedAt,
		&member.User.DeletedAt,
	)

	return member, err
}

func updateOrganizationMember(ctx context.Context, db *pgxpool.Pool, organizationID string, userID string, input UpdateOrganizationMemberInput) (OrganizationMember, error) {
	const query = `
		UPDATE organization_users ou
		SET
			role = COALESCE($3, ou.role),
			status = COALESCE($4, ou.status)
		WHERE ou.organization_id = $1 AND ou.user_id = $2
		RETURNING
			ou.organization_id,
			ou.user_id,
			ou.role,
			ou.status,
			ou.invited_by,
			ou.joined_at,
			ou.created_at,
			ou.updated_at
	`

	var member OrganizationMember
	err := db.QueryRow(ctx, query, organizationID, userID, input.Role, input.Status).Scan(
		&member.OrganizationID,
		&member.UserID,
		&member.Role,
		&member.Status,
		&member.InvitedBy,
		&member.JoinedAt,
		&member.CreatedAt,
		&member.UpdatedAt,
	)
	if err != nil {
		return OrganizationMember{}, err
	}

	member.User, err = findUserForMember(ctx, db, userID)
	if err != nil {
		return OrganizationMember{}, err
	}

	return member, nil
}

func insertOrganizationInvitation(ctx context.Context, db *pgxpool.Pool, organizationID string, invitedBy string, input CreateOrganizationInvitationInput) (OrganizationInvitation, error) {
	const query = `
		INSERT INTO organization_invitations (organization_id, email, role, expires_at, invited_by)
		VALUES ($1, $2, $3, $4, NULLIF($5, '')::uuid)
		RETURNING id, organization_id, email, role, token, expires_at, accepted_at, invited_by, created_at
	`

	var invitation OrganizationInvitation
	err := db.QueryRow(ctx, query, organizationID, input.Email, input.Role, input.ExpiresAt, invitedBy).Scan(
		&invitation.ID,
		&invitation.OrganizationID,
		&invitation.Email,
		&invitation.Role,
		&invitation.Token,
		&invitation.ExpiresAt,
		&invitation.AcceptedAt,
		&invitation.InvitedBy,
		&invitation.CreatedAt,
	)

	return invitation, err
}

func listOrganizationInvitations(ctx context.Context, db *pgxpool.Pool, organizationID string) ([]OrganizationInvitation, error) {
	const query = `
		SELECT id, organization_id, email, role, token, expires_at, accepted_at, invited_by, created_at
		FROM organization_invitations
		WHERE organization_id = $1
		ORDER BY created_at DESC
	`

	rows, err := db.Query(ctx, query, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invitations []OrganizationInvitation
	for rows.Next() {
		var invitation OrganizationInvitation
		if err := rows.Scan(
			&invitation.ID,
			&invitation.OrganizationID,
			&invitation.Email,
			&invitation.Role,
			&invitation.Token,
			&invitation.ExpiresAt,
			&invitation.AcceptedAt,
			&invitation.InvitedBy,
			&invitation.CreatedAt,
		); err != nil {
			return nil, err
		}
		invitations = append(invitations, invitation)
	}

	return invitations, rows.Err()
}

func findOrganizationInvitationByToken(ctx context.Context, q queryRower, token string) (OrganizationInvitation, error) {
	const query = `
		SELECT id, organization_id, email, role, token, expires_at, accepted_at, invited_by, created_at
		FROM organization_invitations
		WHERE token = $1
	`

	var invitation OrganizationInvitation
	err := q.QueryRow(ctx, query, token).Scan(
		&invitation.ID,
		&invitation.OrganizationID,
		&invitation.Email,
		&invitation.Role,
		&invitation.Token,
		&invitation.ExpiresAt,
		&invitation.AcceptedAt,
		&invitation.InvitedBy,
		&invitation.CreatedAt,
	)

	return invitation, err
}

func acceptOrganizationInvitation(ctx context.Context, tx pgx.Tx, invitation OrganizationInvitation, userID string) error {
	const markAccepted = `
		UPDATE organization_invitations
		SET accepted_at = NOW()
		WHERE id = $1 AND accepted_at IS NULL
	`
	const upsertMembership = `
		INSERT INTO organization_users (organization_id, user_id, role, status, invited_by, joined_at)
		VALUES ($1, $2, $3, 'ACTIVE', $4, NOW())
		ON CONFLICT (organization_id, user_id) DO UPDATE
		SET
			role = EXCLUDED.role,
			status = 'ACTIVE',
			invited_by = EXCLUDED.invited_by,
			joined_at = COALESCE(organization_users.joined_at, EXCLUDED.joined_at)
	`

	tag, err := tx.Exec(ctx, markAccepted, invitation.ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	_, err = tx.Exec(ctx, upsertMembership, invitation.OrganizationID, userID, invitation.Role, invitation.InvitedBy)
	return err
}

func insertOrganizationMembership(ctx context.Context, tx pgx.Tx, organizationID string, userID string, role string) error {
	const query = `
		INSERT INTO organization_users (organization_id, user_id, role, status, joined_at)
		VALUES ($1, $2, $3, 'ACTIVE', NOW())
		ON CONFLICT (organization_id, user_id) DO UPDATE
		SET
			role = EXCLUDED.role,
			status = 'ACTIVE',
			joined_at = COALESCE(organization_users.joined_at, EXCLUDED.joined_at)
	`

	_, err := tx.Exec(ctx, query, organizationID, userID, role)
	return err
}

func findUserForMember(ctx context.Context, q queryRower, userID string) (users.User, error) {
	const query = `
		SELECT id, email, COALESCE(first_name, ''), COALESCE(last_name, ''), COALESCE(phone, ''), is_active, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1
	`

	var user users.User
	err := q.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	return user, err
}

type queryRower interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

type queryer interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func jsonbOrEmpty(value json.RawMessage) []byte {
	if len(value) == 0 {
		return []byte("{}")
	}

	return []byte(value)
}
