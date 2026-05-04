package invoicing

import (
	"comune/apps/api/internal/platform/apperror"
	"comune/apps/api/internal/platform/db"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	validAdminOrganizationRoles = map[string]struct{}{
		"OWNER":         {},
		"ORG_ADMIN":     {},
		"BILLING_ADMIN": {},
		"OPERATIONS":    {},
	}
	validAdminCommunityRoles = map[string]struct{}{
		"COMMUNITY_ADMIN": {},
		"MANAGER":         {},
		"BOARD_MEMBER":    {},
	}
)

var (
	ErrInvoicingOrganizationRequired  = apperror.Error{Kind: apperror.KindValidation, Code: "organization_id_required", Message: "organization id is required"}
	ErrInvoicingCommunityRequired     = apperror.Error{Kind: apperror.KindValidation, Code: "community_id_required", Message: "community id is required"}
	ErrInvoicingUserRequired          = apperror.Error{Kind: apperror.KindUnauthorized, Code: "authentication_required", Message: "authentication is required"}
	ErrInvoicePlanNameRequired        = apperror.Error{Kind: apperror.KindValidation, Code: "invoice_plan_name_required", Message: "invoice plan name is required"}
	ErrInvalidIssueDay                = apperror.Error{Kind: apperror.KindValidation, Code: "invalid_issue_day_of_month", Message: "issue day must be between 1 and 28"}
	ErrInvalidDueDay                  = apperror.Error{Kind: apperror.KindValidation, Code: "invalid_due_day_of_month", Message: "due day must be between 1 and 28"}
	ErrInvalidDefaultAmount           = apperror.Error{Kind: apperror.KindValidation, Code: "invalid_default_amount", Message: "default amount must be greater than or equal to zero"}
	ErrStartsOnRequired               = apperror.Error{Kind: apperror.KindValidation, Code: "starts_on_required", Message: "starts on date is required"}
	ErrInvalidDateRange               = apperror.Error{Kind: apperror.KindValidation, Code: "invalid_date_range", Message: "end date must not be before start date"}
	ErrInsufficientInvoicingAccess    = apperror.Error{Kind: apperror.KindUnauthorized, Code: "insufficient_invoicing_permissions", Message: "insufficient permissions to manage invoicing"}
	ErrInvoicePlanNotFound            = apperror.Error{Kind: apperror.KindNotFound, Code: "invoice_plan_not_found", Message: "invoice plan not found"}
	ErrInvalidBillingPeriod           = apperror.Error{Kind: apperror.KindValidation, Code: "invalid_billing_period", Message: "billing period must use YYYY-MM"}
	ErrInactiveInvoicePlan            = apperror.Error{Kind: apperror.KindValidation, Code: "invoice_plan_not_active", Message: "invoice plan is not active"}
	ErrUnitIDRequired                 = apperror.Error{Kind: apperror.KindValidation, Code: "unit_id_required", Message: "unit id is required"}
	ErrInvoiceAmountRequired          = apperror.Error{Kind: apperror.KindValidation, Code: "invoice_amount_required", Message: "invoice amount must be greater than zero"}
	ErrInvoiceDueDateRequired         = apperror.Error{Kind: apperror.KindValidation, Code: "invoice_due_date_required", Message: "due date is required"}
	ErrUnitNotFound                   = apperror.Error{Kind: apperror.KindNotFound, Code: "unit_not_found", Message: "unit not found"}
	ErrResidentInvoiceOverviewMissing = apperror.Error{Kind: apperror.KindNotFound, Code: "resident_invoice_overview_not_found", Message: "no linked resident found for this community"}
	ErrCreateInvoicePlanFailed        = apperror.Error{Kind: apperror.KindInternal, Code: "create_invoice_plan_failed", Message: "internal server error"}
	ErrListInvoicePlansFailed         = apperror.Error{Kind: apperror.KindInternal, Code: "list_invoice_plans_failed", Message: "internal server error"}
	ErrGenerateInvoicesFailed         = apperror.Error{Kind: apperror.KindInternal, Code: "generate_invoices_failed", Message: "internal server error"}
	ErrCreateInvoiceFailed            = apperror.Error{Kind: apperror.KindInternal, Code: "create_invoice_failed", Message: "internal server error"}
	ErrListInvoicesFailed             = apperror.Error{Kind: apperror.KindInternal, Code: "list_invoices_failed", Message: "internal server error"}
	ErrListOverdueHouseholdsFailed    = apperror.Error{Kind: apperror.KindInternal, Code: "list_overdue_households_failed", Message: "internal server error"}
	ErrGetResidentOverviewFailed      = apperror.Error{Kind: apperror.KindInternal, Code: "get_resident_invoice_overview_failed", Message: "internal server error"}
)

