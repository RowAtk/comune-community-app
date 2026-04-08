package units

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Unit struct {
	ID             string     `json:"id"`
	OrganizationID string     `json:"organization_id"`
	CommunityID    string     `json:"community_id"`
	UnitNumber     string     `json:"unit_number"`
	BlockFloor     string     `json:"block_floor,omitempty"`
	UnitType       string     `json:"unit_type,omitempty"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type CreateInput struct {
	OrganizationID string `json:"organization_id"`
	CommunityID    string `json:"community_id"`
	UnitNumber     string `json:"unit_number"`
	BlockFloor     string `json:"block_floor"`
	UnitType       string `json:"unit_type"`
	Status         string `json:"status"`
}

type UpdateInput struct {
	UnitNumber *string `json:"unit_number"`
	BlockFloor *string `json:"block_floor"`
	UnitType   *string `json:"unit_type"`
	Status     *string `json:"status"`
}

type Service struct {
	db *pgxpool.Pool
}
