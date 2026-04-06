package users

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func insertUser(ctx context.Context, tx pgx.Tx, params NewUserParams) (User, error) {
	const query = `
		INSERT INTO users (email, first_name, last_name, phone)
		VALUES ($1, NULLIF($2, ''), NULLIF($3, ''), NULLIF($4, ''))
		RETURNING id, email, COALESCE(first_name, ''), COALESCE(last_name, ''), COALESCE(phone, ''), is_active, created_at, updated_at
	`

	var user User
	err := tx.QueryRow(ctx, query, params.Email, params.FirstName, params.LastName, params.Phone).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return user, err
}



