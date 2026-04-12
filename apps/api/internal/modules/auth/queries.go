package auth

import (
	"comune/apps/api/internal/modules/users"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
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

func listDashboardOrganizations(ctx context.Context, db *pgxpool.Pool, userID string) ([]DashboardOrganizationMembership, error) {
	const query = `
		SELECT
			ou.organization_id,
			o.name,
			o.slug,
			ou.role,
			ou.status,
			COALESCE(ou.joined_at, ou.created_at)
		FROM organization_users ou
		JOIN organizations o ON o.id = ou.organization_id
		WHERE ou.user_id = $1
		  AND ou.status = 'ACTIVE'
		  AND o.deleted_at IS NULL
		ORDER BY o.name ASC
	`

	rows, err := db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var organizations []DashboardOrganizationMembership
	for rows.Next() {
		var membership DashboardOrganizationMembership
		if err := rows.Scan(
			&membership.OrganizationID,
			&membership.Name,
			&membership.Slug,
			&membership.Role,
			&membership.Status,
			&membership.JoinedAt,
		); err != nil {
			return nil, err
		}

		organizations = append(organizations, membership)
	}

	return organizations, rows.Err()
}

func listDashboardCommunities(ctx context.Context, db *pgxpool.Pool, userID string) ([]DashboardCommunityMembership, error) {
	const query = `
		SELECT
			cu.organization_id,
			o.name,
			o.slug,
			cu.community_id,
			c.name,
			c.slug,
			cu.role,
			cu.status,
			COALESCE(cu.joined_at, cu.created_at)
		FROM community_users cu
		JOIN organizations o ON o.id = cu.organization_id
		JOIN communities c ON c.id = cu.community_id
		WHERE cu.user_id = $1
		  AND cu.status = 'ACTIVE'
		  AND o.deleted_at IS NULL
		  AND c.deleted_at IS NULL
		ORDER BY o.name ASC, c.name ASC
	`

	rows, err := db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var communities []DashboardCommunityMembership
	for rows.Next() {
		var membership DashboardCommunityMembership
		if err := rows.Scan(
			&membership.OrganizationID,
			&membership.OrganizationName,
			&membership.OrganizationSlug,
			&membership.CommunityID,
			&membership.CommunityName,
			&membership.CommunitySlug,
			&membership.Role,
			&membership.Status,
			&membership.JoinedAt,
		); err != nil {
			return nil, err
		}

		communities = append(communities, membership)
	}

	return communities, rows.Err()
}
