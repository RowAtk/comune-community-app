package units

import (
	"comune/apps/api/internal/platform/apperror"
	"comune/apps/api/internal/platform/db"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var validStatuses = map[string]struct{}{
	"ACTIVE":   {},
	"INACTIVE": {},
}

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
	ErrUnitNumberRequired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "unit_number_required",
		Message: "unit number is required",
	}
	ErrInvalidStatus = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "invalid_unit_status",
		Message: "invalid unit status",
	}
	ErrUnitNumberTaken = apperror.Error{
		Kind:    apperror.KindConflict,
		Code:    "unit_number_taken",
		Message: "unit number is already taken for this community",
	}
	ErrUnitNotFound = apperror.Error{
		Kind:    apperror.KindNotFound,
		Code:    "unit_not_found",
		Message: "unit not found",
	}
	ErrCreateFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "create_unit_failed",
		Message: "internal server error",
	}
	ErrListFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "list_units_failed",
		Message: "internal server error",
	}
	ErrGetFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "get_unit_failed",
		Message: "internal server error",
	}
	ErrUpdateFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "update_unit_failed",
		Message: "internal server error",
	}
	ErrDeleteFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "delete_unit_failed",
		Message: "internal server error",
	}
)

func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (Unit, error) {
	input = sanitizeCreateInput(input)

	switch {
	case input.OrganizationID == "":
		return Unit{}, ErrOrganizationIDRequired.Wrap(nil)
	case input.CommunityID == "":
		return Unit{}, ErrCommunityIDRequired.Wrap(nil)
	case input.UnitNumber == "":
		return Unit{}, ErrUnitNumberRequired.Wrap(nil)
	case !isValidStatus(input.Status):
		return Unit{}, ErrInvalidStatus.Wrap(nil)
	}

	unit, err := insert(ctx, s.db, input)
	if err != nil {
		if db.IsUniqueViolation(err) && db.IsConstraintViolation(err, "uq_units_community_unit_number") {
			return Unit{}, ErrUnitNumberTaken.Wrap(fmt.Errorf("create unit: %w", err))
		}
		return Unit{}, ErrCreateFailed.Wrap(fmt.Errorf("create unit: %w", err))
	}

	return unit, nil
}

func (s *Service) List(ctx context.Context, organizationID string, communityID string) ([]Unit, error) {
	organizationID = strings.TrimSpace(organizationID)
	communityID = strings.TrimSpace(communityID)

	switch {
	case organizationID == "":
		return nil, ErrOrganizationIDRequired.Wrap(nil)
	case communityID == "":
		return nil, ErrCommunityIDRequired.Wrap(nil)
	}

	units, err := list(ctx, s.db, organizationID, communityID)
	if err != nil {
		return nil, ErrListFailed.Wrap(fmt.Errorf("list units: %w", err))
	}

	return ensureSlice(units), nil
}

func (s *Service) GetByID(ctx context.Context, organizationID string, communityID string, unitID string) (Unit, error) {
	organizationID = strings.TrimSpace(organizationID)
	communityID = strings.TrimSpace(communityID)
	unitID = strings.TrimSpace(unitID)

	switch {
	case organizationID == "":
		return Unit{}, ErrOrganizationIDRequired.Wrap(nil)
	case communityID == "":
		return Unit{}, ErrCommunityIDRequired.Wrap(nil)
	case unitID == "":
		return Unit{}, ErrUnitIDRequired.Wrap(nil)
	}

	unit, err := findByID(ctx, s.db, organizationID, communityID, unitID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Unit{}, ErrUnitNotFound.Wrap(nil)
		}
		return Unit{}, ErrGetFailed.Wrap(fmt.Errorf("get unit: %w", err))
	}

	return unit, nil
}

func (s *Service) Update(ctx context.Context, organizationID string, communityID string, unitID string, input UpdateInput) (Unit, error) {
	organizationID = strings.TrimSpace(organizationID)
	communityID = strings.TrimSpace(communityID)
	unitID = strings.TrimSpace(unitID)
	input = sanitizeUpdateInput(input)

	switch {
	case organizationID == "":
		return Unit{}, ErrOrganizationIDRequired.Wrap(nil)
	case communityID == "":
		return Unit{}, ErrCommunityIDRequired.Wrap(nil)
	case unitID == "":
		return Unit{}, ErrUnitIDRequired.Wrap(nil)
	}

	if input.UnitNumber != nil && *input.UnitNumber == "" {
		return Unit{}, ErrUnitNumberRequired.Wrap(nil)
	}
	if input.Status != nil && !isValidStatus(*input.Status) {
		return Unit{}, ErrInvalidStatus.Wrap(nil)
	}

	unit, err := update(ctx, s.db, organizationID, communityID, unitID, input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Unit{}, ErrUnitNotFound.Wrap(nil)
		}
		if db.IsUniqueViolation(err) && db.IsConstraintViolation(err, "uq_units_community_unit_number") {
			return Unit{}, ErrUnitNumberTaken.Wrap(fmt.Errorf("update unit: %w", err))
		}
		return Unit{}, ErrUpdateFailed.Wrap(fmt.Errorf("update unit: %w", err))
	}

	return unit, nil
}

func (s *Service) Delete(ctx context.Context, organizationID string, communityID string, unitID string) error {
	organizationID = strings.TrimSpace(organizationID)
	communityID = strings.TrimSpace(communityID)
	unitID = strings.TrimSpace(unitID)

	switch {
	case organizationID == "":
		return ErrOrganizationIDRequired.Wrap(nil)
	case communityID == "":
		return ErrCommunityIDRequired.Wrap(nil)
	case unitID == "":
		return ErrUnitIDRequired.Wrap(nil)
	}

	err := remove(ctx, s.db, organizationID, communityID, unitID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUnitNotFound.Wrap(nil)
		}
		return ErrDeleteFailed.Wrap(fmt.Errorf("delete unit: %w", err))
	}

	return nil
}

func sanitizeCreateInput(input CreateInput) CreateInput {
	input.OrganizationID = strings.TrimSpace(input.OrganizationID)
	input.CommunityID = strings.TrimSpace(input.CommunityID)
	input.UnitNumber = strings.TrimSpace(input.UnitNumber)
	input.BlockFloor = strings.TrimSpace(input.BlockFloor)
	input.UnitType = strings.TrimSpace(input.UnitType)
	input.Status = strings.ToUpper(strings.TrimSpace(input.Status))

	if input.Status == "" {
		input.Status = "ACTIVE"
	}

	return input
}

func sanitizeUpdateInput(input UpdateInput) UpdateInput {
	if input.UnitNumber != nil {
		value := strings.TrimSpace(*input.UnitNumber)
		input.UnitNumber = &value
	}
	if input.BlockFloor != nil {
		value := strings.TrimSpace(*input.BlockFloor)
		input.BlockFloor = &value
	}
	if input.UnitType != nil {
		value := strings.TrimSpace(*input.UnitType)
		input.UnitType = &value
	}
	if input.Status != nil {
		value := strings.ToUpper(strings.TrimSpace(*input.Status))
		input.Status = &value
	}

	return input
}

func isValidStatus(status string) bool {
	_, ok := validStatuses[status]
	return ok
}

func ensureSlice[T any](items []T) []T {
	if items == nil {
		return []T{}
	}

	return items
}
