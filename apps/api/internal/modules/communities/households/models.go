package households

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Household struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	CommunityID    string    `json:"community_id"`
	UnitID         string    `json:"unit_id"`
	Name           string    `json:"name,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateInput struct {
	OrganizationID string `json:"organization_id"`
	CommunityID    string `json:"community_id"`
	UnitID         string `json:"unit_id"`
	Name           string `json:"name"`
}

type UpdateInput struct {
	UnitID *string `json:"unit_id"`
	Name   *string `json:"name"`
}

type Service struct {
	db *pgxpool.Pool
}
