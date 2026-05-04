package invoicing

import (
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InvoicePlan struct {
	ID              string     `json:"id"`
	OrganizationID  string     `json:"organization_id"`
	CommunityID     string     `json:"community_id"`
	Name            string     `json:"name"`
	PlanType        string     `json:"plan_type"`
	Status          string     `json:"status"`
	IssueDayOfMonth int        `json:"issue_day_of_month"`
	DueDayOfMonth   int        `json:"due_day_of_month"`
	DefaultAmount   float64    `json:"default_amount"`
	StartsOn        time.Time  `json:"starts_on"`
	EndsOn          *time.Time `json:"ends_on,omitempty"`
	Description     string     `json:"description,omitempty"`
	CreatedBy       *string    `json:"created_by,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
	OverrideCount   int        `json:"override_count"`
	ActiveUnitCount int        `json:"active_unit_count"`
}

type Invoice struct {
	ID                string     `json:"id"`
	OrganizationID    string     `json:"organization_id"`
	CommunityID       string     `json:"community_id"`
	UnitID            string     `json:"unit_id"`
	UnitNumber        string     `json:"unit_number"`
	HouseholdID       *string    `json:"household_id,omitempty"`
	HouseholdName     string     `json:"household_name,omitempty"`
	InvoicePlanID     *string    `json:"invoice_plan_id,omitempty"`
	Source            string     `json:"source"`
	Amount            float64    `json:"amount"`
	PaidAmount        float64    `json:"paid_amount"`
	OutstandingAmount float64    `json:"outstanding_amount"`
	IssuedOn          *time.Time `json:"issued_on,omitempty"`
	DueDate           time.Time  `json:"due_date"`
	Status            string     `json:"status"`
	BillingPeriod     string     `json:"billing_period"`
	Description       string     `json:"description,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	IsOverdue         bool       `json:"is_overdue"`
}

type OverdueHouseholdSummary struct {
	HouseholdID         *string   `json:"household_id,omitempty"`
	HouseholdName       string    `json:"household_name"`
	UnitID              string    `json:"unit_id"`
	UnitNumber          string    `json:"unit_number"`
	PrimaryResidentName string    `json:"primary_resident_name,omitempty"`
	OverdueInvoiceCount int       `json:"overdue_invoice_count"`
	TotalOutstanding    float64   `json:"total_outstanding"`
	EarliestDueDate     time.Time `json:"earliest_due_date"`
}

type ResidentInvoiceOverview struct {
	ResidentID          string     `json:"resident_id"`
	ResidentName        string     `json:"resident_name"`
	ResidentStatus      string     `json:"resident_status"`
	UnitID              string     `json:"unit_id"`
	UnitNumber          string     `json:"unit_number"`
	HouseholdID         *string    `json:"household_id,omitempty"`
	HouseholdName       string     `json:"household_name,omitempty"`
	NextDueDate         *time.Time `json:"next_due_date,omitempty"`
	OverdueInvoiceCount int        `json:"overdue_invoice_count"`
	TotalOutstanding    float64    `json:"total_outstanding"`
	Invoices            []Invoice  `json:"invoices"`
}

type CreateInvoicePlanInput struct {
	OrganizationID  string  `json:"organization_id"`
	CommunityID     string  `json:"community_id"`
	Name            string  `json:"name"`
	IssueDayOfMonth int     `json:"issue_day_of_month"`
	DueDayOfMonth   int     `json:"due_day_of_month"`
	DefaultAmount   float64 `json:"default_amount"`
	StartsOn        string  `json:"starts_on"`
	EndsOn          string  `json:"ends_on"`
	Description     string  `json:"description"`
}

type GenerateInvoicesInput struct {
	BillingPeriod string `json:"billing_period"`
}

type GenerateInvoicesResult struct {
	PlanID        string `json:"plan_id"`
	BillingPeriod string `json:"billing_period"`
	CreatedCount  int64  `json:"created_count"`
}

type CreateInvoiceInput struct {
	OrganizationID string  `json:"organization_id"`
	CommunityID    string  `json:"community_id"`
	UnitID         string  `json:"unit_id"`
	Amount         float64 `json:"amount"`
	DueDate        string  `json:"due_date"`
	BillingPeriod  string  `json:"billing_period"`
	Description    string  `json:"description"`
}

type linkedResident struct {
	ID            string
	FirstName     string
	LastName      string
	Status        string
	UnitID        string
	UnitNumber    string
	HouseholdID   *string
	HouseholdName string
}

type Service struct {
	db *pgxpool.Pool
}

func parseDateString(raw string) (*time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}

	formats := []string{"2006-01-02", time.RFC3339, time.RFC3339Nano}
	var lastErr error
	for _, format := range formats {
		parsed, err := time.Parse(format, raw)
		if err == nil {
			return &parsed, nil
		}
		lastErr = err
	}

	return nil, fmt.Errorf("invalid date %q: %w", raw, lastErr)
}

func parseBillingPeriod(raw string) (time.Time, string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, "", fmt.Errorf("billing period is required")
	}

	parsed, err := time.Parse("2006-01", raw)
	if err != nil {
		return time.Time{}, "", fmt.Errorf("billing period must use YYYY-MM: %w", err)
	}

	return parsed, parsed.Format("2006-01"), nil
}
