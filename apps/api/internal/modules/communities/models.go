package communities

import (
	"comune/apps/api/internal/modules/users"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Community struct {
	ID             string          `json:"id"`
	OrganizationID string          `json:"organization_id"`
	Name           string          `json:"name"`
	Slug           string          `json:"slug"`
	Address        string          `json:"address,omitempty"`
	Timezone       string          `json:"timezone"`
	Status         string          `json:"status"`
	Settings       json.RawMessage `json:"settings"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      *time.Time      `json:"deleted_at,omitempty"`
}

type CreateCommunityInput struct {
	OrganizationID string          `json:"organization_id"`
	Name           string          `json:"name"`
	Slug           string          `json:"slug"`
	Address        string          `json:"address"`
	Timezone       string          `json:"timezone"`
	Status         string          `json:"status"`
	Settings       json.RawMessage `json:"settings"`
}

type UpdateCommunityInput struct {
	Name     *string          `json:"name"`
	Slug     *string          `json:"slug"`
	Address  *string          `json:"address"`
	Timezone *string          `json:"timezone"`
	Status   *string          `json:"status"`
	Settings *json.RawMessage `json:"settings"`
}

type Service struct {
	db *pgxpool.Pool
}

type CommunityMember struct {
	OrganizationID string     `json:"organization_id"`
	CommunityID    string     `json:"community_id"`
	UserID         string     `json:"user_id"`
	Role           string     `json:"role"`
	Status         string     `json:"status"`
	InvitedBy      *string    `json:"invited_by,omitempty"`
	JoinedAt       *time.Time `json:"joined_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	User           users.User `json:"user"`
}

type UpdateCommunityMemberInput struct {
	Role   *string `json:"role"`
	Status *string `json:"status"`
}

type CommunityInvitation struct {
	ID             string     `json:"id"`
	OrganizationID string     `json:"organization_id"`
	CommunityID    string     `json:"community_id"`
	Email          string     `json:"email"`
	Role           string     `json:"role"`
	Token          string     `json:"token"`
	ExpiresAt      time.Time  `json:"expires_at"`
	AcceptedAt     *time.Time `json:"accepted_at,omitempty"`
	InvitedBy      *string    `json:"invited_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

type CreateCommunityInvitationInput struct {
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type acceptCommunityInvitationParams struct {
	UserID string
	Email  string
	Token  string
}
