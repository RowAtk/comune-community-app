package residents

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

var validTypes = map[string]struct{}{
	"OWNER":     {},
	"TENANT":    {},
	"DEPENDENT": {},
	"OCCUPANT":  {},
}

var validHouseholdRoles = map[string]struct{}{
	"HOUSEHOLD_ADMIN":  {},
	"HOUSEHOLD_MEMBER": {},
	"HOUSEHOLD_VIEWER": {},
}

var validStatuses = map[string]struct{}{
	"ACTIVE":    {},
	"INACTIVE":  {},
	"PENDING":   {},
	"MOVED_OUT": {},
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
	ErrResidentIDRequired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "resident_id_required",
		Message: "resident id is required",
	}
	ErrUnitIDRequired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "unit_id_required",
		Message: "unit id is required",
	}
	ErrFirstNameRequired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "resident_first_name_required",
		Message: "resident first name is required",
	}
	ErrLastNameRequired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "resident_last_name_required",
		Message: "resident last name is required",
	}
	ErrResidentTypeRequired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "resident_type_required",
		Message: "resident type is required",
	}
	ErrInvalidResidentType = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "invalid_resident_type",
		Message: "invalid resident type",
	}
	ErrInvalidHouseholdRole = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "invalid_household_role",
		Message: "invalid household role",
	}
	ErrInvalidResidentStatus = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "invalid_resident_status",
		Message: "invalid resident status",
	}
	ErrInvalidMoveDates = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "invalid_resident_move_dates",
		Message: "move out date cannot be before move in date",
	}
	ErrResidentNotFound = apperror.Error{
		Kind:    apperror.KindNotFound,
		Code:    "resident_not_found",
		Message: "resident not found",
	}
	ErrResidentUnitMismatch = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "resident_unit_mismatch",
		Message: "unit does not belong to this community",
	}
	ErrResidentHouseholdMismatch = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "resident_household_mismatch",
		Message: "household does not belong to the selected unit",
	}
	ErrPrimaryContactAlreadyAssigned = apperror.Error{
		Kind:    apperror.KindConflict,
		Code:    "resident_primary_contact_exists",
		Message: "a primary contact already exists for this household",
	}
	ErrCreateFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "create_resident_failed",
		Message: "internal server error",
	}
	ErrListFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "list_residents_failed",
		Message: "internal server error",
	}
	ErrGetFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "get_resident_failed",
		Message: "internal server error",
	}
	ErrUpdateFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "update_resident_failed",
		Message: "internal server error",
	}
	ErrDeleteFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "delete_resident_failed",
		Message: "internal server error",
	}
)

func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (Resident, error) {
	input = sanitizeCreateInput(input)
	moveInDate := createMoveInDate(input)
	moveOutDate := createMoveOutDate(input)

	switch {
	case input.OrganizationID == "":
		return Resident{}, ErrOrganizationIDRequired.Wrap(nil)
	case input.CommunityID == "":
		return Resident{}, ErrCommunityIDRequired.Wrap(nil)
	case input.UnitID == "":
		return Resident{}, ErrUnitIDRequired.Wrap(nil)
	case input.FirstName == "":
		return Resident{}, ErrFirstNameRequired.Wrap(nil)
	case input.LastName == "":
		return Resident{}, ErrLastNameRequired.Wrap(nil)
	case input.ResidentType == "":
		return Resident{}, ErrResidentTypeRequired.Wrap(nil)
	case !isValidType(input.ResidentType):
		return Resident{}, ErrInvalidResidentType.Wrap(nil)
	case input.HouseholdRole != "" && !isValidHouseholdRole(input.HouseholdRole):
		return Resident{}, ErrInvalidHouseholdRole.Wrap(nil)
	case !isValidStatus(input.Status):
		return Resident{}, ErrInvalidResidentStatus.Wrap(nil)
	case !moveDatesAreValid(moveInDate, moveOutDate):
		return Resident{}, ErrInvalidMoveDates.Wrap(nil)
	}

	if err := s.validateRelationships(ctx, input.OrganizationID, input.CommunityID, input.UnitID, input.HouseholdID); err != nil {
		return Resident{}, err
	}

	resident, err := insert(ctx, s.db, input, moveInDate, moveOutDate)
	if err != nil {
		if isPrimaryContactConflict(err) {
			return Resident{}, ErrPrimaryContactAlreadyAssigned.Wrap(fmt.Errorf("create resident: %w", err))
		}
		return Resident{}, ErrCreateFailed.Wrap(fmt.Errorf("create resident: %w", err))
	}

	return resident, nil
}

func (s *Service) List(ctx context.Context, organizationID string, communityID string) ([]Resident, error) {
	organizationID = strings.TrimSpace(organizationID)
	communityID = strings.TrimSpace(communityID)

	switch {
	case organizationID == "":
		return nil, ErrOrganizationIDRequired.Wrap(nil)
	case communityID == "":
		return nil, ErrCommunityIDRequired.Wrap(nil)
	}

	residents, err := list(ctx, s.db, organizationID, communityID)
	if err != nil {
		return nil, ErrListFailed.Wrap(fmt.Errorf("list residents: %w", err))
	}

	return ensureSlice(residents), nil
}

