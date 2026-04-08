package households

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func insert(ctx context.Context, db *pgxpool.Pool, input CreateInput) (Household, error) {
	const query = `
		INSERT INTO households (
			organization_id, community_id, unit_id, name
		)
		VALUES ($1, $2, $3, NULLIF($4, ''))
		RETURNING id, organization_id, community_id, unit_id, COALESCE(name, ''), created_at, updated_at
	`

	var household Household
	err := db.QueryRow(ctx, query, input.OrganizationID, input.CommunityID, input.UnitID, input.Name).Scan(
		&household.ID,
		&household.OrganizationID,
		&household.CommunityID,
		&household.UnitID,
		&household.Name,
		&household.CreatedAt,
		&household.UpdatedAt,
	)

	return household, err
}

func list(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string) ([]Household, error) {
	const query = `
		SELECT id, organization_id, community_id, unit_id, COALESCE(name, ''), created_at, updated_at
		FROM households
		WHERE organization_id = $1 AND community_id = $2
		ORDER BY created_at DESC
	`

	rows, err := db.Query(ctx, query, organizationID, communityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var households []Household
	for rows.Next() {
		var household Household
		if err := rows.Scan(
			&household.ID,
			&household.OrganizationID,
			&household.CommunityID,
			&household.UnitID,
			&household.Name,
			&household.CreatedAt,
			&household.UpdatedAt,
		); err != nil {
			return nil, err
		}
		households = append(households, household)
	}

	return households, rows.Err()
}

func findByID(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string, householdID string) (Household, error) {
	const query = `
		SELECT id, organization_id, community_id, unit_id, COALESCE(name, ''), created_at, updated_at
		FROM households
		WHERE organization_id = $1 AND community_id = $2 AND id = $3
	`

	var household Household
	err := db.QueryRow(ctx, query, organizationID, communityID, householdID).Scan(
		&household.ID,
		&household.OrganizationID,
		&household.CommunityID,
		&household.UnitID,
		&household.Name,
		&household.CreatedAt,
		&household.UpdatedAt,
	)

	return household, err
}

func update(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string, householdID string, input UpdateInput) (Household, error) {
	const query = `
		UPDATE households
		SET
			unit_id = COALESCE($4, unit_id),
			name = CASE WHEN $5 IS NULL THEN name ELSE NULLIF($5, '') END
		WHERE organization_id = $1 AND community_id = $2 AND id = $3
		RETURNING id, organization_id, community_id, unit_id, COALESCE(name, ''), created_at, updated_at
	`

	var household Household
	err := db.QueryRow(ctx, query, organizationID, communityID, householdID, input.UnitID, input.Name).Scan(
		&household.ID,
		&household.OrganizationID,
		&household.CommunityID,
		&household.UnitID,
		&household.Name,
		&household.CreatedAt,
		&household.UpdatedAt,
	)

	return household, err
}

func remove(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string, householdID string) error {
	const query = `
		DELETE FROM households
		WHERE organization_id = $1 AND community_id = $2 AND id = $3
	`

	tag, err := db.Exec(ctx, query, organizationID, communityID, householdID)
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
