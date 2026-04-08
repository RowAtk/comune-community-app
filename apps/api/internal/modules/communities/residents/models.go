package residents

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Resident struct {
	ID               string     `json:"id"`
	OrganizationID   string     `json:"organization_id"`
	CommunityID      string     `json:"community_id"`
	UnitID           string     `json:"unit_id"`
	UserID           *string    `json:"user_id,omitempty"`
	HouseholdID      *string    `json:"household_id,omitempty"`
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	Email            string     `json:"email,omitempty"`
	Phone            string     `json:"phone,omitempty"`
	ResidentType     string     `json:"resident_type"`
	Status           string     `json:"status"`
	IsPrimaryContact bool       `json:"is_primary_contact"`
	MoveInDate       *time.Time `json:"move_in_date,omitempty"`
	MoveOutDate      *time.Time `json:"move_out_date,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

type CreateInput struct {
	OrganizationID   string     `json:"organization_id"`
	CommunityID      string     `json:"community_id"`
	UnitID           string     `json:"unit_id"`
	UserID           string     `json:"user_id"`
	HouseholdID      string     `json:"household_id"`
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	Email            string     `json:"email"`
	Phone            string     `json:"phone"`
	ResidentType     string     `json:"resident_type"`
	Status           string     `json:"status"`
	IsPrimaryContact bool       `json:"is_primary_contact"`
	MoveInDate       *time.Time `json:"move_in_date"`
	MoveOutDate      *time.Time `json:"move_out_date"`
}

type UpdateInput struct {
	UnitID           *string     `json:"unit_id"`
	UserID           *string     `json:"user_id"`
	HouseholdID      *string     `json:"household_id"`
	FirstName        *string     `json:"first_name"`
	LastName         *string     `json:"last_name"`
	Email            *string     `json:"email"`
	Phone            *string     `json:"phone"`
	ResidentType     *string     `json:"resident_type"`
	Status           *string     `json:"status"`
	IsPrimaryContact *bool       `json:"is_primary_contact"`
	MoveInDate       **time.Time `json:"move_in_date"`
	MoveOutDate      **time.Time `json:"move_out_date"`
}

type Service struct {
	db *pgxpool.Pool
}