func (s *Service) GetByID(ctx context.Context, organizationID string, communityID string, residentID string) (Resident, error) {
	organizationID = strings.TrimSpace(organizationID)
	communityID = strings.TrimSpace(communityID)
	residentID = strings.TrimSpace(residentID)

	switch {
	case organizationID == "":
		return Resident{}, ErrOrganizationIDRequired.Wrap(nil)
	case communityID == "":
		return Resident{}, ErrCommunityIDRequired.Wrap(nil)
	case residentID == "":
		return Resident{}, ErrResidentIDRequired.Wrap(nil)
	}

	resident, err := findByID(ctx, s.db, organizationID, communityID, residentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Resident{}, ErrResidentNotFound.Wrap(nil)
		}
		return Resident{}, ErrGetFailed.Wrap(fmt.Errorf("get resident: %w", err))
	}

	return resident, nil
}

func (s *Service) Update(ctx context.Context, organizationID string, communityID string, residentID string, input UpdateInput) (Resident, error) {
	organizationID = strings.TrimSpace(organizationID)
	communityID = strings.TrimSpace(communityID)
	residentID = strings.TrimSpace(residentID)
	input = sanitizeUpdateInput(input)

	switch {
	case organizationID == "":
		return Resident{}, ErrOrganizationIDRequired.Wrap(nil)
	case communityID == "":
		return Resident{}, ErrCommunityIDRequired.Wrap(nil)
	case residentID == "":
		return Resident{}, ErrResidentIDRequired.Wrap(nil)
	}

	current, err := findByID(ctx, s.db, organizationID, communityID, residentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Resident{}, ErrResidentNotFound.Wrap(nil)
		}
		return Resident{}, ErrUpdateFailed.Wrap(fmt.Errorf("load resident: %w", err))
	}

	if input.FirstName != nil && *input.FirstName == "" {
		return Resident{}, ErrFirstNameRequired.Wrap(nil)
	}
	if input.LastName != nil && *input.LastName == "" {
		return Resident{}, ErrLastNameRequired.Wrap(nil)
	}
	if input.ResidentType != nil {
		if *input.ResidentType == "" {
			return Resident{}, ErrResidentTypeRequired.Wrap(nil)
		}
		if !isValidType(*input.ResidentType) {
			return Resident{}, ErrInvalidResidentType.Wrap(nil)
		}
	}
	if input.HouseholdRole != nil && *input.HouseholdRole != "" && !isValidHouseholdRole(*input.HouseholdRole) {
		return Resident{}, ErrInvalidHouseholdRole.Wrap(nil)
	}
	if input.Status != nil && !isValidStatus(*input.Status) {
		return Resident{}, ErrInvalidResidentStatus.Wrap(nil)
	}

	effectiveMoveInDate := current.MoveInDate
	parsedMoveInDate := updateMoveInDate(input)
	if parsedMoveInDate != nil {
		effectiveMoveInDate = *parsedMoveInDate
	}
	effectiveMoveOutDate := current.MoveOutDate
	parsedMoveOutDate := updateMoveOutDate(input)
	if parsedMoveOutDate != nil {
		effectiveMoveOutDate = *parsedMoveOutDate
	}
	if !moveDatesAreValid(effectiveMoveInDate, effectiveMoveOutDate) {
		return Resident{}, ErrInvalidMoveDates.Wrap(nil)
	}

	effectiveUnitID := current.UnitID
	if input.UnitID != nil {
		if *input.UnitID == "" {
			return Resident{}, ErrUnitIDRequired.Wrap(nil)
		}
		effectiveUnitID = *input.UnitID
	}

	effectiveHouseholdID := stringPtrValue(current.HouseholdID)
	if input.HouseholdID != nil {
		effectiveHouseholdID = strings.TrimSpace(*input.HouseholdID)
	}

	if err := s.validateRelationships(ctx, organizationID, communityID, effectiveUnitID, effectiveHouseholdID); err != nil {
		return Resident{}, err
	}

	resident, err := update(ctx, s.db, organizationID, communityID, residentID, input, parsedMoveInDate, parsedMoveOutDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Resident{}, ErrResidentNotFound.Wrap(nil)
		}
		if isPrimaryContactConflict(err) {
			return Resident{}, ErrPrimaryContactAlreadyAssigned.Wrap(fmt.Errorf("update resident: %w", err))
		}
		return Resident{}, ErrUpdateFailed.Wrap(fmt.Errorf("update resident: %w", err))
	}

	return resident, nil
}

func (s *Service) Delete(ctx context.Context, organizationID string, communityID string, residentID string) error {
	organizationID = strings.TrimSpace(organizationID)
	communityID = strings.TrimSpace(communityID)
	residentID = strings.TrimSpace(residentID)

	switch {
	case organizationID == "":
		return ErrOrganizationIDRequired.Wrap(nil)
	case communityID == "":
		return ErrCommunityIDRequired.Wrap(nil)
	case residentID == "":
		return ErrResidentIDRequired.Wrap(nil)
	}

	err := remove(ctx, s.db, organizationID, communityID, residentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrResidentNotFound.Wrap(nil)
		}
		return ErrDeleteFailed.Wrap(fmt.Errorf("delete resident: %w", err))
	}

	return nil
}

