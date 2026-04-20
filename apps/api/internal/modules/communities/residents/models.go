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
	LinkedUserEmail  string     `json:"linked_user_email,omitempty"`
	LinkedUserName   string     `json:"linked_user_name,omitempty"`
	LinkedUserPhone  string     `json:"linked_user_phone,omitempty"`
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
	UnitID           *string            `json:"unit_id"`
	UserID           *string            `json:"user_id"`
	HouseholdID      *string            `json:"household_id"`
	FirstName        *string            `json:"first_name"`
	LastName         *string            `json:"last_name"`
	Email            *string            `json:"email"`
	Phone            *string            `json:"phone"`
	ResidentType     *string            `json:"resident_type"`
	HouseholdRole    *string            `json:"household_role"`
	Status           *string            `json:"status"`
	IsPrimaryContact *bool              `json:"is_primary_contact"`
	MoveInDate       *NullableDateInput `json:"move_in_date"`
	MoveOutDate      *NullableDateInput `json:"move_out_date"`
}

type ResidentInvitation struct {
	ID               string     `json:"id"`
	OrganizationID   string     `json:"organization_id"`
	CommunityID      string     `json:"community_id"`
	ResidentID       string     `json:"resident_id"`
	Token            string     `json:"token"`
	ExpiresAt        time.Time  `json:"expires_at"`
	AcceptedAt       *time.Time `json:"accepted_at,omitempty"`
	InvitedBy        *string    `json:"invited_by,omitempty"`
	AcceptedByUserID *string    `json:"accepted_by_user_id,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}

type CreateInvitationInput struct {
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type ResidentInvitationPreview struct {
	ID             string     `json:"id"`
	OrganizationID string     `json:"organization_id"`
	Organization   string     `json:"organization"`
	CommunityID    string     `json:"community_id"`
	Community      string     `json:"community"`
	ResidentID     string     `json:"resident_id"`
	ResidentName   string     `json:"resident_name"`
	ResidentType   string     `json:"resident_type"`
	HouseholdID    *string    `json:"household_id,omitempty"`
	HouseholdName  string     `json:"household_name,omitempty"`
	UnitID         string     `json:"unit_id"`
	UnitNumber     string     `json:"unit_number"`
	HouseholdRole  *string    `json:"household_role,omitempty"`
	Status         string     `json:"status"`
	Token          string     `json:"token"`
	ExpiresAt      time.Time  `json:"expires_at"`
	AcceptedAt     *time.Time `json:"accepted_at,omitempty"`
	InvitedBy      *string    `json:"invited_by,omitempty"`
}

type acceptInvitationParams struct {
	UserID string
	Token  string
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
