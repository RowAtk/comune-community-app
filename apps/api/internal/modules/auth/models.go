package auth

import (
	"comune/apps/api/internal/modules/users"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthAccount struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	Provider        string    `json:"provider"`
	ProviderUserID  string    `json:"provider_user_id"`
	EmailAtProvider string    `json:"email_at_provider,omitempty"`
	PasswordHash    string    `json:"-"`
	LastLoginAt     time.Time `json:"last_login_at,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Session struct {
	Token     string    `json:"token"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type SignupInput struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResult struct {
	User        users.User  `json:"user"`
	AuthAccount AuthAccount `json:"auth_account"`
	Session     Session     `json:"session"`
}

type DashboardOrganizationMembership struct {
	OrganizationID string    `json:"organization_id"`
	Name           string    `json:"name"`
	Slug           string    `json:"slug"`
	Role           string    `json:"role"`
	Status         string    `json:"status"`
	JoinedAt       time.Time `json:"joined_at"`
}

type DashboardCommunityMembership struct {
	OrganizationID   string    `json:"organization_id"`
	OrganizationName string    `json:"organization_name"`
	OrganizationSlug string    `json:"organization_slug"`
	CommunityID      string    `json:"community_id"`
	CommunityName    string    `json:"community_name"`
	CommunitySlug    string    `json:"community_slug"`
	Role             string    `json:"role"`
	Status           string    `json:"status"`
	JoinedAt         time.Time `json:"joined_at"`
}

type DashboardResult struct {
	Organizations []DashboardOrganizationMembership `json:"organizations"`
	Communities   []DashboardCommunityMembership    `json:"communities"`
}

type Service struct {
	db              *pgxpool.Pool
	users           *users.Service
	sessionSecret   []byte
	sessionDuration time.Duration
}
