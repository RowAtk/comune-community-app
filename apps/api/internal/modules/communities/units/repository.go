package units

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func insert(ctx context.Context, db *pgxpool.Pool, input CreateInput) (Unit, error) {
	const query = `
		INSERT INTO units (
			organization_id, community_id, unit_number, block_floor, unit_type, status
		)
		VALUES ($1, $2, $3, NULLIF($4, ''), NULLIF($5, ''), $6)
		RETURNING id, organization_id, community_id, unit_number, COALESCE(block_floor, ''), COALESCE(unit_type, ''), status, created_at, updated_at, deleted_at
	`

	var unit Unit
	err := db.QueryRow(
		ctx,
		query,
		input.OrganizationID,
		input.CommunityID,
		input.UnitNumber,
		input.BlockFloor,
		input.UnitType,
		input.Status,
	).Scan(
		&unit.ID,
		&unit.OrganizationID,
		&unit.CommunityID,
		&unit.UnitNumber,
		&unit.BlockFloor,
		&unit.UnitType,
		&unit.Status,
		&unit.CreatedAt,
		&unit.UpdatedAt,
		&unit.DeletedAt,
	)

	return unit, err
}

func list(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string) ([]Unit, error) {
	const query = `
		SELECT id, organization_id, community_id, unit_number, COALESCE(block_floor, ''), COALESCE(unit_type, ''), status, created_at, updated_at, deleted_at
		FROM units
		WHERE organization_id = $1 AND community_id = $2 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := db.Query(ctx, query, organizationID, communityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []Unit
	for rows.Next() {
		var unit Unit
		if err := rows.Scan(
			&unit.ID,
			&unit.OrganizationID,
			&unit.CommunityID,
			&unit.UnitNumber,
			&unit.BlockFloor,
			&unit.UnitType,
			&unit.Status,
			&unit.CreatedAt,
			&unit.UpdatedAt,
			&unit.DeletedAt,
		); err != nil {
			return nil, err
		}
		units = append(units, unit)
	}

	return units, rows.Err()
}

func findByID(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string, unitID string) (Unit, error) {
	const query = `
		SELECT id, organization_id, community_id, unit_number, COALESCE(block_floor, ''), COALESCE(unit_type, ''), status, created_at, updated_at, deleted_at
		FROM units
		WHERE organization_id = $1 AND community_id = $2 AND id = $3 AND deleted_at IS NULL
	`

	var unit Unit
	err := db.QueryRow(ctx, query, organizationID, communityID, unitID).Scan(
		&unit.ID,
		&unit.OrganizationID,
		&unit.CommunityID,
		&unit.UnitNumber,
		&unit.BlockFloor,
		&unit.UnitType,
		&unit.Status,
		&unit.CreatedAt,
		&unit.UpdatedAt,
		&unit.DeletedAt,
	)

	return unit, err
}

func update(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string, unitID string, input UpdateInput) (Unit, error) {
	const query = `
		UPDATE units
		SET
			unit_number = COALESCE($4, unit_number),
			block_floor = CASE WHEN $5 IS NULL THEN block_floor ELSE NULLIF($5, '') END,
			unit_type = CASE WHEN $6 IS NULL THEN unit_type ELSE NULLIF($6, '') END,
			status = COALESCE($7, status)
		WHERE organization_id = $1 AND community_id = $2 AND id = $3 AND deleted_at IS NULL
		RETURNING id, organization_id, community_id, unit_number, COALESCE(block_floor, ''), COALESCE(unit_type, ''), status, created_at, updated_at, deleted_at
	`

	var unit Unit
	err := db.QueryRow(
		ctx,
		query,
		organizationID,
		communityID,
		unitID,
		input.UnitNumber,
		input.BlockFloor,
		input.UnitType,
		input.Status,
	).Scan(
		&unit.ID,
		&unit.OrganizationID,
		&unit.CommunityID,
		&unit.UnitNumber,
		&unit.BlockFloor,
		&unit.UnitType,
		&unit.Status,
		&unit.CreatedAt,
		&unit.UpdatedAt,
		&unit.DeletedAt,
	)

	return unit, err
}

func remove(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string, unitID string) error {
	const query = `
		UPDATE units
		SET deleted_at = NOW()
		WHERE organization_id = $1 AND community_id = $2 AND id = $3 AND deleted_at IS NULL
	`

	tag, err := db.Exec(ctx, query, organizationID, communityID, unitID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
