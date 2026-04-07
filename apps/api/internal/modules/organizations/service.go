package organizations

import (
	"comune/apps/api/internal/platform/apperror"
	"comune/apps/api/internal/platform/db"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var validStatuses = map[string]struct{}{
	"ACTIVE":    {},
	"SUSPENDED": {},
	"INACTIVE":  {},
}

var validMembershipRoles = map[string]struct{}{
	"OWNER":         {},
	"ORG_ADMIN":     {},
	"BILLING_ADMIN": {},
	"OPERATIONS":    {},
	"SUPPORT":       {},
}

var validMembershipStatuses = map[string]struct{}{
	"INVITED":   {},
	"ACTIVE":    {},
	"SUSPENDED": {},
	"REMOVED":   {},
}

var (
	ErrOrganizationNameRequired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "organization_name_required",
		Message: "organization name is required",
	}
	ErrOrganizationSlugRequired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "organization_slug_required",
		Message: "organization slug is required",
	}
	ErrInvalidOrganizationStatus = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "invalid_organization_status",
		Message: "invalid organization status",
	}
	ErrOrganizationSlugTaken = apperror.Error{
		Kind:    apperror.KindConflict,
		Code:    "organization_slug_taken",
		Message: "organization slug is already taken",
	}
	ErrOrganizationNotFound = apperror.Error{
		Kind:    apperror.KindNotFound,
		Code:    "organization_not_found",
		Message: "organization not found",
	}
	ErrCreateOrganizationFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "create_organization_failed",
		Message: "internal server error",
	}
	ErrListOrganizationsFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "list_organizations_failed",
		Message: "internal server error",
	}
	ErrGetOrganizationFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "get_organization_failed",
		Message: "internal server error",
	}
	ErrUpdateOrganizationFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "update_organization_failed",
		Message: "internal server error",
	}
	ErrDeleteOrganizationFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "delete_organization_failed",
		Message: "internal server error",
	}
	ErrOrganizationIDRequired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "organization_id_required",
		Message: "organization id is required",
	}
	ErrUserIDRequired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "user_id_required",
		Message: "user id is required",
	}
	ErrMemberRoleRequired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "member_role_required",
		Message: "member role is required",
	}
	ErrInvalidMemberRole = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "invalid_member_role",
		Message: "invalid organization member role",
	}
	ErrInvalidMemberStatus = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "invalid_member_status",
		Message: "invalid organization member status",
	}
	ErrMemberNotFound = apperror.Error{
		Kind:    apperror.KindNotFound,
		Code:    "organization_member_not_found",
		Message: "organization member not found",
	}
	ErrListMembersFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "list_organization_members_failed",
		Message: "internal server error",
	}
	ErrUpdateMemberFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "update_organization_member_failed",
		Message: "internal server error",
	}
	ErrInvitationEmailRequired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "invitation_email_required",
		Message: "invitation email is required",
	}
	ErrInvitationRoleRequired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "invitation_role_required",
		Message: "invitation role is required",
	}
	ErrInvitationNotFound = apperror.Error{
		Kind:    apperror.KindNotFound,
		Code:    "organization_invitation_not_found",
		Message: "organization invitation not found",
	}
	ErrInvitationExpired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "organization_invitation_expired",
		Message: "organization invitation has expired",
	}
	ErrInvitationAlreadyAccepted = apperror.Error{
		Kind:    apperror.KindConflict,
		Code:    "organization_invitation_already_accepted",
		Message: "organization invitation has already been accepted",
	}
	ErrInvitationEmailMismatch = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "organization_invitation_email_mismatch",
		Message: "invitation email does not match the authenticated user",
	}
	ErrCreateInvitationFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "create_organization_invitation_failed",
		Message: "internal server error",
	}
	ErrListInvitationsFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "list_organization_invitations_failed",
		Message: "internal server error",
	}
	ErrAcceptInvitationFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "accept_organization_invitation_failed",
		Message: "internal server error",
	}
)

func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