func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

func (s *Service) ListPlans(ctx context.Context, organizationID string, communityID string, actorUserID string) ([]InvoicePlan, error) {
	if err := s.authorizeAdmin(ctx, organizationID, communityID, actorUserID); err != nil {
		return nil, err
	}

	plans, err := listInvoicePlans(ctx, s.db, organizationID, communityID)
	if err != nil {
		return nil, ErrListInvoicePlansFailed.Wrap(fmt.Errorf("list invoice plans: %w", err))
	}

	return plans, nil
}

func (s *Service) CreatePlan(ctx context.Context, actorUserID string, input CreateInvoicePlanInput) (InvoicePlan, error) {
	actorUserID = strings.TrimSpace(actorUserID)
	input.OrganizationID = strings.TrimSpace(input.OrganizationID)
	input.CommunityID = strings.TrimSpace(input.CommunityID)
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)

	if actorUserID == "" {
		return InvoicePlan{}, ErrInvoicingUserRequired.Wrap(nil)
	}
	if input.OrganizationID == "" {
		return InvoicePlan{}, ErrInvoicingOrganizationRequired.Wrap(nil)
	}
	if input.CommunityID == "" {
		return InvoicePlan{}, ErrInvoicingCommunityRequired.Wrap(nil)
	}
	if input.Name == "" {
		return InvoicePlan{}, ErrInvoicePlanNameRequired.Wrap(nil)
	}
	if input.IssueDayOfMonth < 1 || input.IssueDayOfMonth > 28 {
		return InvoicePlan{}, ErrInvalidIssueDay.Wrap(nil)
	}
	if input.DueDayOfMonth < 1 || input.DueDayOfMonth > 28 {
		return InvoicePlan{}, ErrInvalidDueDay.Wrap(nil)
	}
	if input.DefaultAmount < 0 {
		return InvoicePlan{}, ErrInvalidDefaultAmount.Wrap(nil)
	}
	if err := s.authorizeAdmin(ctx, input.OrganizationID, input.CommunityID, actorUserID); err != nil {
		return InvoicePlan{}, err
	}

	startsOn, err := parseDateString(input.StartsOn)
	if err != nil || startsOn == nil {
		return InvoicePlan{}, ErrStartsOnRequired.Wrap(err)
	}
	endsOn, err := parseDateString(input.EndsOn)
	if err != nil {
		return InvoicePlan{}, ErrInvalidDateRange.Wrap(err)
	}
	if endsOn != nil && endsOn.Before(*startsOn) {
		return InvoicePlan{}, ErrInvalidDateRange.Wrap(nil)
	}

	plan, err := insertInvoicePlan(ctx, s.db, actorUserID, input, *startsOn, endsOn)
	if err != nil {
		return InvoicePlan{}, ErrCreateInvoicePlanFailed.Wrap(fmt.Errorf("create invoice plan: %w", err))
	}

	return plan, nil
}

