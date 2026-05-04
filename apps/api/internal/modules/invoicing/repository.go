package invoicing

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type queryExecutor interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}

func insertInvoicePlan(ctx context.Context, db *pgxpool.Pool, createdBy string, input CreateInvoicePlanInput, startsOn time.Time, endsOn *time.Time) (InvoicePlan, error) {
	const query = `
		INSERT INTO invoice_plans (
			organization_id, community_id, name, plan_type, status, issue_day_of_month, due_day_of_month, default_amount, starts_on, ends_on, description, created_by
		)
		VALUES ($1, $2, $3, 'MAINTENANCE', 'ACTIVE', $4, $5, $6, $7, $8, NULLIF($9, ''), NULLIF($10, '')::uuid)
		RETURNING id, organization_id, community_id, name, plan_type, status, issue_day_of_month, due_day_of_month, default_amount, starts_on, ends_on, COALESCE(description, ''), created_by, created_at, updated_at, deleted_at
	`

	var plan InvoicePlan
	err := db.QueryRow(
		ctx,
		query,
		input.OrganizationID,
		input.CommunityID,
		input.Name,
		input.IssueDayOfMonth,
		input.DueDayOfMonth,
		input.DefaultAmount,
		startsOn,
		endsOn,
		input.Description,
		createdBy,
	).Scan(
		&plan.ID,
		&plan.OrganizationID,
		&plan.CommunityID,
		&plan.Name,
		&plan.PlanType,
		&plan.Status,
		&plan.IssueDayOfMonth,
		&plan.DueDayOfMonth,
		&plan.DefaultAmount,
		&plan.StartsOn,
		&plan.EndsOn,
		&plan.Description,
		&plan.CreatedBy,
		&plan.CreatedAt,
		&plan.UpdatedAt,
		&plan.DeletedAt,
	)
	if err != nil {
		return InvoicePlan{}, err
	}

	return plan, nil
}