func (s *Service) validateRelationships(ctx context.Context, organizationID string, communityID string, unitID string, householdID string) error {
	ok, err := unitExists(ctx, s.db, organizationID, communityID, unitID)
	if err != nil {
		return ErrCreateFailed.Wrap(fmt.Errorf("validate resident unit: %w", err))
	}
	if !ok {
		return ErrResidentUnitMismatch.Wrap(nil)
	}

	if householdID == "" {
		return nil
	}

	householdUnitID, ok, err := findHouseholdUnitID(ctx, s.db, organizationID, communityID, householdID)
	if err != nil {
		return ErrCreateFailed.Wrap(fmt.Errorf("validate resident household: %w", err))
	}
	if !ok || householdUnitID != unitID {
		return ErrResidentHouseholdMismatch.Wrap(nil)
	}

	return nil
}

func sanitizeCreateInput(input CreateInput) CreateInput {
	input.OrganizationID = strings.TrimSpace(input.OrganizationID)
	input.CommunityID = strings.TrimSpace(input.CommunityID)
	input.UnitID = strings.TrimSpace(input.UnitID)
	input.UserID = strings.TrimSpace(input.UserID)
	input.HouseholdID = strings.TrimSpace(input.HouseholdID)
	input.FirstName = strings.TrimSpace(input.FirstName)
	input.LastName = strings.TrimSpace(input.LastName)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.Phone = strings.TrimSpace(input.Phone)
	input.ResidentType = strings.ToUpper(strings.TrimSpace(input.ResidentType))
	input.HouseholdRole = strings.ToUpper(strings.TrimSpace(input.HouseholdRole))
	input.Status = strings.ToUpper(strings.TrimSpace(input.Status))

	if input.Status == "" {
		input.Status = "ACTIVE"
	}

	return input
}

func sanitizeUpdateInput(input UpdateInput) UpdateInput {
	if input.UnitID != nil {
		value := strings.TrimSpace(*input.UnitID)
		input.UnitID = &value
	}
	if input.UserID != nil {
		value := strings.TrimSpace(*input.UserID)
		input.UserID = &value
	}
	if input.HouseholdID != nil {
		value := strings.TrimSpace(*input.HouseholdID)
		input.HouseholdID = &value
	}
	if input.FirstName != nil {
		value := strings.TrimSpace(*input.FirstName)
		input.FirstName = &value
	}
	if input.LastName != nil {
		value := strings.TrimSpace(*input.LastName)
		input.LastName = &value
	}
	if input.Email != nil {
		value := strings.ToLower(strings.TrimSpace(*input.Email))
		input.Email = &value
	}
	if input.Phone != nil {
		value := strings.TrimSpace(*input.Phone)
		input.Phone = &value
	}
	if input.ResidentType != nil {
		value := strings.ToUpper(strings.TrimSpace(*input.ResidentType))
		input.ResidentType = &value
	}
	if input.HouseholdRole != nil {
		value := strings.ToUpper(strings.TrimSpace(*input.HouseholdRole))
		input.HouseholdRole = &value
	}
	if input.Status != nil {
		value := strings.ToUpper(strings.TrimSpace(*input.Status))
		input.Status = &value
	}

	return input
}

func createMoveInDate(input CreateInput) *time.Time {
	if input.MoveInDate == nil {
		return nil
	}

	return &input.MoveInDate.Time
}

func createMoveOutDate(input CreateInput) *time.Time {
	if input.MoveOutDate == nil {
		return nil
	}

	return &input.MoveOutDate.Time
}

func updateMoveInDate(input UpdateInput) **time.Time {
	if input.MoveInDate == nil {
		return nil
	}

	return &input.MoveInDate.Value
}

func updateMoveOutDate(input UpdateInput) **time.Time {
	if input.MoveOutDate == nil {
		return nil
	}

	return &input.MoveOutDate.Value
}

func isValidType(value string) bool {
	_, ok := validTypes[value]
	return ok
}

func isValidHouseholdRole(value string) bool {
	_, ok := validHouseholdRoles[value]
	return ok
}

func isValidStatus(value string) bool {
	_, ok := validStatuses[value]
	return ok
}

func moveDatesAreValid(moveInDate *time.Time, moveOutDate *time.Time) bool {
	if moveInDate == nil || moveOutDate == nil {
		return true
	}
	return !moveOutDate.Before(*moveInDate)
}

func stringPtrValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func isPrimaryContactConflict(err error) bool {
	return db.IsUniqueViolation(err) && db.IsConstraintViolation(err, "uq_residents_primary_contact_per_household")
}

func ensureSlice[T any](items []T) []T {
	if items == nil {
		return []T{}
	}

	return items
}
