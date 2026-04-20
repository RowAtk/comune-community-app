package residents

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type queryRower interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func insert(ctx context.Context, db *pgxpool.Pool, input CreateInput, moveInDate *time.Time, moveOutDate *time.Time) (Resident, error) {
	const query = `
		INSERT INTO residents (
			organization_id, community_id, unit_id, user_id, household_id, first_name, last_name, email, phone, resident_type, household_role, status, is_primary_contact, move_in_date, move_out_date
		)
		VALUES ($1, $2, $3, NULLIF($4, '')::uuid, NULLIF($5, '')::uuid, $6, $7, NULLIF($8, ''), NULLIF($9, ''), $10, NULLIF($11, ''), $12, $13, $14, $15)
		RETURNING id, organization_id, community_id, unit_id, user_id, household_id, first_name, last_name, COALESCE(email, ''), COALESCE(phone, ''), resident_type, household_role, status, is_primary_contact, move_in_date, move_out_date, created_at, updated_at, deleted_at
	`

	var resident Resident
	err := db.QueryRow(
		ctx,
		query,
		input.OrganizationID,
		input.CommunityID,
		input.UnitID,
		input.UserID,
		input.HouseholdID,
		input.FirstName,
		input.LastName,
		input.Email,
		input.Phone,
		input.ResidentType,
		input.HouseholdRole,
		input.Status,
		input.IsPrimaryContact,
		moveInDate,
		moveOutDate,
	).Scan(scanResidentFields(&resident)...)

	return resident, err
}

