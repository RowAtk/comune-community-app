package communities

import (
	"comune/apps/api/internal/modules/auth"
	"comune/apps/api/internal/platform/apperror"
	"comune/apps/api/internal/platform/httpx"
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(mux *http.ServeMux, logger *zap.Logger, requireAuth httpx.Middleware) {
	mux.Handle("POST /v1/organizations/{organizationID}/communities", httpx.Adapt(logger, httpx.Chain(h.handleCreate, requireAuth)))
	mux.Handle("GET /v1/organizations/{organizationID}/communities", httpx.Adapt(logger, httpx.Chain(h.handleList, requireAuth)))
	mux.Handle("GET /v1/organizations/{organizationID}/communities/{id}", httpx.Adapt(logger, httpx.Chain(h.handleGetByID, requireAuth)))
	mux.Handle("PATCH /v1/organizations/{organizationID}/communities/{id}", httpx.Adapt(logger, httpx.Chain(h.handleUpdate, requireAuth)))
	mux.Handle("DELETE /v1/organizations/{organizationID}/communities/{id}", httpx.Adapt(logger, httpx.Chain(h.handleDelete, requireAuth)))
	mux.Handle("GET /v1/organizations/{organizationID}/communities/{id}/members", httpx.Adapt(logger, httpx.Chain(h.handleListMembers, requireAuth)))
	mux.Handle("PATCH /v1/organizations/{organizationID}/communities/{id}/members/{userID}", httpx.Adapt(logger, httpx.Chain(h.handleUpdateMember, requireAuth)))
	mux.Handle("POST /v1/organizations/{organizationID}/communities/{id}/invitations", httpx.Adapt(logger, httpx.Chain(h.handleCreateInvitation, requireAuth)))
	mux.Handle("GET /v1/organizations/{organizationID}/communities/{id}/invitations", httpx.Adapt(logger, httpx.Chain(h.handleListInvitations, requireAuth)))
	mux.Handle("POST /v1/community-invitations/{token}/accept", httpx.Adapt(logger, httpx.Chain(h.handleAcceptInvitation, requireAuth)))
}

func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) error {
	var input CreateCommunityInput
	if err := httpx.DecodeJSON(r, &input); err != nil {
		return apperror.Validation("invalid_json", "invalid JSON body", err)
	}

	input.OrganizationID = r.PathValue("organizationID")
	community, err := h.service.Create(r.Context(), input)
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusCreated, httpx.ResponseEnvelope{Data: community})
	return nil
}

func (h *Handler) handleList(w http.ResponseWriter, r *http.Request) error {
	communities, err := h.service.List(r.Context(), r.PathValue("organizationID"))
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.ResponseEnvelope{Data: communities})
	return nil
}

func (h *Handler) handleGetByID(w http.ResponseWriter, r *http.Request) error {
	community, err := h.service.GetByID(r.Context(), r.PathValue("organizationID"), r.PathValue("id"))
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.ResponseEnvelope{Data: community})
	return nil
}

func (h *Handler) handleUpdate(w http.ResponseWriter, r *http.Request) error {
	var input UpdateCommunityInput
	if err := httpx.DecodeJSON(r, &input); err != nil {
		return apperror.Validation("invalid_json", "invalid JSON body", err)
	}

	community, err := h.service.Update(r.Context(), r.PathValue("organizationID"), r.PathValue("id"), input)
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.ResponseEnvelope{Data: community})
	return nil
}

func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.Delete(r.Context(), r.PathValue("organizationID"), r.PathValue("id")); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *Handler) handleListMembers(w http.ResponseWriter, r *http.Request) error {
	members, err := h.service.ListMembers(r.Context(), r.PathValue("organizationID"), r.PathValue("id"))
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.ResponseEnvelope{Data: members})
	return nil
}

func (h *Handler) handleUpdateMember(w http.ResponseWriter, r *http.Request) error {
	var input UpdateCommunityMemberInput
	if err := httpx.DecodeJSON(r, &input); err != nil {
		return apperror.Validation("invalid_json", "invalid JSON body", err)
	}

	member, err := h.service.UpdateMember(r.Context(), r.PathValue("organizationID"), r.PathValue("id"), r.PathValue("userID"), input)
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.ResponseEnvelope{Data: member})
	return nil
}

func (h *Handler) handleCreateInvitation(w http.ResponseWriter, r *http.Request) error {
	var input CreateCommunityInvitationInput
	if err := httpx.DecodeJSON(r, &input); err != nil {
		return apperror.Validation("invalid_json", "invalid JSON body", err)
	}

	current, ok := auth.CurrentAuthResult(r)
	if !ok {
		return apperror.Unauthorized("not_authenticated", "not authenticated", nil)
	}

	invitation, err := h.service.CreateInvitation(r.Context(), r.PathValue("organizationID"), r.PathValue("id"), current.User.ID, input)
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusCreated, httpx.ResponseEnvelope{Data: invitation})
	return nil
}

func (h *Handler) handleListInvitations(w http.ResponseWriter, r *http.Request) error {
	invitations, err := h.service.ListInvitations(r.Context(), r.PathValue("organizationID"), r.PathValue("id"))
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.ResponseEnvelope{Data: invitations})
	return nil
}

func (h *Handler) handleAcceptInvitation(w http.ResponseWriter, r *http.Request) error {
	current, ok := auth.CurrentAuthResult(r)
	if !ok {
		return apperror.Unauthorized("not_authenticated", "not authenticated", nil)
	}

	member, err := h.service.AcceptInvitation(r.Context(), acceptCommunityInvitationParams{
		UserID: current.User.ID,
		Email:  current.User.Email,
		Token:  r.PathValue("token"),
	})
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.ResponseEnvelope{Data: member})
	return nil
}
