package households

import (
	"comune/apps/api/internal/platform/apperror"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrOrganizationIDRequired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "organization_id_required",
		Message: "organization id is required",
	}
	ErrCommunityIDRequired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "community_id_required",
		Message: "community id is required",
	}
	ErrUnitIDRequired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "unit_id_required",
		Message: "unit id is required",
	}
	ErrHouseholdIDRequired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "household_id_required",
		Message: "household id is required",
	}
	ErrHouseholdNameRequired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "household_name_required",
		Message: "household name is required",
	}
	ErrHouseholdNotFound = apperror.Error{
		Kind:    apperror.KindNotFound,
		Code:    "household_not_found",
		Message: "household not found",
	}
	ErrHouseholdUnitMismatch = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "household_unit_mismatch",
		Message: "unit does not belong to this community",
	}
	ErrCreateFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "create_household_failed",
		Message: "internal server error",
	}
	ErrListFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "list_households_failed",
		Message: "internal server error",
	}
	ErrGetFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "get_household_failed",
		Message: "internal server error",
	}
	ErrUpdateFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "update_household_failed",
		Message: "internal server error",
	}
	ErrDeleteFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "delete_household_failed",
		Message: "internal server error",
	}
)

func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (Household, error) {
	input = sanitizeCreateInput(input)

	switch {
	case input.OrganizationID == "":
		return Household{}, ErrOrganizationIDRequired.Wrap(nil)
	case input.CommunityID == "":
		return Household{}, ErrCommunityIDRequired.Wrap(nil)
	case input.UnitID == "":
		return Household{}, ErrUnitIDRequired.Wrap(nil)
	case input.Name == "":
		return Household{}, ErrHouseholdNameRequired.Wrap(nil)
	}

	if ok, err := unitExists(ctx, s.db, input.OrganizationID, input.CommunityID, input.UnitID); err != nil {
		return Household{}, ErrCreateFailed.Wrap(fmt.Errorf("validate household unit: %w", err))
	} else if !ok {
		return Household{}, ErrHouseholdUnitMismatch.Wrap(nil)
	}

	household, err := insert(ctx, s.db, input)
	if err != nil {
		return Household{}, ErrCreateFailed.Wrap(fmt.Errorf("create household: %w", err))
	}

	return household, nil
}

func (s *Service) List(ctx context.Context, organizationID string, communityID string) ([]Household, error) {
	organizationID = strings.TrimSpace(organizationID)
	communityID = strings.TrimSpace(communityID)

	switch {
	case organizationID == "":
		return nil, ErrOrganizationIDRequired.Wrap(nil)
	case communityID == "":
		return nil, ErrCommunityIDRequired.Wrap(nil)
	}

	households, err := list(ctx, s.db, organizationID, communityID)
	if err != nil {
		return nil, ErrListFailed.Wrap(fmt.Errorf("list households: %w", err))
	}

	return ensureSlice(households), nil
}

func (s *Service) GetByID(ctx context.Context, organizationID string, communityID string, householdID string) (Household, error) {
	organizationID = strings.TrimSpace(organizationID)
	communityID = strings.TrimSpace(communityID)
	householdID = strings.TrimSpace(householdID)

	switch {
	case organizationID == "":
		return Household{}, ErrOrganizationIDRequired.Wrap(nil)
	case communityID == "":
		return Household{}, ErrCommunityIDRequired.Wrap(nil)
	case householdID == "":
		return Household{}, ErrHouseholdIDRequired.Wrap(nil)
	}

	household, err := findByID(ctx, s.db, organizationID, communityID, householdID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Household{}, ErrHouseholdNotFound.Wrap(nil)
		}
		return Household{}, ErrGetFailed.Wrap(fmt.Errorf("get household: %w", err))
	}

	return household, nil
}

func (s *Service) Update(ctx context.Context, organizationID string, communityID string, householdID string, input UpdateInput) (Household, error) {
	organizationID = strings.TrimSpace(organizationID)
	communityID = strings.TrimSpace(communityID)
	householdID = strings.TrimSpace(householdID)
	input = sanitizeUpdateInput(input)

	switch {
	case organizationID == "":
		return Household{}, ErrOrganizationIDRequired.Wrap(nil)
	case communityID == "":
		return Household{}, ErrCommunityIDRequired.Wrap(nil)
	case householdID == "":
		return Household{}, ErrHouseholdIDRequired.Wrap(nil)
	}

	if input.Name != nil && *input.Name == "" {
		return Household{}, ErrHouseholdNameRequired.Wrap(nil)
	}
	if input.UnitID != nil {
		if *input.UnitID == "" {
			return Household{}, ErrUnitIDRequired.Wrap(nil)
		}
		if ok, err := unitExists(ctx, s.db, organizationID, communityID, *input.UnitID); err != nil {
			return Household{}, ErrUpdateFailed.Wrap(fmt.Errorf("validate household unit: %w", err))
		} else if !ok {
			return Household{}, ErrHouseholdUnitMismatch.Wrap(nil)
		}
	}

	household, err := update(ctx, s.db, organizationID, communityID, householdID, input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Household{}, ErrHouseholdNotFound.Wrap(nil)
		}
		return Household{}, ErrUpdateFailed.Wrap(fmt.Errorf("update household: %w", err))
	}

	return household, nil
}

func (s *Service) Delete(ctx context.Context, organizationID string, communityID string, householdID string) error {
	organizationID = strings.TrimSpace(organizationID)
	communityID = strings.TrimSpace(communityID)
	householdID = strings.TrimSpace(householdID)

	switch {
	case organizationID == "":
		return ErrOrganizationIDRequired.Wrap(nil)
	case communityID == "":
		return ErrCommunityIDRequired.Wrap(nil)
	case householdID == "":
		return ErrHouseholdIDRequired.Wrap(nil)
	}

	err := remove(ctx, s.db, organizationID, communityID, householdID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrHouseholdNotFound.Wrap(nil)
		}
		return ErrDeleteFailed.Wrap(fmt.Errorf("delete household: %w", err))
	}

	return nil
}

func sanitizeCreateInput(input CreateInput) CreateInput {
	input.OrganizationID = strings.TrimSpace(input.OrganizationID)
	input.CommunityID = strings.TrimSpace(input.CommunityID)
	input.UnitID = strings.TrimSpace(input.UnitID)
	input.Name = strings.TrimSpace(input.Name)
	return input
}

func sanitizeUpdateInput(input UpdateInput) UpdateInput {
	if input.UnitID != nil {
		value := strings.TrimSpace(*input.UnitID)
		input.UnitID = &value
	}
	if input.Name != nil {
		value := strings.TrimSpace(*input.Name)
		input.Name = &value
	}
	return input
}

func ensureSlice[T any](items []T) []T {
	if items == nil {
		return []T{}
	}

	return items
}