func (s *Service) Create(ctx context.Context, input CreateOrganizationInput) (Organization, error) {
	input = sanitizeCreateInput(input)

	switch {
	case input.Name == "":
		return Organization{}, ErrOrganizationNameRequired.Wrap(nil)
	case input.Slug == "":
		return Organization{}, ErrOrganizationSlugRequired.Wrap(nil)
	case !isValidStatus(input.Status):
		return Organization{}, ErrInvalidOrganizationStatus.Wrap(nil)
	}

	org, err := insertOrganization(ctx, s.db, input)
	if err != nil {
		if db.IsUniqueViolation(err) {
			return Organization{}, ErrOrganizationSlugTaken.Wrap(fmt.Errorf("create organization: %w", err))
		}
		return Organization{}, ErrCreateOrganizationFailed.Wrap(fmt.Errorf("create organization: %w", err))
	}

	return org, nil
}

func (s *Service) List(ctx context.Context) ([]Organization, error) {
	organizations, err := listOrganizations(ctx, s.db)
	if err != nil {
		return nil, ErrListOrganizationsFailed.Wrap(fmt.Errorf("list organizations: %w", err))
	}

	return organizations, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (Organization, error) {
	org, err := findOrganizationByID(ctx, s.db, strings.TrimSpace(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Organization{}, ErrOrganizationNotFound.Wrap(nil)
		}
		return Organization{}, ErrGetOrganizationFailed.Wrap(fmt.Errorf("get organization by id: %w", err))
	}

	return org, nil
}

func (s *Service) Update(ctx context.Context, id string, input UpdateOrganizationInput) (Organization, error) {
	id = strings.TrimSpace(id)
	input = sanitizeUpdateInput(input)

	if input.Name != nil && *input.Name == "" {
		return Organization{}, ErrOrganizationNameRequired.Wrap(nil)
	}
	if input.Slug != nil && *input.Slug == "" {
		return Organization{}, ErrOrganizationSlugRequired.Wrap(nil)
	}
	if input.Status != nil && !isValidStatus(*input.Status) {
		return Organization{}, ErrInvalidOrganizationStatus.Wrap(nil)
	}

	org, err := updateOrganization(ctx, s.db, id, input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Organization{}, ErrOrganizationNotFound.Wrap(nil)
		}
		if db.IsUniqueViolation(err) {
			return Organization{}, ErrOrganizationSlugTaken.Wrap(fmt.Errorf("update organization: %w", err))
		}
		return Organization{}, ErrUpdateOrganizationFailed.Wrap(fmt.Errorf("update organization: %w", err))
	}

	return org, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	err := deleteOrganization(ctx, s.db, strings.TrimSpace(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrOrganizationNotFound.Wrap(nil)
		}
		return ErrDeleteOrganizationFailed.Wrap(fmt.Errorf("delete organization: %w", err))
	}

	return nil
}

func (s *Service) ListMembers(ctx context.Context, organizationID string) ([]OrganizationMember, error) {
	organizationID = strings.TrimSpace(organizationID)
	if organizationID == "" {
		return nil, ErrOrganizationIDRequired.Wrap(nil)
	}

	members, err := listOrganizationMembers(ctx, s.db, organizationID)
	if err != nil {
		return nil, ErrListMembersFailed.Wrap(fmt.Errorf("list organization members: %w", err))
	}

	return members, nil
}

func (s *Service) UpdateMember(ctx context.Context, organizationID string, userID string, input UpdateOrganizationMemberInput) (OrganizationMember, error) {
	organizationID = strings.TrimSpace(organizationID)
	userID = strings.TrimSpace(userID)
	input = sanitizeMemberInput(input)

	switch {
	case organizationID == "":
		return OrganizationMember{}, ErrOrganizationIDRequired.Wrap(nil)
	case userID == "":
		return OrganizationMember{}, ErrUserIDRequired.Wrap(nil)
	case input.Role == nil && input.Status == nil:
		return OrganizationMember{}, ErrMemberRoleRequired.Wrap(nil)
	}

	if input.Role != nil && !isValidMembershipRole(*input.Role) {
		return OrganizationMember{}, ErrInvalidMemberRole.Wrap(nil)
	}
	if input.Status != nil && !isValidMembershipStatus(*input.Status) {
		return OrganizationMember{}, ErrInvalidMemberStatus.Wrap(nil)
	}

	member, err := updateOrganizationMember(ctx, s.db, organizationID, userID, input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return OrganizationMember{}, ErrMemberNotFound.Wrap(nil)
		}
		return OrganizationMember{}, ErrUpdateMemberFailed.Wrap(fmt.Errorf("update organization member: %w", err))
	}

	return member, nil
}