func listInvoicePlans(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string) ([]InvoicePlan, error) {
	const query = `
		SELECT
			p.id,
			p.organization_id,
			p.community_id,
			p.name,
			p.plan_type,
			p.status,
			p.issue_day_of_month,
			p.due_day_of_month,
			p.default_amount,
			p.starts_on,
			p.ends_on,
			COALESCE(p.description, ''),
			p.created_by,
			p.created_at,
			p.updated_at,
			p.deleted_at,
			COUNT(DISTINCT o.unit_id) AS override_count,
			COUNT(DISTINCT u.id) FILTER (WHERE u.status = 'ACTIVE' AND u.deleted_at IS NULL) AS active_unit_count
		FROM invoice_plans p
		LEFT JOIN invoice_plan_unit_overrides o
			ON o.invoice_plan_id = p.id
		LEFT JOIN units u
			ON u.organization_id = p.organization_id
		   AND u.community_id = p.community_id
		WHERE p.organization_id = $1
		  AND p.community_id = $2
		  AND p.deleted_at IS NULL
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`

	rows, err := db.Query(ctx, query, organizationID, communityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []InvoicePlan
	for rows.Next() {
		var plan InvoicePlan
		if err := rows.Scan(
			&plan.ID,
			&plan.OrganizationID,
			&plan.CommunityID,
			&plan.Name,
			&plan.PlanType,
			&plan.Status,
			&plan.IssueDayOfMonth,
			&plan.DueDayOfMonth,
			&plan.DefaultAmount,
			&plan.StartsOn,
			&plan.EndsOn,
			&plan.Description,
			&plan.CreatedBy,
			&plan.CreatedAt,
			&plan.UpdatedAt,
			&plan.DeletedAt,
			&plan.OverrideCount,
			&plan.ActiveUnitCount,
		); err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}

	return plans, rows.Err()
}

func findInvoicePlanByID(ctx context.Context, q queryExecutor, organizationID string, communityID string, planID string) (InvoicePlan, error) {
	const query = `
		SELECT
			id, organization_id, community_id, name, plan_type, status, issue_day_of_month, due_day_of_month, default_amount, starts_on, ends_on,
			COALESCE(description, ''), created_by, created_at, updated_at, deleted_at
		FROM invoice_plans
		WHERE organization_id = $1 AND community_id = $2 AND id = $3 AND deleted_at IS NULL
	`

	var plan InvoicePlan
	err := q.QueryRow(ctx, query, organizationID, communityID, planID).Scan(
		&plan.ID,
		&plan.OrganizationID,
		&plan.CommunityID,
		&plan.Name,
		&plan.PlanType,
		&plan.Status,
		&plan.IssueDayOfMonth,
		&plan.DueDayOfMonth,
		&plan.DefaultAmount,
		&plan.StartsOn,
		&plan.EndsOn,
		&plan.Description,
		&plan.CreatedBy,
		&plan.CreatedAt,
		&plan.UpdatedAt,
		&plan.DeletedAt,
	)

	return plan, err
}

func countActiveUnitsForCommunity(ctx context.Context, q queryExecutor, organizationID string, communityID string) (int64, error) {
	const query = `
		SELECT COUNT(*)
		FROM units
		WHERE organization_id = $1
		  AND community_id = $2
		  AND status = 'ACTIVE'
		  AND deleted_at IS NULL
	`

	var count int64
	err := q.QueryRow(ctx, query, organizationID, communityID).Scan(&count)
	return count, err
}

func generateInvoicesForPlan(ctx context.Context, tx pgx.Tx, plan InvoicePlan, billingPeriod string, issuedOn time.Time, dueDate time.Time) (int64, error) {
	const query = `
		INSERT INTO invoices (
			organization_id, community_id, unit_id, invoice_plan_id, source, amount, paid_amount, issued_on, due_date, status, billing_period, description
		)
		SELECT
			p.organization_id,
			p.community_id,
			u.id,
			p.id,
			'SCHEDULED',
			COALESCE(o.amount, p.default_amount),
			0.00,
			$4,
			$5,
			CASE WHEN $5 < CURRENT_DATE THEN 'OVERDUE' ELSE 'UNPAID' END,
			$3,
			COALESCE(NULLIF(p.description, ''), p.name || ' invoice')
		FROM invoice_plans p
		JOIN units u
			ON u.organization_id = p.organization_id
		   AND u.community_id = p.community_id
		   AND u.status = 'ACTIVE'
		   AND u.deleted_at IS NULL
		LEFT JOIN invoice_plan_unit_overrides o
			ON o.invoice_plan_id = p.id
		   AND o.unit_id = u.id
		WHERE p.id = $1
		  AND p.organization_id = $2
		  AND p.status = 'ACTIVE'
		  AND p.deleted_at IS NULL
		  AND p.starts_on <= $5
		  AND (p.ends_on IS NULL OR p.ends_on >= $4)
		ON CONFLICT (invoice_plan_id, unit_id, billing_period) WHERE source = 'SCHEDULED'
		DO NOTHING
	`

	tag, err := tx.Exec(ctx, query, plan.ID, plan.OrganizationID, billingPeriod, issuedOn, dueDate)
	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}

func insertManualInvoice(ctx context.Context, db *pgxpool.Pool, input CreateInvoiceInput, issuedOn time.Time, dueDate time.Time) (Invoice, error) {
	const query = `
		INSERT INTO invoices (
			organization_id, community_id, unit_id, invoice_plan_id, source, amount, paid_amount, issued_on, due_date, status, billing_period, description
		)
		VALUES (
			$1, $2, $3, NULL, 'MANUAL', $4, 0.00, $5, $6,
			CASE WHEN $6 < CURRENT_DATE THEN 'OVERDUE' ELSE 'UNPAID' END,
			$7, NULLIF($8, '')
		)
		RETURNING id
	`

	var invoiceID string
	if err := db.QueryRow(
		ctx,
		query,
		input.OrganizationID,
		input.CommunityID,
		input.UnitID,
		input.Amount,
		issuedOn,
		dueDate,
		input.BillingPeriod,
		input.Description,
	).Scan(&invoiceID); err != nil {
		return Invoice{}, err
	}

	return findInvoiceByID(ctx, db, input.OrganizationID, input.CommunityID, invoiceID)
}

func listInvoices(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string) ([]Invoice, error) {
	const query = `
		SELECT
			i.id,
			i.organization_id,
			i.community_id,
			i.unit_id,
			u.unit_number,
			h.id,
			COALESCE(h.name, ''),
			i.invoice_plan_id,
			i.source,
			i.amount,
			i.paid_amount,
			GREATEST(i.amount - i.paid_amount, 0),
			i.issued_on,
			i.due_date,
			i.status,
			i.billing_period,
			COALESCE(i.description, ''),
			i.created_at,
			i.updated_at,
			CASE
				WHEN i.status <> 'PAID' AND i.status <> 'VOID' AND i.due_date < CURRENT_DATE AND i.amount > i.paid_amount
				THEN TRUE
				ELSE FALSE
			END AS is_overdue
		FROM invoices i
		JOIN units u ON u.id = i.unit_id
		LEFT JOIN LATERAL (
			SELECT hh.id, hh.name
			FROM households hh
			WHERE hh.organization_id = i.organization_id
			  AND hh.community_id = i.community_id
			  AND hh.unit_id = i.unit_id
			ORDER BY hh.created_at ASC
			LIMIT 1
		) h ON TRUE
		WHERE i.organization_id = $1
		  AND i.community_id = $2
		ORDER BY i.due_date ASC, i.created_at DESC
	`

	rows, err := db.Query(ctx, query, organizationID, communityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []Invoice
	for rows.Next() {
		var invoice Invoice
		if err := rows.Scan(
			&invoice.ID,
			&invoice.OrganizationID,
			&invoice.CommunityID,
			&invoice.UnitID,
			&invoice.UnitNumber,
			&invoice.HouseholdID,
			&invoice.HouseholdName,
			&invoice.InvoicePlanID,
			&invoice.Source,
			&invoice.Amount,
			&invoice.PaidAmount,
			&invoice.OutstandingAmount,
			&invoice.IssuedOn,
			&invoice.DueDate,
			&invoice.Status,
			&invoice.BillingPeriod,
			&invoice.Description,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
			&invoice.IsOverdue,
		); err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}

	return invoices, rows.Err()
}

func findInvoiceByID(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string, invoiceID string) (Invoice, error) {
	const query = `
		SELECT
			i.id,
			i.organization_id,
			i.community_id,
			i.unit_id,
			u.unit_number,
			h.id,
			COALESCE(h.name, ''),
			i.invoice_plan_id,
			i.source,
			i.amount,
			i.paid_amount,
			GREATEST(i.amount - i.paid_amount, 0),
			i.issued_on,
			i.due_date,
			i.status,
			i.billing_period,
			COALESCE(i.description, ''),
			i.created_at,
			i.updated_at,
			CASE
				WHEN i.status <> 'PAID' AND i.status <> 'VOID' AND i.due_date < CURRENT_DATE AND i.amount > i.paid_amount
				THEN TRUE
				ELSE FALSE
			END AS is_overdue
		FROM invoices i
		JOIN units u ON u.id = i.unit_id
		LEFT JOIN LATERAL (
			SELECT hh.id, hh.name
			FROM households hh
			WHERE hh.organization_id = i.organization_id
			  AND hh.community_id = i.community_id
			  AND hh.unit_id = i.unit_id
			ORDER BY hh.created_at ASC
			LIMIT 1
		) h ON TRUE
		WHERE i.organization_id = $1
		  AND i.community_id = $2
		  AND i.id = $3
	`

	var invoice Invoice
	err := db.QueryRow(ctx, query, organizationID, communityID, invoiceID).Scan(
		&invoice.ID,
		&invoice.OrganizationID,
		&invoice.CommunityID,
		&invoice.UnitID,
		&invoice.UnitNumber,
		&invoice.HouseholdID,
		&invoice.HouseholdName,
		&invoice.InvoicePlanID,
		&invoice.Source,
		&invoice.Amount,
		&invoice.PaidAmount,
		&invoice.OutstandingAmount,
		&invoice.IssuedOn,
		&invoice.DueDate,
		&invoice.Status,
		&invoice.BillingPeriod,
		&invoice.Description,
		&invoice.CreatedAt,
		&invoice.UpdatedAt,
		&invoice.IsOverdue,
	)

	return invoice, err
}

func listOverdueHouseholds(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string) ([]OverdueHouseholdSummary, error) {
	const query = `
		SELECT
			h.id,
			COALESCE(NULLIF(h.name, ''), 'Unit ' || u.unit_number),
			u.id,
			u.unit_number,
			COALESCE(pr.first_name || ' ' || pr.last_name, ''),
			COUNT(i.id) AS overdue_invoice_count,
			SUM(GREATEST(i.amount - i.paid_amount, 0)) AS total_outstanding,
			MIN(i.due_date) AS earliest_due_date
		FROM invoices i
		JOIN units u
			ON u.id = i.unit_id
		LEFT JOIN LATERAL (
			SELECT hh.id, hh.name
			FROM households hh
			WHERE hh.organization_id = i.organization_id
			  AND hh.community_id = i.community_id
			  AND hh.unit_id = i.unit_id
			ORDER BY hh.created_at ASC
			LIMIT 1
		) h ON TRUE
		LEFT JOIN LATERAL (
			SELECT r.first_name, r.last_name
			FROM residents r
			WHERE r.organization_id = i.organization_id
			  AND r.community_id = i.community_id
			  AND r.unit_id = i.unit_id
			  AND r.deleted_at IS NULL
			  AND r.status IN ('ACTIVE', 'PENDING')
			ORDER BY r.is_primary_contact DESC, r.created_at ASC
			LIMIT 1
		) pr ON TRUE
		WHERE i.organization_id = $1
		  AND i.community_id = $2
		  AND i.status <> 'PAID'
		  AND i.status <> 'VOID'
		  AND i.due_date < CURRENT_DATE
		  AND i.amount > i.paid_amount
		GROUP BY h.id, h.name, u.id, u.unit_number, pr.first_name, pr.last_name
		ORDER BY MIN(i.due_date) ASC, SUM(GREATEST(i.amount - i.paid_amount, 0)) DESC
	`

	rows, err := db.Query(ctx, query, organizationID, communityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []OverdueHouseholdSummary
	for rows.Next() {
		var summary OverdueHouseholdSummary
		if err := rows.Scan(
			&summary.HouseholdID,
			&summary.HouseholdName,
			&summary.UnitID,
			&summary.UnitNumber,
			&summary.PrimaryResidentName,
			&summary.OverdueInvoiceCount,
			&summary.TotalOutstanding,
			&summary.EarliestDueDate,
		); err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}

	return summaries, rows.Err()
}

func findLinkedResidentForUser(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string, userID string) (linkedResident, error) {
	const query = `
		SELECT
			r.id,
			r.first_name,
			r.last_name,
			r.status,
			r.unit_id,
			u.unit_number,
			r.household_id,
			COALESCE(h.name, '')
		FROM residents r
		JOIN units u ON u.id = r.unit_id
		LEFT JOIN households h ON h.id = r.household_id
		WHERE r.organization_id = $1
		  AND r.community_id = $2
		  AND r.user_id = $3
		  AND r.deleted_at IS NULL
		  AND r.status IN ('ACTIVE', 'PENDING')
		ORDER BY r.is_primary_contact DESC, r.created_at ASC
		LIMIT 1
	`

	var resident linkedResident
	err := db.QueryRow(ctx, query, organizationID, communityID, userID).Scan(
		&resident.ID,
		&resident.FirstName,
		&resident.LastName,
		&resident.Status,
		&resident.UnitID,
		&resident.UnitNumber,
		&resident.HouseholdID,
		&resident.HouseholdName,
	)

	return resident, err
}

func listInvoicesForUnit(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string, unitID string) ([]Invoice, error) {
	const query = `
		SELECT
			i.id,
			i.organization_id,
			i.community_id,
			i.unit_id,
			u.unit_number,
			h.id,
			COALESCE(h.name, ''),
			i.invoice_plan_id,
			i.source,
			i.amount,
			i.paid_amount,
			GREATEST(i.amount - i.paid_amount, 0),
			i.issued_on,
			i.due_date,
			i.status,
			i.billing_period,
			COALESCE(i.description, ''),
			i.created_at,
			i.updated_at,
			CASE
				WHEN i.status <> 'PAID' AND i.status <> 'VOID' AND i.due_date < CURRENT_DATE AND i.amount > i.paid_amount
				THEN TRUE
				ELSE FALSE
			END AS is_overdue
		FROM invoices i
		JOIN units u ON u.id = i.unit_id
		LEFT JOIN households h ON h.id = (
			SELECT hh.id
			FROM households hh
			WHERE hh.organization_id = i.organization_id
			  AND hh.community_id = i.community_id
			  AND hh.unit_id = i.unit_id
			ORDER BY hh.created_at ASC
			LIMIT 1
		)
		WHERE i.organization_id = $1
		  AND i.community_id = $2
		  AND i.unit_id = $3
		ORDER BY i.due_date ASC, i.created_at DESC
	`

	rows, err := db.Query(ctx, query, organizationID, communityID, unitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []Invoice
	for rows.Next() {
		var invoice Invoice
		if err := rows.Scan(
			&invoice.ID,
			&invoice.OrganizationID,
			&invoice.CommunityID,
			&invoice.UnitID,
			&invoice.UnitNumber,
			&invoice.HouseholdID,
			&invoice.HouseholdName,
			&invoice.InvoicePlanID,
			&invoice.Source,
			&invoice.Amount,
			&invoice.PaidAmount,
			&invoice.OutstandingAmount,
			&invoice.IssuedOn,
			&invoice.DueDate,
			&invoice.Status,
			&invoice.BillingPeriod,
			&invoice.Description,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
			&invoice.IsOverdue,
		); err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}

	return invoices, rows.Err()
}

func communityUnitExists(ctx context.Context, db *pgxpool.Pool, organizationID string, communityID string, unitID string) (bool, error) {
	const query = `
		SELECT 1
		FROM units
		WHERE organization_id = $1
		  AND community_id = $2
		  AND id = $3
		  AND deleted_at IS NULL
	`

	var exists int
	err := db.QueryRow(ctx, query, organizationID, communityID, unitID).Scan(&exists)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return exists == 1, nil
}

func findOrganizationRoleForUser(ctx context.Context, q queryExecutor, organizationID string, userID string) (string, error) {
	const query = `
		SELECT role
		FROM organization_users
		WHERE organization_id = $1 AND user_id = $2 AND status = 'ACTIVE'
	`

	var role string
	err := q.QueryRow(ctx, query, organizationID, userID).Scan(&role)
	return role, err
}

func findCommunityRoleForUser(ctx context.Context, q queryExecutor, organizationID string, communityID string, userID string) (string, error) {
	const query = `
		SELECT role
		FROM community_users
		WHERE organization_id = $1 AND community_id = $2 AND user_id = $3 AND status = 'ACTIVE'
	`

	var role string
	err := q.QueryRow(ctx, query, organizationID, communityID, userID).Scan(&role)
	return role, err
}

func buildPeriodDates(periodStart time.Time, issueDay int, dueDay int) (time.Time, time.Time) {
	location := time.UTC
	issuedOn := time.Date(periodStart.Year(), periodStart.Month(), issueDay, 0, 0, 0, 0, location)
	dueDate := time.Date(periodStart.Year(), periodStart.Month(), dueDay, 0, 0, 0, 0, location)
	return issuedOn, dueDate
}

func computeResidentOverview(resident linkedResident, invoices []Invoice) ResidentInvoiceOverview {
	var (
		nextUpcomingDueDate *time.Time
		fallbackDueDate     *time.Time
		overdueInvoiceCount int
		totalOutstanding    float64
	)
	today := time.Now().UTC().Truncate(24 * time.Hour)

	for i := range invoices {
		invoice := invoices[i]
		if invoice.Status == "VOID" || invoice.OutstandingAmount <= 0 {
			continue
		}

		totalOutstanding += invoice.OutstandingAmount
		if fallbackDueDate == nil || invoice.DueDate.Before(*fallbackDueDate) {
			dueDate := invoice.DueDate
			fallbackDueDate = &dueDate
		}
		if !invoice.DueDate.Before(today) && (nextUpcomingDueDate == nil || invoice.DueDate.Before(*nextUpcomingDueDate)) {
			dueDate := invoice.DueDate
			nextUpcomingDueDate = &dueDate
		}
		if invoice.IsOverdue {
			overdueInvoiceCount++
		}
	}

	nextDueDate := nextUpcomingDueDate
	if nextDueDate == nil {
		nextDueDate = fallbackDueDate
	}

	return ResidentInvoiceOverview{
		ResidentID:          resident.ID,
		ResidentName:        fmt.Sprintf("%s %s", resident.FirstName, resident.LastName),
		ResidentStatus:      resident.Status,
		UnitID:              resident.UnitID,
		UnitNumber:          resident.UnitNumber,
		HouseholdID:         resident.HouseholdID,
		HouseholdName:       resident.HouseholdName,
		NextDueDate:         nextDueDate,
		OverdueInvoiceCount: overdueInvoiceCount,
		TotalOutstanding:    totalOutstanding,
		Invoices:            invoices,
	}
}