func (s *Service) GenerateInvoices(ctx context.Context, organizationID string, communityID string, planID string, actorUserID string, input GenerateInvoicesInput) (GenerateInvoicesResult, error) {
	if err := s.authorizeAdmin(ctx, organizationID, communityID, actorUserID); err != nil {
		return GenerateInvoicesResult{}, err
	}

	periodStart, billingPeriod, err := parseBillingPeriod(input.BillingPeriod)
	if err != nil {
		return GenerateInvoicesResult{}, ErrInvalidBillingPeriod.Wrap(err)
	}

	var result GenerateInvoicesResult
	err = db.RunInTx(ctx, s.db, func(tx pgx.Tx) error {
		plan, err := findInvoicePlanByID(ctx, tx, organizationID, communityID, planID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrInvoicePlanNotFound.Wrap(nil)
			}
			return ErrGenerateInvoicesFailed.Wrap(fmt.Errorf("load invoice plan: %w", err))
		}
		if plan.Status != "ACTIVE" {
			return ErrInactiveInvoicePlan.Wrap(nil)
		}

		issuedOn, dueDate := buildPeriodDates(periodStart, plan.IssueDayOfMonth, plan.DueDayOfMonth)
		createdCount, err := generateInvoicesForPlan(ctx, tx, plan, billingPeriod, issuedOn, dueDate)
		if err != nil {
			return ErrGenerateInvoicesFailed.Wrap(fmt.Errorf("generate invoices: %w", err))
		}

		result = GenerateInvoicesResult{
			PlanID:        plan.ID,
			BillingPeriod: billingPeriod,
			CreatedCount:  createdCount,
		}

		return nil
	})
	if err != nil {
		return GenerateInvoicesResult{}, err
	}

	return result, nil
}

func (s *Service) ListInvoices(ctx context.Context, organizationID string, communityID string, actorUserID string) ([]Invoice, error) {
	if err := s.authorizeAdmin(ctx, organizationID, communityID, actorUserID); err != nil {
		return nil, err
	}

	invoices, err := listInvoices(ctx, s.db, organizationID, communityID)
	if err != nil {
		return nil, ErrListInvoicesFailed.Wrap(fmt.Errorf("list invoices: %w", err))
	}

	return invoices, nil
}

func (s *Service) CreateManualInvoice(ctx context.Context, actorUserID string, input CreateInvoiceInput) (Invoice, error) {
	actorUserID = strings.TrimSpace(actorUserID)
	input.OrganizationID = strings.TrimSpace(input.OrganizationID)
	input.CommunityID = strings.TrimSpace(input.CommunityID)
	input.UnitID = strings.TrimSpace(input.UnitID)
	input.BillingPeriod = strings.TrimSpace(input.BillingPeriod)
	input.Description = strings.TrimSpace(input.Description)

	if actorUserID == "" {
		return Invoice{}, ErrInvoicingUserRequired.Wrap(nil)
	}
	if input.OrganizationID == "" {
		return Invoice{}, ErrInvoicingOrganizationRequired.Wrap(nil)
	}
	if input.CommunityID == "" {
		return Invoice{}, ErrInvoicingCommunityRequired.Wrap(nil)
	}
	if input.UnitID == "" {
		return Invoice{}, ErrUnitIDRequired.Wrap(nil)
	}
	if input.Amount <= 0 {
		return Invoice{}, ErrInvoiceAmountRequired.Wrap(nil)
	}
	if err := s.authorizeAdmin(ctx, input.OrganizationID, input.CommunityID, actorUserID); err != nil {
		return Invoice{}, err
	}
	if _, _, err := parseBillingPeriod(input.BillingPeriod); err != nil {
		return Invoice{}, ErrInvalidBillingPeriod.Wrap(err)
	}
	dueDate, err := parseDateString(input.DueDate)
	if err != nil || dueDate == nil {
		return Invoice{}, ErrInvoiceDueDateRequired.Wrap(err)
	}
	unitExists, err := communityUnitExists(ctx, s.db, input.OrganizationID, input.CommunityID, input.UnitID)
	if err != nil {
		return Invoice{}, ErrCreateInvoiceFailed.Wrap(fmt.Errorf("check unit existence: %w", err))
	}
	if !unitExists {
		return Invoice{}, ErrUnitNotFound.Wrap(nil)
	}

	invoice, err := insertManualInvoice(ctx, s.db, input, time.Now().UTC(), *dueDate)
	if err != nil {
		return Invoice{}, ErrCreateInvoiceFailed.Wrap(fmt.Errorf("create manual invoice: %w", err))
	}

	return invoice, nil
}

