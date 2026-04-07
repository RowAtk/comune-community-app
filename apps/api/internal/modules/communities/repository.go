package communities

import (
	"comune/apps/api/internal/modules/users"
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func insertCommunity(ctx context.Context, db *pgxpool.Pool, input CreateCommunityInput) (Community, error) {
	const query = `
		INSERT INTO communities (
			organization_id, name, slug, address, timezone, status, settings
		)
		VALUES ($1, $2, $3, NULLIF($4, ''), $5, $6, $7::jsonb)
		RETURNING id, organization_id, name, slug, COALESCE(address, ''), timezone, status, settings, created_at, updated_at, deleted_at
	`

	var community Community
	err := db.QueryRow(
		ctx,
		query,
		input.OrganizationID,
		input.Name,
		input.Slug,
		input.Address,
		input.Timezone,
		input.Status,
		jsonbOrEmpty(input.Settings),
	).Scan(
		&community.ID,
		&community.OrganizationID,
		&community.Name,
		&community.Slug,
		&community.Address,
		&community.Timezone,
		&community.Status,
		&community.Settings,
		&community.CreatedAt,
		&community.UpdatedAt,
		&community.DeletedAt,
	)

	return community, err
}

func listCommunities(ctx context.Context, db *pgxpool.Pool, organizationID string) ([]Community, error) {
	const query = `
		SELECT id, organization_id, name, slug, COALESCE(address, ''), timezone, status, settings, created_at, updated_at, deleted_at
		FROM communities
		WHERE organization_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := db.Query(ctx, query, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var communities []Community
	for rows.Next() {
		var community Community
		if err := rows.Scan(
			&community.ID,
			&community.OrganizationID,
			&community.Name,
			&community.Slug,
			&community.Address,
			&community.Timezone,
			&community.Status,
			&community.Settings,
			&community.CreatedAt,
			&community.UpdatedAt,
			&community.DeletedAt,
		); err != nil {
			return nil, err
		}

		communities = append(communities, community)
	}

	return communities, rows.Err()
}

func findCommunityByID(ctx context.Context, db *pgxpool.Pool, organizationID string, id string) (Community, error) {
	const query = `
		SELECT id, organization_id, name, slug, COALESCE(address, ''), timezone, status, settings, created_at, updated_at, deleted_at
		FROM communities
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var community Community
	err := db.QueryRow(ctx, query, organizationID, id).Scan(
		&community.ID,
		&community.OrganizationID,
		&community.Name,
		&community.Slug,
		&community.Address,
		&community.Timezone,
		&community.Status,
		&community.Settings,
		&community.CreatedAt,
		&community.UpdatedAt,
		&community.DeletedAt,
	)

	return community, err
}

func updateCommunity(ctx context.Context, db *pgxpool.Pool, organizationID string, id string, input UpdateCommunityInput) (Community, error) {
	const query = `
		UPDATE communities
		SET
			name = COALESCE($3, name),
			slug = COALESCE($4, slug),
			address = CASE WHEN $5 IS NULL THEN address ELSE NULLIF($5, '') END,
			timezone = COALESCE($6, timezone),
			status = COALESCE($7, status),
			settings = COALESCE($8::jsonb, settings)
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
		RETURNING id, organization_id, name, slug, COALESCE(address, ''), timezone, status, settings, created_at, updated_at, deleted_at
	`

	var (
		community Community
		settings  any
	)

	if input.Settings != nil {
		settings = jsonbOrEmpty(*input.Settings)
	}

	err := db.QueryRow(
		ctx,
		query,
		organizationID,
		id,
		input.Name,
		input.Slug,
		input.Address,
		input.Timezone,
		input.Status,
		settings,
	).Scan(
		&community.ID,
		&community.OrganizationID,
		&community.Name,
		&community.Slug,
		&community.Address,
		&community.Timezone,
		&community.Status,
		&community.Settings,
		&community.CreatedAt,
		&community.UpdatedAt,
		&community.DeletedAt,
	)

	return community, err
}

func deleteCommunity(ctx context.Context, db *pgxpool.Pool, organizationID string, id string) error {
	const query = `
		UPDATE communities
		SET deleted_at = NOW()
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	tag, err := db.Exec(ctx, query, organizationID, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func listCommunityMembers(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string) ([]CommunityMember, error) {
	const query = `
		SELECT
			cu.organization_id,
			cu.community_id,
			cu.user_id,
			cu.role,
			cu.status,
			cu.invited_by,
			cu.joined_at,
			cu.created_at,
			cu.updated_at,
			u.id,
			u.email,
			COALESCE(u.first_name, ''),
			COALESCE(u.last_name, ''),
			COALESCE(u.phone, ''),
			u.is_active,
			u.created_at,
			u.updated_at,
			u.deleted_at
		FROM community_users cu
		JOIN users u ON u.id = cu.user_id
		WHERE cu.organization_id = $1 AND cu.community_id = $2
		ORDER BY cu.created_at DESC
	`

	rows, err := db.Query(ctx, query, organizationID, communityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []CommunityMember
	for rows.Next() {
		var member CommunityMember
		if err := rows.Scan(
			&member.OrganizationID,
			&member.CommunityID,
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

func findCommunityMemberByID(ctx context.Context, q queryRower, organizationID string, communityID string, userID string) (CommunityMember, error) {
	const query = `
		SELECT
			cu.organization_id,
			cu.community_id,
			cu.user_id,
			cu.role,
			cu.status,
			cu.invited_by,
			cu.joined_at,
			cu.created_at,
			cu.updated_at,
			u.id,
			u.email,
			COALESCE(u.first_name, ''),
			COALESCE(u.last_name, ''),
			COALESCE(u.phone, ''),
			u.is_active,
			u.created_at,
			u.updated_at,
			u.deleted_at
		FROM community_users cu
		JOIN users u ON u.id = cu.user_id
		WHERE cu.organization_id = $1 AND cu.community_id = $2 AND cu.user_id = $3
	`

	var member CommunityMember
	err := q.QueryRow(ctx, query, organizationID, communityID, userID).Scan(
		&member.OrganizationID,
		&member.CommunityID,
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

func updateCommunityMember(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string, userID string, input UpdateCommunityMemberInput) (CommunityMember, error) {
	const query = `
		UPDATE community_users cu
		SET
			role = COALESCE($4, cu.role),
			status = COALESCE($5, cu.status)
		WHERE cu.organization_id = $1 AND cu.community_id = $2 AND cu.user_id = $3
		RETURNING cu.organization_id, cu.community_id, cu.user_id, cu.role, cu.status, cu.invited_by, cu.joined_at, cu.created_at, cu.updated_at
	`

	var member CommunityMember
	err := db.QueryRow(ctx, query, organizationID, communityID, userID, input.Role, input.Status).Scan(
		&member.OrganizationID,
		&member.CommunityID,
		&member.UserID,
		&member.Role,
		&member.Status,
		&member.InvitedBy,
		&member.JoinedAt,
		&member.CreatedAt,
		&member.UpdatedAt,
	)
	if err != nil {
		return CommunityMember{}, err
	}

	member.User, err = findUserForMember(ctx, db, userID)
	if err != nil {
		return CommunityMember{}, err
	}

	return member, nil
}

func insertCommunityInvitation(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string, invitedBy string, input CreateCommunityInvitationInput) (CommunityInvitation, error) {
	const query = `
		INSERT INTO community_invitations (organization_id, community_id, email, role, expires_at, invited_by)
		VALUES ($1, $2, $3, $4, $5, NULLIF($6, '')::uuid)
		RETURNING id, organization_id, community_id, email, role, token, expires_at, accepted_at, invited_by, created_at
	`

	var invitation CommunityInvitation
	err := db.QueryRow(ctx, query, organizationID, communityID, input.Email, input.Role, input.ExpiresAt, invitedBy).Scan(
		&invitation.ID,
		&invitation.OrganizationID,
		&invitation.CommunityID,
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

func listCommunityInvitations(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string) ([]CommunityInvitation, error) {
	const query = `
		SELECT id, organization_id, community_id, email, role, token, expires_at, accepted_at, invited_by, created_at
		FROM community_invitations
		WHERE organization_id = $1 AND community_id = $2
		ORDER BY created_at DESC
	`

	rows, err := db.Query(ctx, query, organizationID, communityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invitations []CommunityInvitation
	for rows.Next() {
		var invitation CommunityInvitation
		if err := rows.Scan(
			&invitation.ID,
			&invitation.OrganizationID,
			&invitation.CommunityID,
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

func findCommunityInvitationByToken(ctx context.Context, q queryRower, token string) (CommunityInvitation, error) {
	const query = `
		SELECT id, organization_id, community_id, email, role, token, expires_at, accepted_at, invited_by, created_at
		FROM community_invitations
		WHERE token = $1
	`

	var invitation CommunityInvitation
	err := q.QueryRow(ctx, query, token).Scan(
		&invitation.ID,
		&invitation.OrganizationID,
		&invitation.CommunityID,
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

func acceptCommunityInvitation(ctx context.Context, tx pgx.Tx, invitation CommunityInvitation, userID string) error {
	const markAccepted = `
		UPDATE community_invitations
		SET accepted_at = NOW()
		WHERE id = $1 AND accepted_at IS NULL
	`
	const upsertMembership = `
		INSERT INTO community_users (organization_id, community_id, user_id, role, status, invited_by, joined_at)
		VALUES ($1, $2, $3, $4, 'ACTIVE', $5, NOW())
		ON CONFLICT (community_id, user_id) DO UPDATE
		SET
			role = EXCLUDED.role,
			status = 'ACTIVE',
			invited_by = EXCLUDED.invited_by,
			joined_at = COALESCE(community_users.joined_at, EXCLUDED.joined_at)
	`

	tag, err := tx.Exec(ctx, markAccepted, invitation.ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	_, err = tx.Exec(ctx, upsertMembership, invitation.OrganizationID, invitation.CommunityID, userID, invitation.Role, invitation.InvitedBy)
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

func jsonbOrEmpty(value json.RawMessage) []byte {
	if len(value) == 0 {
		return []byte("{}")
	}

	return []byte(value)
}
