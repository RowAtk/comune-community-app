package residents

import (
	"encoding/json"
	"fmt"
	"strings"
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
	HouseholdRole    *string    `json:"household_role,omitempty"`
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
	HouseholdRole    string     `json:"household_role"`
	Status           string     `json:"status"`
	IsPrimaryContact bool       `json:"is_primary_contact"`
	MoveInDate       *DateInput `json:"move_in_date"`
	MoveOutDate      *DateInput `json:"move_out_date"`
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
	HouseholdRole    *string     `json:"household_role"`
	Status           *string     `json:"status"`
	IsPrimaryContact *bool       `json:"is_primary_contact"`
	MoveInDate       *NullableDateInput `json:"move_in_date"`
	MoveOutDate      *NullableDateInput `json:"move_out_date"`
}

type Service struct {
	db *pgxpool.Pool
}

type DateInput struct {
	Time time.Time
}

func (d *DateInput) UnmarshalJSON(data []byte) error {
	parsed, err := parseDateInput(data)
	if err != nil {
		return err
	}
	if parsed == nil {
		return nil
	}

	d.Time = *parsed
	return nil
}

type NullableDateInput struct {
	Value *time.Time
}

func (d *NullableDateInput) UnmarshalJSON(data []byte) error {
	parsed, err := parseDateInput(data)
	if err != nil {
		return err
	}

	d.Value = parsed
	return nil
}

func parseDateInput(data []byte) (*time.Time, error) {
	if string(data) == "null" {
		return nil, nil
	}

	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("date must be a JSON string in YYYY-MM-DD or RFC3339 format: %w", err)
	}

	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}

	formats := []string{
		"2006-01-02",
		time.RFC3339,
		time.RFC3339Nano,
	}

	var lastErr error
	for _, format := range formats {
		parsed, err := time.Parse(format, raw)
		if err == nil {
			return &parsed, nil
		}
		lastErr = err
	}

	return nil, fmt.Errorf("invalid date %q: expected YYYY-MM-DD or RFC3339: %w", raw, lastErr)
}
