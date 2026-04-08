package residents

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func insert(ctx context.Context, db *pgxpool.Pool, input CreateInput) (Resident, error) {
	const query = `
		INSERT INTO residents (
			organization_id, community_id, unit_id, user_id, household_id, first_name, last_name, email, phone, resident_type, status, is_primary_contact, move_in_date, move_out_date
		)
		VALUES ($1, $2, $3, NULLIF($4, '')::uuid, NULLIF($5, '')::uuid, $6, $7, NULLIF($8, ''), NULLIF($9, ''), $10, $11, $12, $13, $14)
		RETURNING id, organization_id, community_id, unit_id, user_id, household_id, first_name, last_name, COALESCE(email, ''), COALESCE(phone, ''), resident_type, status, is_primary_contact, move_in_date, move_out_date, created_at, updated_at, deleted_at
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
		input.Status,
		input.IsPrimaryContact,
		input.MoveInDate,
		input.MoveOutDate,
	).Scan(
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
		&resident.Status,
		&resident.IsPrimaryContact,
		&resident.MoveInDate,
		&resident.MoveOutDate,
		&resident.CreatedAt,
		&resident.UpdatedAt,
		&resident.DeletedAt,
	)

	return resident, err
}

func list(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string) ([]Resident, error) {
	const query = `
		SELECT id, organization_id, community_id, unit_id, user_id, household_id, first_name, last_name, COALESCE(email, ''), COALESCE(phone, ''), resident_type, status, is_primary_contact, move_in_date, move_out_date, created_at, updated_at, deleted_at
		FROM residents
		WHERE organization_id = $1 AND community_id = $2 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := db.Query(ctx, query, organizationID, communityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var residents []Resident
	for rows.Next() {
		var resident Resident
		if err := rows.Scan(
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
			&resident.Status,
			&resident.IsPrimaryContact,
			&resident.MoveInDate,
			&resident.MoveOutDate,
			&resident.CreatedAt,
			&resident.UpdatedAt,
			&resident.DeletedAt,
		); err != nil {
			return nil, err
		}
		residents = append(residents, resident)
	}

	return residents, rows.Err()
}

func findByID(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string, residentID string) (Resident, error) {
	const query = `
		SELECT id, organization_id, community_id, unit_id, user_id, household_id, first_name, last_name, COALESCE(email, ''), COALESCE(phone, ''), resident_type, status, is_primary_contact, move_in_date, move_out_date, created_at, updated_at, deleted_at
		FROM residents
		WHERE organization_id = $1 AND community_id = $2 AND id = $3 AND deleted_at IS NULL
	`

	var resident Resident
	err := db.QueryRow(ctx, query, organizationID, communityID, residentID).Scan(
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
		&resident.Status,
		&resident.IsPrimaryContact,
		&resident.MoveInDate,
		&resident.MoveOutDate,
		&resident.CreatedAt,
		&resident.UpdatedAt,
		&resident.DeletedAt,
	)

	return resident, err
}

func update(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string, residentID string, input UpdateInput) (Resident, error) {
	const query = `
		UPDATE residents
		SET
			unit_id = COALESCE($4, unit_id),
			user_id = CASE WHEN $5 IS NULL THEN user_id ELSE NULLIF($5, '')::uuid END,
			household_id = CASE WHEN $6 IS NULL THEN household_id ELSE NULLIF($6, '')::uuid END,
			first_name = COALESCE($7, first_name),
			last_name = COALESCE($8, last_name),
			email = CASE WHEN $9 IS NULL THEN email ELSE NULLIF($9, '') END,
			phone = CASE WHEN $10 IS NULL THEN phone ELSE NULLIF($10, '') END,
			resident_type = COALESCE($11, resident_type),
			status = COALESCE($12, status),
			is_primary_contact = COALESCE($13, is_primary_contact),
			move_in_date = CASE WHEN $14 IS NULL THEN move_in_date ELSE $14 END,
			move_out_date = CASE WHEN $15 IS NULL THEN move_out_date ELSE $15 END
		WHERE organization_id = $1 AND community_id = $2 AND id = $3 AND deleted_at IS NULL
		RETURNING id, organization_id, community_id, unit_id, user_id, household_id, first_name, last_name, COALESCE(email, ''), COALESCE(phone, ''), resident_type, status, is_primary_contact, move_in_date, move_out_date, created_at, updated_at, deleted_at
	`

	var (
		moveInDate  any
		moveOutDate any
	)
	if input.MoveInDate != nil {
		moveInDate = *input.MoveInDate
	}
	if input.MoveOutDate != nil {
		moveOutDate = *input.MoveOutDate
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
		input.Status,
		input.IsPrimaryContact,
		moveInDate,
		moveOutDate,
	).Scan(
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
		&resident.Status,
		&resident.IsPrimaryContact,
		&resident.MoveInDate,
		&resident.MoveOutDate,
		&resident.CreatedAt,
		&resident.UpdatedAt,
		&resident.DeletedAt,
	)

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