func (s *Service) CreateInvitation(ctx context.Context, organizationID string, invitedBy string, input CreateOrganizationInvitationInput) (OrganizationInvitation, error) {
	organizationID = strings.TrimSpace(organizationID)
	invitedBy = strings.TrimSpace(invitedBy)
	input = sanitizeInvitationInput(input)

	switch {
	case organizationID == "":
		return OrganizationInvitation{}, ErrOrganizationIDRequired.Wrap(nil)
	case input.Email == "":
		return OrganizationInvitation{}, ErrInvitationEmailRequired.Wrap(nil)
	case input.Role == "":
		return OrganizationInvitation{}, ErrInvitationRoleRequired.Wrap(nil)
	case !isValidMembershipRole(input.Role):
		return OrganizationInvitation{}, ErrInvalidMemberRole.Wrap(nil)
	}

	if input.ExpiresAt == nil {
		expiresAt := time.Now().UTC().Add(7 * 24 * time.Hour)
		input.ExpiresAt = &expiresAt
	}

	invitation, err := insertOrganizationInvitation(ctx, s.db, organizationID, invitedBy, input)
	if err != nil {
		return OrganizationInvitation{}, ErrCreateInvitationFailed.Wrap(fmt.Errorf("create organization invitation: %w", err))
	}

	return invitation, nil
}

func (s *Service) ListInvitations(ctx context.Context, organizationID string) ([]OrganizationInvitation, error) {
	organizationID = strings.TrimSpace(organizationID)
	if organizationID == "" {
		return nil, ErrOrganizationIDRequired.Wrap(nil)
	}

	invitations, err := listOrganizationInvitations(ctx, s.db, organizationID)
	if err != nil {
		return nil, ErrListInvitationsFailed.Wrap(fmt.Errorf("list organization invitations: %w", err))
	}

	return invitations, nil
}

func (s *Service) AcceptInvitation(ctx context.Context, params acceptOrganizationInvitationParams) (OrganizationMember, error) {
	params.UserID = strings.TrimSpace(params.UserID)
	params.Email = strings.ToLower(strings.TrimSpace(params.Email))
	params.Token = strings.TrimSpace(params.Token)

	switch {
	case params.UserID == "":
		return OrganizationMember{}, ErrUserIDRequired.Wrap(nil)
	case params.Email == "":
		return OrganizationMember{}, ErrInvitationEmailRequired.Wrap(nil)
	case params.Token == "":
		return OrganizationMember{}, ErrInvitationNotFound.Wrap(nil)
	}

	var member OrganizationMember
	err := db.RunInTx(ctx, s.db, func(tx pgx.Tx) error {
		invitation, err := findOrganizationInvitationByToken(ctx, tx, params.Token)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrInvitationNotFound.Wrap(nil)
			}
			return ErrAcceptInvitationFailed.Wrap(fmt.Errorf("find organization invitation: %w", err))
		}

		switch {
		case invitation.AcceptedAt != nil:
			return ErrInvitationAlreadyAccepted.Wrap(nil)
		case time.Now().UTC().After(invitation.ExpiresAt):
			return ErrInvitationExpired.Wrap(nil)
		case strings.ToLower(strings.TrimSpace(invitation.Email)) != params.Email:
			return ErrInvitationEmailMismatch.Wrap(nil)
		}

		if err := acceptOrganizationInvitation(ctx, tx, invitation, params.UserID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrInvitationAlreadyAccepted.Wrap(nil)
			}
			return ErrAcceptInvitationFailed.Wrap(fmt.Errorf("accept organization invitation: %w", err))
		}

		member, err = findOrganizationMemberByID(ctx, tx, invitation.OrganizationID, params.UserID)
		if err != nil {
			return ErrAcceptInvitationFailed.Wrap(fmt.Errorf("load organization member: %w", err))
		}

		return nil
	})
	if err != nil {
		return OrganizationMember{}, err
	}

	return member, nil
}

