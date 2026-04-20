package residents

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
	mux.Handle("POST /v1/organizations/{organizationID}/communities/{id}/residents", httpx.Adapt(logger, httpx.Chain(h.handleCreate, requireAuth)))
	mux.Handle("GET /v1/organizations/{organizationID}/communities/{id}/residents", httpx.Adapt(logger, httpx.Chain(h.handleList, requireAuth)))
	mux.Handle("GET /v1/organizations/{organizationID}/communities/{id}/residents/{residentID}", httpx.Adapt(logger, httpx.Chain(h.handleGetByID, requireAuth)))
	mux.Handle("PATCH /v1/organizations/{organizationID}/communities/{id}/residents/{residentID}", httpx.Adapt(logger, httpx.Chain(h.handleUpdate, requireAuth)))
	mux.Handle("DELETE /v1/organizations/{organizationID}/communities/{id}/residents/{residentID}", httpx.Adapt(logger, httpx.Chain(h.handleDelete, requireAuth)))
	mux.Handle("POST /v1/organizations/{organizationID}/communities/{id}/residents/{residentID}/invitations", httpx.Adapt(logger, httpx.Chain(h.handleCreateInvitation, requireAuth)))
	mux.Handle("GET /v1/organizations/{organizationID}/communities/{id}/residents/{residentID}/invitations", httpx.Adapt(logger, httpx.Chain(h.handleListInvitations, requireAuth)))
	mux.Handle("POST /v1/organizations/{organizationID}/communities/{id}/residents/{residentID}/unlink-user", httpx.Adapt(logger, httpx.Chain(h.handleUnlinkUser, requireAuth)))
	mux.Handle("GET /v1/resident-invitations/{token}", httpx.Adapt(logger, h.handleGetInvitationPreview))
	mux.Handle("POST /v1/resident-invitations/{token}/accept", httpx.Adapt(logger, httpx.Chain(h.handleAcceptInvitation, requireAuth)))
}

func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) error {
	var input CreateInput
	if err := httpx.DecodeJSON(r, &input); err != nil {
		return apperror.Validation("invalid_json", "invalid JSON body", err)
	}

	input.OrganizationID = r.PathValue("organizationID")
	input.CommunityID = r.PathValue("id")

	resident, err := h.service.Create(r.Context(), input)
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusCreated, httpx.ResponseEnvelope{Data: resident})
	return nil
}

func (h *Handler) handleList(w http.ResponseWriter, r *http.Request) error {
	residents, err := h.service.List(r.Context(), r.PathValue("organizationID"), r.PathValue("id"))
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.ResponseEnvelope{Data: residents})
	return nil
}

func (h *Handler) handleGetByID(w http.ResponseWriter, r *http.Request) error {
	resident, err := h.service.GetByID(r.Context(), r.PathValue("organizationID"), r.PathValue("id"), r.PathValue("residentID"))
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.ResponseEnvelope{Data: resident})
	return nil
}

func (h *Handler) handleUpdate(w http.ResponseWriter, r *http.Request) error {
	var input UpdateInput
	if err := httpx.DecodeJSON(r, &input); err != nil {
		return apperror.Validation("invalid_json", "invalid JSON body", err)
	}

	resident, err := h.service.Update(r.Context(), r.PathValue("organizationID"), r.PathValue("id"), r.PathValue("residentID"), input)
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.ResponseEnvelope{Data: resident})
	return nil
}

func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.Delete(r.Context(), r.PathValue("organizationID"), r.PathValue("id"), r.PathValue("residentID")); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *Handler) handleCreateInvitation(w http.ResponseWriter, r *http.Request) error {
	var input CreateInvitationInput
	if err := httpx.DecodeJSON(r, &input); err != nil {
		return apperror.Validation("invalid_json", "invalid JSON body", err)
	}

	current, ok := auth.CurrentAuthResult(r)
	if !ok {
		return apperror.Unauthorized("not_authenticated", "not authenticated", nil)
	}

	invitation, err := h.service.CreateInvitation(
		r.Context(),
		r.PathValue("organizationID"),
		r.PathValue("id"),
		r.PathValue("residentID"),
		current.User.ID,
		input,
	)
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusCreated, httpx.ResponseEnvelope{Data: invitation})
	return nil
}

func (h *Handler) handleListInvitations(w http.ResponseWriter, r *http.Request) error {
	current, ok := auth.CurrentAuthResult(r)
	if !ok {
		return apperror.Unauthorized("not_authenticated", "not authenticated", nil)
	}

	invitations, err := h.service.ListInvitations(
		r.Context(),
		r.PathValue("organizationID"),
		r.PathValue("id"),
		r.PathValue("residentID"),
		current.User.ID,
	)
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.ResponseEnvelope{Data: invitations})
	return nil
}

func (h *Handler) handleUnlinkUser(w http.ResponseWriter, r *http.Request) error {
	current, ok := auth.CurrentAuthResult(r)
	if !ok {
		return apperror.Unauthorized("not_authenticated", "not authenticated", nil)
	}

	resident, err := h.service.UnlinkUser(
		r.Context(),
		r.PathValue("organizationID"),
		r.PathValue("id"),
		r.PathValue("residentID"),
		current.User.ID,
	)
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.ResponseEnvelope{Data: resident})
	return nil
}

func (h *Handler) handleGetInvitationPreview(w http.ResponseWriter, r *http.Request) error {
	preview, err := h.service.GetInvitationPreview(r.Context(), r.PathValue("token"))
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.ResponseEnvelope{Data: preview})
	return nil
}

func (h *Handler) handleAcceptInvitation(w http.ResponseWriter, r *http.Request) error {
	current, ok := auth.CurrentAuthResult(r)
	if !ok {
		return apperror.Unauthorized("not_authenticated", "not authenticated", nil)
	}

	resident, err := h.service.AcceptInvitation(r.Context(), acceptInvitationParams{
		UserID: current.User.ID,
		Token:  r.PathValue("token"),
	})
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.ResponseEnvelope{Data: resident})
	return nil
}
