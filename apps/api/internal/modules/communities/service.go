package communities

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
	"ACTIVE":   {},
	"INACTIVE": {},
	"ARCHIVED": {},
}

var validMembershipRoles = map[string]struct{}{
	"COMMUNITY_ADMIN":   {},
	"RESIDENT":          {},
	"SECURITY":          {},
	"MANAGER":           {},
	"BOARD_MEMBER":      {},
	"MAINTENANCE_STAFF": {},
}

var validMembershipStatuses = map[string]struct{}{
	"INVITED":   {},
	"ACTIVE":    {},
	"SUSPENDED": {},
	"REMOVED":   {},
}

var (
	ErrOrganizationIDRequired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "organization_id_required",
		Message: "organization id is required",
	}
	ErrCommunityNameRequired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "community_name_required",
		Message: "community name is required",
	}
	ErrCommunitySlugRequired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "community_slug_required",
		Message: "community slug is required",
	}
	ErrInvalidCommunityStatus = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "invalid_community_status",
		Message: "invalid community status",
	}
	ErrCommunitySlugTaken = apperror.Error{
		Kind:    apperror.KindConflict,
		Code:    "community_slug_taken",
		Message: "community slug is already taken for this organization",
	}
	ErrCommunityNameTaken = apperror.Error{
		Kind:    apperror.KindConflict,
		Code:    "community_name_taken",
		Message: "community name is already taken for this organization",
	}
	ErrCommunityNotFound = apperror.Error{
		Kind:    apperror.KindNotFound,
		Code:    "community_not_found",
		Message: "community not found",
	}
	ErrCreateCommunityFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "create_community_failed",
		Message: "internal server error",
	}
	ErrListCommunitiesFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "list_communities_failed",
		Message: "internal server error",
	}
	ErrGetCommunityFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "get_community_failed",
		Message: "internal server error",
	}
	ErrUpdateCommunityFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "update_community_failed",
		Message: "internal server error",
	}
	ErrDeleteCommunityFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "delete_community_failed",
		Message: "internal server error",
	}
	ErrCommunityIDRequired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "community_id_required",
		Message: "community id is required",
	}
	ErrUserIDRequired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "user_id_required",
		Message: "user id is required",
	}
	ErrMemberRoleRequired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "community_member_role_required",
		Message: "community member role is required",
	}
	ErrInvalidMemberRole = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "invalid_community_member_role",
		Message: "invalid community member role",
	}
	ErrInvalidMemberStatus = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "invalid_community_member_status",
		Message: "invalid community member status",
	}
	ErrMemberNotFound = apperror.Error{
		Kind:    apperror.KindNotFound,
		Code:    "community_member_not_found",
		Message: "community member not found",
	}
	ErrListMembersFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "list_community_members_failed",
		Message: "internal server error",
	}
	ErrUpdateMemberFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "update_community_member_failed",
		Message: "internal server error",
	}
	ErrInvitationEmailRequired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "community_invitation_email_required",
		Message: "invitation email is required",
	}
	ErrInvitationRoleRequired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "community_invitation_role_required",
		Message: "invitation role is required",
	}
	ErrInvitationNotFound = apperror.Error{
		Kind:    apperror.KindNotFound,
		Code:    "community_invitation_not_found",
		Message: "community invitation not found",
	}
	ErrInvitationExpired = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "community_invitation_expired",
		Message: "community invitation has expired",
	}
	ErrInvitationAlreadyAccepted = apperror.Error{
		Kind:    apperror.KindConflict,
		Code:    "community_invitation_already_accepted",
		Message: "community invitation has already been accepted",
	}
	ErrInvitationEmailMismatch = apperror.Error{
		Kind:    apperror.KindValidation,
		Code:    "community_invitation_email_mismatch",
		Message: "invitation email does not match the authenticated user",
	}
	ErrCreateInvitationFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "create_community_invitation_failed",
		Message: "internal server error",
	}
	ErrListInvitationsFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "list_community_invitations_failed",
		Message: "internal server error",
	}
	ErrAcceptInvitationFailed = apperror.Error{
		Kind:    apperror.KindInternal,
		Code:    "accept_community_invitation_failed",
		Message: "internal server error",
	}
)

