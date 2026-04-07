package organizations

import (
	"comune/apps/api/internal/modules/users"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Organization struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Slug         string          `json:"slug"`
	LegalName    string          `json:"legal_name,omitempty"`
	BillingEmail string          `json:"billing_email,omitempty"`
	Phone        string          `json:"phone,omitempty"`
	CountryCode  string          `json:"country_code"`
	Timezone     string          `json:"timezone"`
	Status       string          `json:"status"`
	Settings     json.RawMessage `json:"settings"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    *time.Time      `json:"deleted_at,omitempty"`
}

type CreateOrganizationInput struct {
	Name         string          `json:"name"`
	Slug         string          `json:"slug"`
	LegalName    string          `json:"legal_name"`
	BillingEmail string          `json:"billing_email"`
	Phone        string          `json:"phone"`
	CountryCode  string          `json:"country_code"`
	Timezone     string          `json:"timezone"`
	Status       string          `json:"status"`
	Settings     json.RawMessage `json:"settings"`
}

type UpdateOrganizationInput struct {
	Name         *string          `json:"name"`
	Slug         *string          `json:"slug"`
	LegalName    *string          `json:"legal_name"`
	BillingEmail *string          `json:"billing_email"`
	Phone        *string          `json:"phone"`
	CountryCode  *string          `json:"country_code"`
	Timezone     *string          `json:"timezone"`
	Status       *string          `json:"status"`
	Settings     *json.RawMessage `json:"settings"`
}

type Service struct {
	db *pgxpool.Pool
}

type OrganizationMember struct {
	OrganizationID string     `json:"organization_id"`
	UserID         string     `json:"user_id"`
	Role           string     `json:"role"`
	Status         string     `json:"status"`
	InvitedBy      *string    `json:"invited_by,omitempty"`
	JoinedAt       *time.Time `json:"joined_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	User           users.User `json:"user"`
}

type UpdateOrganizationMemberInput struct {
	Role   *string `json:"role"`
	Status *string `json:"status"`
}

type OrganizationInvitation struct {
	ID             string     `json:"id"`
	OrganizationID string     `json:"organization_id"`
	Email          string     `json:"email"`
	Role           string     `json:"role"`
	Token          string     `json:"token"`
	ExpiresAt      time.Time  `json:"expires_at"`
	AcceptedAt     *time.Time `json:"accepted_at,omitempty"`
	InvitedBy      *string    `json:"invited_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

type CreateOrganizationInvitationInput struct {
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type acceptOrganizationInvitationParams struct {
	UserID string
	Email  string
	Token  string
}
