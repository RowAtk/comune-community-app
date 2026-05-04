package invoicing

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
	mux.Handle("GET /v1/organizations/{organizationID}/communities/{id}/invoicing/plans", httpx.Adapt(logger, httpx.Chain(h.handleListPlans, requireAuth)))
	mux.Handle("POST /v1/organizations/{organizationID}/communities/{id}/invoicing/plans", httpx.Adapt(logger, httpx.Chain(h.handleCreatePlan, requireAuth)))
	mux.Handle("POST /v1/organizations/{organizationID}/communities/{id}/invoicing/plans/{planID}/generate", httpx.Adapt(logger, httpx.Chain(h.handleGenerateInvoices, requireAuth)))
	mux.Handle("GET /v1/organizations/{organizationID}/communities/{id}/invoicing/invoices", httpx.Adapt(logger, httpx.Chain(h.handleListInvoices, requireAuth)))
	mux.Handle("POST /v1/organizations/{organizationID}/communities/{id}/invoicing/invoices", httpx.Adapt(logger, httpx.Chain(h.handleCreateInvoice, requireAuth)))
	mux.Handle("GET /v1/organizations/{organizationID}/communities/{id}/invoicing/households/overdue", httpx.Adapt(logger, httpx.Chain(h.handleListOverdueHouseholds, requireAuth)))
	mux.Handle("GET /v1/organizations/{organizationID}/communities/{id}/invoicing/me", httpx.Adapt(logger, httpx.Chain(h.handleGetResidentOverview, requireAuth)))
}

func (h *Handler) handleListPlans(w http.ResponseWriter, r *http.Request) error {
	current, ok := auth.CurrentAuthResult(r)
	if !ok {
		return apperror.Unauthorized("not_authenticated", "not authenticated", nil)
	}

	plans, err := h.service.ListPlans(r.Context(), r.PathValue("organizationID"), r.PathValue("id"), current.User.ID)
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.ResponseEnvelope{Data: plans})
	return nil
}

func (h *Handler) handleCreatePlan(w http.ResponseWriter, r *http.Request) error {
	var input CreateInvoicePlanInput
	if err := httpx.DecodeJSON(r, &input); err != nil {
		return apperror.Validation("invalid_json", "invalid JSON body", err)
	}

	current, ok := auth.CurrentAuthResult(r)
	if !ok {
		return apperror.Unauthorized("not_authenticated", "not authenticated", nil)
	}

	input.OrganizationID = r.PathValue("organizationID")
	input.CommunityID = r.PathValue("id")
	plan, err := h.service.CreatePlan(r.Context(), current.User.ID, input)
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusCreated, httpx.ResponseEnvelope{Data: plan})
	return nil
}

func (h *Handler) handleGenerateInvoices(w http.ResponseWriter, r *http.Request) error {
	var input GenerateInvoicesInput
	if err := httpx.DecodeJSON(r, &input); err != nil {
		return apperror.Validation("invalid_json", "invalid JSON body", err)
	}

	current, ok := auth.CurrentAuthResult(r)
	if !ok {
		return apperror.Unauthorized("not_authenticated", "not authenticated", nil)
	}

	result, err := h.service.GenerateInvoices(
		r.Context(),
		r.PathValue("organizationID"),
		r.PathValue("id"),
		r.PathValue("planID"),
		current.User.ID,
		input,
	)
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.ResponseEnvelope{Data: result})
	return nil
}

func (h *Handler) handleListInvoices(w http.ResponseWriter, r *http.Request) error {
	current, ok := auth.CurrentAuthResult(r)
	if !ok {
		return apperror.Unauthorized("not_authenticated", "not authenticated", nil)
	}

	invoices, err := h.service.ListInvoices(r.Context(), r.PathValue("organizationID"), r.PathValue("id"), current.User.ID)
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.ResponseEnvelope{Data: invoices})
	return nil
}

func (h *Handler) handleCreateInvoice(w http.ResponseWriter, r *http.Request) error {
	var input CreateInvoiceInput
	if err := httpx.DecodeJSON(r, &input); err != nil {
		return apperror.Validation("invalid_json", "invalid JSON body", err)
	}

	current, ok := auth.CurrentAuthResult(r)
	if !ok {
		return apperror.Unauthorized("not_authenticated", "not authenticated", nil)
	}

	input.OrganizationID = r.PathValue("organizationID")
	input.CommunityID = r.PathValue("id")
	invoice, err := h.service.CreateManualInvoice(r.Context(), current.User.ID, input)
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusCreated, httpx.ResponseEnvelope{Data: invoice})
	return nil
}

func (h *Handler) handleListOverdueHouseholds(w http.ResponseWriter, r *http.Request) error {
	current, ok := auth.CurrentAuthResult(r)
	if !ok {
		return apperror.Unauthorized("not_authenticated", "not authenticated", nil)
	}

	summaries, err := h.service.ListOverdueHouseholds(r.Context(), r.PathValue("organizationID"), r.PathValue("id"), current.User.ID)
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.ResponseEnvelope{Data: summaries})
	return nil
}

func (h *Handler) handleGetResidentOverview(w http.ResponseWriter, r *http.Request) error {
	current, ok := auth.CurrentAuthResult(r)
	if !ok {
		return apperror.Unauthorized("not_authenticated", "not authenticated", nil)
	}

	overview, err := h.service.GetResidentOverview(r.Context(), r.PathValue("organizationID"), r.PathValue("id"), current.User.ID)
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.ResponseEnvelope{Data: overview})
	return nil
}