func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

func (s *Service) Create(ctx context.Context, input CreateCommunityInput) (Community, error) {
	input = sanitizeCreateInput(input)

	switch {
	case input.OrganizationID == "":
		return Community{}, ErrOrganizationIDRequired.Wrap(nil)
	case input.Name == "":
		return Community{}, ErrCommunityNameRequired.Wrap(nil)
	case input.Slug == "":
		return Community{}, ErrCommunitySlugRequired.Wrap(nil)
	case !isValidStatus(input.Status):
		return Community{}, ErrInvalidCommunityStatus.Wrap(nil)
	}

	community, err := insertCommunity(ctx, s.db, input)
	if err != nil {
		if db.IsUniqueViolation(err) {
			switch {
			case db.IsConstraintViolation(err, "uq_communities_org_slug"):
				return Community{}, ErrCommunitySlugTaken.Wrap(fmt.Errorf("create community: %w", err))
			case db.IsConstraintViolation(err, "uq_communities_org_name"):
				return Community{}, ErrCommunityNameTaken.Wrap(fmt.Errorf("create community: %w", err))
			}
		}
		return Community{}, ErrCreateCommunityFailed.Wrap(fmt.Errorf("create community: %w", err))
	}

	return community, nil
}

func (s *Service) List(ctx context.Context, organizationID string) ([]Community, error) {
	organizationID = strings.TrimSpace(organizationID)
	if organizationID == "" {
		return nil, ErrOrganizationIDRequired.Wrap(nil)
	}

	communities, err := listCommunities(ctx, s.db, organizationID)
	if err != nil {
		return nil, ErrListCommunitiesFailed.Wrap(fmt.Errorf("list communities: %w", err))
	}

	return communities, nil
}

func (s *Service) GetByID(ctx context.Context, organizationID string, id string) (Community, error) {
	organizationID = strings.TrimSpace(organizationID)
	id = strings.TrimSpace(id)

	if organizationID == "" {
		return Community{}, ErrOrganizationIDRequired.Wrap(nil)
	}

	community, err := findCommunityByID(ctx, s.db, organizationID, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Community{}, ErrCommunityNotFound.Wrap(nil)
		}
		return Community{}, ErrGetCommunityFailed.Wrap(fmt.Errorf("get community by id: %w", err))
	}

	return community, nil
}

func (s *Service) Update(ctx context.Context, organizationID string, id string, input UpdateCommunityInput) (Community, error) {
	organizationID = strings.TrimSpace(organizationID)
	id = strings.TrimSpace(id)
	input = sanitizeUpdateInput(input)

	if organizationID == "" {
		return Community{}, ErrOrganizationIDRequired.Wrap(nil)
	}
	if input.Name != nil && *input.Name == "" {
		return Community{}, ErrCommunityNameRequired.Wrap(nil)
	}
	if input.Slug != nil && *input.Slug == "" {
		return Community{}, ErrCommunitySlugRequired.Wrap(nil)
	}
	if input.Status != nil && !isValidStatus(*input.Status) {
		return Community{}, ErrInvalidCommunityStatus.Wrap(nil)
	}

	community, err := updateCommunity(ctx, s.db, organizationID, id, input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Community{}, ErrCommunityNotFound.Wrap(nil)
		}
		if db.IsUniqueViolation(err) {
			switch {
			case db.IsConstraintViolation(err, "uq_communities_org_slug"):
				return Community{}, ErrCommunitySlugTaken.Wrap(fmt.Errorf("update community: %w", err))
			case db.IsConstraintViolation(err, "uq_communities_org_name"):
				return Community{}, ErrCommunityNameTaken.Wrap(fmt.Errorf("update community: %w", err))
			}
		}
		return Community{}, ErrUpdateCommunityFailed.Wrap(fmt.Errorf("update community: %w", err))
	}

	return community, nil
}