func sanitizeCreateInput(input CreateOrganizationInput) CreateOrganizationInput {
	input.Name = strings.TrimSpace(input.Name)
	input.Slug = strings.ToLower(strings.TrimSpace(input.Slug))
	input.LegalName = strings.TrimSpace(input.LegalName)
	input.BillingEmail = strings.ToLower(strings.TrimSpace(input.BillingEmail))
	input.Phone = strings.TrimSpace(input.Phone)
	input.CountryCode = strings.ToUpper(strings.TrimSpace(input.CountryCode))
	input.Timezone = strings.TrimSpace(input.Timezone)
	input.Status = strings.ToUpper(strings.TrimSpace(input.Status))

	if input.CountryCode == "" {
		input.CountryCode = "JM"
	}
	if input.Timezone == "" {
		input.Timezone = "America/Jamaica"
	}
	if input.Status == "" {
		input.Status = "ACTIVE"
	}
	if !json.Valid(jsonbOrEmpty(input.Settings)) {
		input.Settings = json.RawMessage("{}")
	}

	return input
}

func sanitizeUpdateInput(input UpdateOrganizationInput) UpdateOrganizationInput {
	if input.Name != nil {
		value := strings.TrimSpace(*input.Name)
		input.Name = &value
	}
	if input.Slug != nil {
		value := strings.ToLower(strings.TrimSpace(*input.Slug))
		input.Slug = &value
	}
	if input.LegalName != nil {
		value := strings.TrimSpace(*input.LegalName)
		input.LegalName = &value
	}
	if input.BillingEmail != nil {
		value := strings.ToLower(strings.TrimSpace(*input.BillingEmail))
		input.BillingEmail = &value
	}
	if input.Phone != nil {
		value := strings.TrimSpace(*input.Phone)
		input.Phone = &value
	}
	if input.CountryCode != nil {
		value := strings.ToUpper(strings.TrimSpace(*input.CountryCode))
		if value == "" {
			value = "JM"
		}
		input.CountryCode = &value
	}
	if input.Timezone != nil {
		value := strings.TrimSpace(*input.Timezone)
		if value == "" {
			value = "America/Jamaica"
		}
		input.Timezone = &value
	}
	if input.Status != nil {
		value := strings.ToUpper(strings.TrimSpace(*input.Status))
		input.Status = &value
	}
	if input.Settings != nil && !json.Valid(jsonbOrEmpty(*input.Settings)) {
		value := json.RawMessage("{}")
		input.Settings = &value
	}

	return input
}

func sanitizeMemberInput(input UpdateOrganizationMemberInput) UpdateOrganizationMemberInput {
	if input.Role != nil {
		value := strings.ToUpper(strings.TrimSpace(*input.Role))
		input.Role = &value
	}
	if input.Status != nil {
		value := strings.ToUpper(strings.TrimSpace(*input.Status))
		input.Status = &value
	}

	return input
}

func sanitizeInvitationInput(input CreateOrganizationInvitationInput) CreateOrganizationInvitationInput {
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.Role = strings.ToUpper(strings.TrimSpace(input.Role))
	return input
}

func isValidStatus(status string) bool {
	_, ok := validStatuses[status]
	return ok
}

func isValidMembershipRole(role string) bool {
	_, ok := validMembershipRoles[role]
	return ok
}

func isValidMembershipStatus(status string) bool {
	_, ok := validMembershipStatuses[status]
	return ok
}
