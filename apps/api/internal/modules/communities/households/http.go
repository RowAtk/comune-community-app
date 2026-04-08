package households

import (
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
	mux.Handle("POST /v1/organizations/{organizationID}/communities/{id}/households", httpx.Adapt(logger, httpx.Chain(h.handleCreate, requireAuth)))
	mux.Handle("GET /v1/organizations/{organizationID}/communities/{id}/households", httpx.Adapt(logger, httpx.Chain(h.handleList, requireAuth)))
	mux.Handle("GET /v1/organizations/{organizationID}/communities/{id}/households/{householdID}", httpx.Adapt(logger, httpx.Chain(h.handleGetByID, requireAuth)))
	mux.Handle("PATCH /v1/organizations/{organizationID}/communities/{id}/households/{householdID}", httpx.Adapt(logger, httpx.Chain(h.handleUpdate, requireAuth)))
	mux.Handle("DELETE /v1/organizations/{organizationID}/communities/{id}/households/{householdID}", httpx.Adapt(logger, httpx.Chain(h.handleDelete, requireAuth)))
}

func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) error {
	var input CreateInput
	if err := httpx.DecodeJSON(r, &input); err != nil {
		return apperror.Validation("invalid_json", "invalid JSON body", err)
	}

	input.OrganizationID = r.PathValue("organizationID")
	input.CommunityID = r.PathValue("id")

	household, err := h.service.Create(r.Context(), input)
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusCreated, httpx.ResponseEnvelope{Data: household})
	return nil
}

func (h *Handler) handleList(w http.ResponseWriter, r *http.Request) error {
	households, err := h.service.List(r.Context(), r.PathValue("organizationID"), r.PathValue("id"))
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.ResponseEnvelope{Data: households})
	return nil
}

func (h *Handler) handleGetByID(w http.ResponseWriter, r *http.Request) error {
	household, err := h.service.GetByID(r.Context(), r.PathValue("organizationID"), r.PathValue("id"), r.PathValue("householdID"))
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.ResponseEnvelope{Data: household})
	return nil
}

func (h *Handler) handleUpdate(w http.ResponseWriter, r *http.Request) error {
	var input UpdateInput
	if err := httpx.DecodeJSON(r, &input); err != nil {
		return apperror.Validation("invalid_json", "invalid JSON body", err)
	}

	household, err := h.service.Update(r.Context(), r.PathValue("organizationID"), r.PathValue("id"), r.PathValue("householdID"), input)
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.ResponseEnvelope{Data: household})
	return nil
}

func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) error {
	if err := h.service.Delete(r.Context(), r.PathValue("organizationID"), r.PathValue("id"), r.PathValue("householdID")); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