func (s *Service) Delete(ctx context.Context, organizationID string, id string) error {
	organizationID = strings.TrimSpace(organizationID)
	id = strings.TrimSpace(id)

	if organizationID == "" {
		return ErrOrganizationIDRequired.Wrap(nil)
	}

	err := deleteCommunity(ctx, s.db, organizationID, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrCommunityNotFound.Wrap(nil)
		}
		return ErrDeleteCommunityFailed.Wrap(fmt.Errorf("delete community: %w", err))
	}

	return nil
}

func (s *Service) ListMembers(ctx context.Context, organizationID string, communityID string) ([]CommunityMember, error) {
	organizationID = strings.TrimSpace(organizationID)
	communityID = strings.TrimSpace(communityID)
	if organizationID == "" {
		return nil, ErrOrganizationIDRequired.Wrap(nil)
	}
	if communityID == "" {
		return nil, ErrCommunityIDRequired.Wrap(nil)
	}

	members, err := listCommunityMembers(ctx, s.db, organizationID, communityID)
	if err != nil {
		return nil, ErrListMembersFailed.Wrap(fmt.Errorf("list community members: %w", err))
	}

	return members, nil
}

func (s *Service) UpdateMember(ctx context.Context, organizationID string, communityID string, userID string, input UpdateCommunityMemberInput) (CommunityMember, error) {
	organizationID = strings.TrimSpace(organizationID)
	communityID = strings.TrimSpace(communityID)
	userID = strings.TrimSpace(userID)
	input = sanitizeMemberInput(input)

	switch {
	case organizationID == "":
		return CommunityMember{}, ErrOrganizationIDRequired.Wrap(nil)
	case communityID == "":
		return CommunityMember{}, ErrCommunityIDRequired.Wrap(nil)
	case userID == "":
		return CommunityMember{}, ErrUserIDRequired.Wrap(nil)
	case input.Role == nil && input.Status == nil:
		return CommunityMember{}, ErrMemberRoleRequired.Wrap(nil)
	}

	if input.Role != nil && !isValidMembershipRole(*input.Role) {
		return CommunityMember{}, ErrInvalidMemberRole.Wrap(nil)
	}
	if input.Status != nil && !isValidMembershipStatus(*input.Status) {
		return CommunityMember{}, ErrInvalidMemberStatus.Wrap(nil)
	}

	member, err := updateCommunityMember(ctx, s.db, organizationID, communityID, userID, input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return CommunityMember{}, ErrMemberNotFound.Wrap(nil)
		}
		return CommunityMember{}, ErrUpdateMemberFailed.Wrap(fmt.Errorf("update community member: %w", err))
	}

	return member, nil
}

func (s *Service) CreateInvitation(ctx context.Context, organizationID string, communityID string, invitedBy string, input CreateCommunityInvitationInput) (CommunityInvitation, error) {
	organizationID = strings.TrimSpace(organizationID)
	communityID = strings.TrimSpace(communityID)
	invitedBy = strings.TrimSpace(invitedBy)
	input = sanitizeInvitationInput(input)

	switch {
	case organizationID == "":
		return CommunityInvitation{}, ErrOrganizationIDRequired.Wrap(nil)
	case communityID == "":
		return CommunityInvitation{}, ErrCommunityIDRequired.Wrap(nil)
	case input.Email == "":
		return CommunityInvitation{}, ErrInvitationEmailRequired.Wrap(nil)
	case input.Role == "":
		return CommunityInvitation{}, ErrInvitationRoleRequired.Wrap(nil)
	case !isValidMembershipRole(input.Role):
		return CommunityInvitation{}, ErrInvalidMemberRole.Wrap(nil)
	}

	if input.ExpiresAt == nil {
		expiresAt := time.Now().UTC().Add(7 * 24 * time.Hour)
		input.ExpiresAt = &expiresAt
	}

	invitation, err := insertCommunityInvitation(ctx, s.db, organizationID, communityID, invitedBy, input)
	if err != nil {
		return CommunityInvitation{}, ErrCreateInvitationFailed.Wrap(fmt.Errorf("create community invitation: %w", err))
	}

	return invitation, nil
}