func (s *Service) ListOverdueHouseholds(ctx context.Context, organizationID string, communityID string, actorUserID string) ([]OverdueHouseholdSummary, error) {
	if err := s.authorizeAdmin(ctx, organizationID, communityID, actorUserID); err != nil {
		return nil, err
	}

	summaries, err := listOverdueHouseholds(ctx, s.db, organizationID, communityID)
	if err != nil {
		return nil, ErrListOverdueHouseholdsFailed.Wrap(fmt.Errorf("list overdue households: %w", err))
	}

	return summaries, nil
}

func (s *Service) GetResidentOverview(ctx context.Context, organizationID string, communityID string, actorUserID string) (ResidentInvoiceOverview, error) {
	organizationID = strings.TrimSpace(organizationID)
	communityID = strings.TrimSpace(communityID)
	actorUserID = strings.TrimSpace(actorUserID)

	if organizationID == "" {
		return ResidentInvoiceOverview{}, ErrInvoicingOrganizationRequired.Wrap(nil)
	}
	if communityID == "" {
		return ResidentInvoiceOverview{}, ErrInvoicingCommunityRequired.Wrap(nil)
	}
	if actorUserID == "" {
		return ResidentInvoiceOverview{}, ErrInvoicingUserRequired.Wrap(nil)
	}

	resident, err := findLinkedResidentForUser(ctx, s.db, organizationID, communityID, actorUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ResidentInvoiceOverview{}, ErrResidentInvoiceOverviewMissing.Wrap(nil)
		}
		return ResidentInvoiceOverview{}, ErrGetResidentOverviewFailed.Wrap(fmt.Errorf("load linked resident: %w", err))
	}

	invoices, err := listInvoicesForUnit(ctx, s.db, organizationID, communityID, resident.UnitID)
	if err != nil {
		return ResidentInvoiceOverview{}, ErrGetResidentOverviewFailed.Wrap(fmt.Errorf("list resident invoices: %w", err))
	}

	return computeResidentOverview(resident, invoices), nil
}

func (s *Service) authorizeAdmin(ctx context.Context, organizationID string, communityID string, actorUserID string) error {
	organizationID = strings.TrimSpace(organizationID)
	communityID = strings.TrimSpace(communityID)
	actorUserID = strings.TrimSpace(actorUserID)

	if organizationID == "" {
		return ErrInvoicingOrganizationRequired.Wrap(nil)
	}
	if communityID == "" {
		return ErrInvoicingCommunityRequired.Wrap(nil)
	}
	if actorUserID == "" {
		return ErrInvoicingUserRequired.Wrap(nil)
	}

	if orgRole, err := findOrganizationRoleForUser(ctx, s.db, organizationID, actorUserID); err == nil {
		if _, ok := validAdminOrganizationRoles[orgRole]; ok {
			return nil
		}
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return ErrInsufficientInvoicingAccess.Wrap(fmt.Errorf("load organization role: %w", err))
	}

	if communityRole, err := findCommunityRoleForUser(ctx, s.db, organizationID, communityID, actorUserID); err == nil {
		if _, ok := validAdminCommunityRoles[communityRole]; ok {
			return nil
		}
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return ErrInsufficientInvoicingAccess.Wrap(fmt.Errorf("load community role: %w", err))
	}

	return ErrInsufficientInvoicingAccess.Wrap(nil)
}