func list(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string) ([]Resident, error) {
	const query = `
		SELECT
			r.id,
			r.organization_id,
			r.community_id,
			r.unit_id,
			r.user_id,
			COALESCE(u.email, ''),
			TRIM(CONCAT(COALESCE(u.first_name, ''), ' ', COALESCE(u.last_name, ''))),
			COALESCE(u.phone, ''),
			r.household_id,
			r.first_name,
			r.last_name,
			COALESCE(r.email, ''),
			COALESCE(r.phone, ''),
			r.resident_type,
			r.household_role,
			r.status,
			r.is_primary_contact,
			r.move_in_date,
			r.move_out_date,
			r.created_at,
			r.updated_at,
			r.deleted_at
		FROM residents r
		LEFT JOIN users u ON u.id = r.user_id
		WHERE organization_id = $1 AND community_id = $2 AND r.deleted_at IS NULL
		ORDER BY r.created_at DESC
	`

	rows, err := db.Query(ctx, query, organizationID, communityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var residents []Resident
	for rows.Next() {
		var resident Resident
		if err := rows.Scan(scanResidentFieldsWithLinkedUser(&resident)...); err != nil {
			return nil, err
		}
		residents = append(residents, resident)
	}

	return residents, rows.Err()
}

func findByID(ctx context.Context, q queryRower, organizationID string, communityID string, residentID string) (Resident, error) {
	const query = `
		SELECT
			r.id,
			r.organization_id,
			r.community_id,
			r.unit_id,
			r.user_id,
			COALESCE(u.email, ''),
			TRIM(CONCAT(COALESCE(u.first_name, ''), ' ', COALESCE(u.last_name, ''))),
			COALESCE(u.phone, ''),
			r.household_id,
			r.first_name,
			r.last_name,
			COALESCE(r.email, ''),
			COALESCE(r.phone, ''),
			r.resident_type,
			r.household_role,
			r.status,
			r.is_primary_contact,
			r.move_in_date,
			r.move_out_date,
			r.created_at,
			r.updated_at,
			r.deleted_at
		FROM residents r
		LEFT JOIN users u ON u.id = r.user_id
		WHERE r.organization_id = $1 AND r.community_id = $2 AND r.id = $3 AND r.deleted_at IS NULL
	`

	var resident Resident
	err := q.QueryRow(ctx, query, organizationID, communityID, residentID).Scan(scanResidentFieldsWithLinkedUser(&resident)...)

	return resident, err
}

func update(
	ctx context.Context,
	db *pgxpool.Pool,
	organizationID string,
	communityID string,
	residentID string,
	input UpdateInput,
	parsedMoveInDate **time.Time,
	parsedMoveOutDate **time.Time,
) (Resident, error) {
	const query = `
		UPDATE residents
		SET
			unit_id = COALESCE($4::text::uuid, unit_id),
			user_id = CASE WHEN $5::text IS NULL THEN user_id ELSE NULLIF($5::text, '')::uuid END,
			household_id = CASE WHEN $6::text IS NULL THEN household_id ELSE NULLIF($6::text, '')::uuid END,
			first_name = COALESCE($7::text, first_name),
			last_name = COALESCE($8::text, last_name),
			email = CASE WHEN $9::text IS NULL THEN email ELSE NULLIF($9::text, '') END,
			phone = CASE WHEN $10::text IS NULL THEN phone ELSE NULLIF($10::text, '') END,
			resident_type = COALESCE($11::text, resident_type),
			household_role = CASE WHEN $12::text IS NULL THEN household_role ELSE NULLIF($12::text, '') END,
			status = COALESCE($13::text, status),
			is_primary_contact = COALESCE($14::boolean, is_primary_contact),
			move_in_date = CASE WHEN $15::date IS NULL THEN move_in_date ELSE $15::date END,
			move_out_date = CASE WHEN $16::date IS NULL THEN move_out_date ELSE $16::date END
		WHERE organization_id = $1 AND community_id = $2 AND id = $3 AND deleted_at IS NULL
		RETURNING id, organization_id, community_id, unit_id, user_id, household_id, first_name, last_name, COALESCE(email, ''), COALESCE(phone, ''), resident_type, household_role, status, is_primary_contact, move_in_date, move_out_date, created_at, updated_at, deleted_at
	`

	var (
		moveInDateArg  any
		moveOutDateArg any
	)
	if parsedMoveInDate != nil {
		moveInDateArg = *parsedMoveInDate
	}
	if parsedMoveOutDate != nil {
		moveOutDateArg = *parsedMoveOutDate
	}

	var resident Resident
	err := db.QueryRow(
		ctx,
		query,
		organizationID,
		communityID,
		residentID,
		input.UnitID,
		input.UserID,
		input.HouseholdID,
		input.FirstName,
		input.LastName,
		input.Email,
		input.Phone,
		input.ResidentType,
		input.HouseholdRole,
		input.Status,
		input.IsPrimaryContact,
		moveInDateArg,
		moveOutDateArg,
	).Scan(scanResidentFields(&resident)...)

	return resident, err
}

func remove(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string, residentID string) error {
	const query = `
		UPDATE residents
		SET deleted_at = NOW()
		WHERE organization_id = $1 AND community_id = $2 AND id = $3 AND deleted_at IS NULL
	`

	tag, err := db.Exec(ctx, query, organizationID, communityID, residentID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func unitExists(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string, unitID string) (bool, error) {
	const query = `
		SELECT 1
		FROM units
		WHERE organization_id = $1 AND community_id = $2 AND id = $3 AND deleted_at IS NULL
	`

	var exists int
	err := db.QueryRow(ctx, query, organizationID, communityID, unitID).Scan(&exists)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func findHouseholdUnitID(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string, householdID string) (string, bool, error) {
	const query = `
		SELECT unit_id
		FROM households
		WHERE organization_id = $1 AND community_id = $2 AND id = $3
	`

	var unitID string
	err := db.QueryRow(ctx, query, organizationID, communityID, householdID).Scan(&unitID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", false, nil
		}
		return "", false, err
	}

	return unitID, true, nil
}

func insertInvitation(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string, residentID string, invitedBy string, input CreateInvitationInput) (ResidentInvitation, error) {
	const query = `
		INSERT INTO resident_invitations (organization_id, community_id, resident_id, expires_at, invited_by)
		VALUES ($1, $2, $3, $4, NULLIF($5, '')::uuid)
		RETURNING id, organization_id, community_id, resident_id, token, expires_at, accepted_at, invited_by, accepted_by_user_id, created_at
	`

	var invitation ResidentInvitation
	err := db.QueryRow(ctx, query, organizationID, communityID, residentID, input.ExpiresAt, invitedBy).Scan(
		&invitation.ID,
		&invitation.OrganizationID,
		&invitation.CommunityID,
		&invitation.ResidentID,
		&invitation.Token,
		&invitation.ExpiresAt,
		&invitation.AcceptedAt,
		&invitation.InvitedBy,
		&invitation.AcceptedByUserID,
		&invitation.CreatedAt,
	)

	return invitation, err
}

func listInvitations(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string, residentID string) ([]ResidentInvitation, error) {
	const query = `
		SELECT id, organization_id, community_id, resident_id, token, expires_at, accepted_at, invited_by, accepted_by_user_id, created_at
		FROM resident_invitations
		WHERE organization_id = $1 AND community_id = $2 AND resident_id = $3
		ORDER BY created_at DESC
	`

	rows, err := db.Query(ctx, query, organizationID, communityID, residentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invitations []ResidentInvitation
	for rows.Next() {
		var invitation ResidentInvitation
		if err := rows.Scan(
			&invitation.ID,
			&invitation.OrganizationID,
			&invitation.CommunityID,
			&invitation.ResidentID,
			&invitation.Token,
			&invitation.ExpiresAt,
			&invitation.AcceptedAt,
			&invitation.InvitedBy,
			&invitation.AcceptedByUserID,
			&invitation.CreatedAt,
		); err != nil {
			return nil, err
		}
		invitations = append(invitations, invitation)
	}

	return invitations, rows.Err()
}

func findInvitationByToken(ctx context.Context, q queryRower, token string) (ResidentInvitation, error) {
	const query = `
		SELECT id, organization_id, community_id, resident_id, token, expires_at, accepted_at, invited_by, accepted_by_user_id, created_at
		FROM resident_invitations
		WHERE token = $1
	`

	var invitation ResidentInvitation
	err := q.QueryRow(ctx, query, token).Scan(
		&invitation.ID,
		&invitation.OrganizationID,
		&invitation.CommunityID,
		&invitation.ResidentID,
		&invitation.Token,
		&invitation.ExpiresAt,
		&invitation.AcceptedAt,
		&invitation.InvitedBy,
		&invitation.AcceptedByUserID,
		&invitation.CreatedAt,
	)

	return invitation, err
}

func findInvitationPreviewByToken(ctx context.Context, q queryRower, token string) (ResidentInvitationPreview, error) {
	const query = `
		SELECT
			ri.id,
			ri.organization_id,
			o.name,
			ri.community_id,
			c.name,
			ri.resident_id,
			TRIM(CONCAT(r.first_name, ' ', r.last_name)),
			r.resident_type,
			r.household_id,
			COALESCE(h.name, ''),
			r.unit_id,
			u.unit_number,
			r.household_role,
			r.status,
			ri.token,
			ri.expires_at,
			ri.accepted_at,
			ri.invited_by
		FROM resident_invitations ri
		JOIN organizations o ON o.id = ri.organization_id
		JOIN communities c ON c.id = ri.community_id
		JOIN residents r ON r.id = ri.resident_id
		JOIN units u ON u.id = r.unit_id
		LEFT JOIN households h ON h.id = r.household_id
		WHERE ri.token = $1 AND r.deleted_at IS NULL
	`

	var preview ResidentInvitationPreview
	err := q.QueryRow(ctx, query, token).Scan(
		&preview.ID,
		&preview.OrganizationID,
		&preview.Organization,
		&preview.CommunityID,
		&preview.Community,
		&preview.ResidentID,
		&preview.ResidentName,
		&preview.ResidentType,
		&preview.HouseholdID,
		&preview.HouseholdName,
		&preview.UnitID,
		&preview.UnitNumber,
		&preview.HouseholdRole,
		&preview.Status,
		&preview.Token,
		&preview.ExpiresAt,
		&preview.AcceptedAt,
		&preview.InvitedBy,
	)

	return preview, err
}

func scanResidentFields(resident *Resident) []any {
	return []any{
		&resident.ID,
		&resident.OrganizationID,
		&resident.CommunityID,
		&resident.UnitID,
		&resident.UserID,
		&resident.HouseholdID,
		&resident.FirstName,
		&resident.LastName,
		&resident.Email,
		&resident.Phone,
		&resident.ResidentType,
		&resident.HouseholdRole,
		&resident.Status,
		&resident.IsPrimaryContact,
		&resident.MoveInDate,
		&resident.MoveOutDate,
		&resident.CreatedAt,
		&resident.UpdatedAt,
		&resident.DeletedAt,
	}
}

func scanResidentFieldsWithLinkedUser(resident *Resident) []any {
	return []any{
		&resident.ID,
		&resident.OrganizationID,
		&resident.CommunityID,
		&resident.UnitID,
		&resident.UserID,
		&resident.LinkedUserEmail,
		&resident.LinkedUserName,
		&resident.LinkedUserPhone,
		&resident.HouseholdID,
		&resident.FirstName,
		&resident.LastName,
		&resident.Email,
		&resident.Phone,
		&resident.ResidentType,
		&resident.HouseholdRole,
		&resident.Status,
		&resident.IsPrimaryContact,
		&resident.MoveInDate,
		&resident.MoveOutDate,
		&resident.CreatedAt,
		&resident.UpdatedAt,
		&resident.DeletedAt,
	}
}

func acceptInvitation(ctx context.Context, tx pgx.Tx, invitation ResidentInvitation, userID string) error {
	const markAccepted = `
		UPDATE resident_invitations
		SET accepted_at = NOW(), accepted_by_user_id = $2
		WHERE id = $1 AND accepted_at IS NULL
	`
	const linkResident = `
		UPDATE residents
		SET user_id = CASE
			WHEN user_id IS NULL THEN $2
			WHEN user_id = $2 THEN user_id
			ELSE user_id
		END
		WHERE id = $1 AND deleted_at IS NULL
	`
	const upsertMembership = `
		INSERT INTO community_users (organization_id, community_id, user_id, role, status, invited_by, joined_at)
		VALUES ($1, $2, $3, 'RESIDENT', 'ACTIVE', $4, NOW())
		ON CONFLICT (community_id, user_id) DO UPDATE
		SET
			role = EXCLUDED.role,
			status = 'ACTIVE',
			invited_by = EXCLUDED.invited_by,
			joined_at = COALESCE(community_users.joined_at, EXCLUDED.joined_at)
	`

	tag, err := tx.Exec(ctx, markAccepted, invitation.ID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	if _, err := tx.Exec(ctx, linkResident, invitation.ResidentID, userID); err != nil {
		return err
	}

	_, err = tx.Exec(ctx, upsertMembership, invitation.OrganizationID, invitation.CommunityID, userID, invitation.InvitedBy)
	return err
}

func findActiveUserIDByResident(ctx context.Context, q queryRower, organizationID string, communityID string, residentID string) (*string, error) {
	const query = `
		SELECT user_id
		FROM residents
		WHERE organization_id = $1 AND community_id = $2 AND id = $3 AND deleted_at IS NULL
	`

	var userID *string
	err := q.QueryRow(ctx, query, organizationID, communityID, residentID).Scan(&userID)
	return userID, err
}

func findOrganizationRoleForUser(ctx context.Context, q queryRower, organizationID string, userID string) (string, error) {
	const query = `
		SELECT role
		FROM organization_users
		WHERE organization_id = $1 AND user_id = $2 AND status = 'ACTIVE'
	`

	var role string
	err := q.QueryRow(ctx, query, organizationID, userID).Scan(&role)
	return role, err
}

func findCommunityRoleForUser(ctx context.Context, q queryRower, organizationID string, communityID string, userID string) (string, error) {
	const query = `
		SELECT role
		FROM community_users
		WHERE organization_id = $1 AND community_id = $2 AND user_id = $3 AND status = 'ACTIVE'
	`

	var role string
	err := q.QueryRow(ctx, query, organizationID, communityID, userID).Scan(&role)
	return role, err
}