func (s *Service) ListInvitations(ctx context.Context, organizationID string, communityID string) ([]CommunityInvitation, error) {
	organizationID = strings.TrimSpace(organizationID)
	communityID = strings.TrimSpace(communityID)
	if organizationID == "" {
		return nil, ErrOrganizationIDRequired.Wrap(nil)
	}
	if communityID == "" {
		return nil, ErrCommunityIDRequired.Wrap(nil)
	}

	invitations, err := listCommunityInvitations(ctx, s.db, organizationID, communityID)
	if err != nil {
		return nil, ErrListInvitationsFailed.Wrap(fmt.Errorf("list community invitations: %w", err))
	}

	return invitations, nil
}

func (s *Service) AcceptInvitation(ctx context.Context, params acceptCommunityInvitationParams) (CommunityMember, error) {
	params.UserID = strings.TrimSpace(params.UserID)
	params.Email = strings.ToLower(strings.TrimSpace(params.Email))
	params.Token = strings.TrimSpace(params.Token)

	switch {
	case params.UserID == "":
		return CommunityMember{}, ErrUserIDRequired.Wrap(nil)
	case params.Email == "":
		return CommunityMember{}, ErrInvitationEmailRequired.Wrap(nil)
	case params.Token == "":
		return CommunityMember{}, ErrInvitationNotFound.Wrap(nil)
	}

	var member CommunityMember
	err := db.RunInTx(ctx, s.db, func(tx pgx.Tx) error {
		invitation, err := findCommunityInvitationByToken(ctx, tx, params.Token)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrInvitationNotFound.Wrap(nil)
			}
			return ErrAcceptInvitationFailed.Wrap(fmt.Errorf("find community invitation: %w", err))
		}

		switch {
		case invitation.AcceptedAt != nil:
			return ErrInvitationAlreadyAccepted.Wrap(nil)
		case time.Now().UTC().After(invitation.ExpiresAt):
			return ErrInvitationExpired.Wrap(nil)
		case strings.ToLower(strings.TrimSpace(invitation.Email)) != params.Email:
			return ErrInvitationEmailMismatch.Wrap(nil)
		}

		if err := acceptCommunityInvitation(ctx, tx, invitation, params.UserID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrInvitationAlreadyAccepted.Wrap(nil)
			}
			return ErrAcceptInvitationFailed.Wrap(fmt.Errorf("accept community invitation: %w", err))
		}

		member, err = findCommunityMemberByID(ctx, tx, invitation.OrganizationID, invitation.CommunityID, params.UserID)
		if err != nil {
			return ErrAcceptInvitationFailed.Wrap(fmt.Errorf("load community member: %w", err))
		}

		return nil
	})
	if err != nil {
		return CommunityMember{}, err
	}

	return member, nil
}

func sanitizeCreateInput(input CreateCommunityInput) CreateCommunityInput {
	input.OrganizationID = strings.TrimSpace(input.OrganizationID)
	input.Name = strings.TrimSpace(input.Name)
	input.Slug = strings.ToLower(strings.TrimSpace(input.Slug))
	input.Address = strings.TrimSpace(input.Address)
	input.Timezone = strings.TrimSpace(input.Timezone)
	input.Status = strings.ToUpper(strings.TrimSpace(input.Status))

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

func sanitizeUpdateInput(input UpdateCommunityInput) UpdateCommunityInput {
	if input.Name != nil {
		value := strings.TrimSpace(*input.Name)
		input.Name = &value
	}
	if input.Slug != nil {
		value := strings.ToLower(strings.TrimSpace(*input.Slug))
		input.Slug = &value
	}
	if input.Address != nil {
		value := strings.TrimSpace(*input.Address)
		input.Address = &value
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

func sanitizeMemberInput(input UpdateCommunityMemberInput) UpdateCommunityMemberInput {
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

func sanitizeInvitationInput(input CreateCommunityInvitationInput) CreateCommunityInvitationInput {
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
