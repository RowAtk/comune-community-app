package auth

import (
	"comune/apps/api/internal/modules/users"
	"context"
)

func (s *Service) findUserByEmailWithPasswordProvider(ctx context.Context, email string) (users.User, AuthAccount, error) {
	const query = `
		SELECT
			u.id,
			u.email,
			COALESCE(u.first_name, ''),
			COALESCE(u.last_name, ''),
			COALESCE(u.phone, ''),
			u.is_active,
			u.created_at,
			u.updated_at,
			u.deleted_at,
			a.id,
			a.user_id,
			a.provider,
			a.provider_user_id,
			COALESCE(a.email_at_provider, ''),
			a.password_hash,
			a.last_login_at,
			a.created_at,
			a.updated_at
		FROM users u
		JOIN auth_accounts a ON a.user_id = u.id
		WHERE u.email = $1
		  AND u.deleted_at IS NULL
		  AND a.provider = $2
		  AND a.provider_user_id = $1
	`

	return scanUserWithAccount(s.db.QueryRow(ctx, query, email, PasswordProvider))
}

func (s *Service) findUserByIDWithPasswordProvider(ctx context.Context, userID string) (users.User, AuthAccount, error) {
	const query = `
		SELECT
			u.id,
			u.email,
			COALESCE(u.first_name, ''),
			COALESCE(u.last_name, ''),
			COALESCE(u.phone, ''),
			u.is_active,
			u.created_at,
			u.updated_at,
			u.deleted_at,
			a.id,
			a.user_id,
			a.provider,
			a.provider_user_id,
			COALESCE(a.email_at_provider, ''),
			a.password_hash,
			a.last_login_at,
			a.created_at,
			a.updated_at
		FROM users u
		JOIN auth_accounts a ON a.user_id = u.id
		WHERE u.id = $1
		  AND u.deleted_at IS NULL
		  AND a.provider = $2
	`

	return scanUserWithAccount(s.db.QueryRow(ctx, query, userID, PasswordProvider))
}

func scanUserWithAccount(row rowScanner) (users.User, AuthAccount, error) {
	var (
		user    users.User
		account AuthAccount
	)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
		&account.ID,
		&account.UserID,
		&account.Provider,
		&account.ProviderUserID,
		&account.EmailAtProvider,
		&account.PasswordHash,
		&account.LastLoginAt,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	return user, account, err
}

type rowScanner interface {
	Scan(dest ...any) error
}
