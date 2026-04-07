package auth

import (
	"comune/apps/api/internal/modules/users"
	"comune/apps/api/internal/platform/apperror"
	"comune/apps/api/internal/platform/httpx"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const sessionCookieName = "comune_session"

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(mux *http.ServeMux, logger *zap.Logger, requireAuth httpx.Middleware) {
	mux.Handle("POST /v1/auth/signup", httpx.Adapt(logger, h.handleSignup))
	mux.Handle("POST /v1/auth/login", httpx.Adapt(logger, h.handleLogin))
	mux.Handle("GET /v1/auth/me", httpx.Adapt(logger, httpx.Chain(h.handleMe, requireAuth)))
	mux.Handle("POST /v1/auth/logout", httpx.Adapt(logger, httpx.Chain(h.handleLogout, requireAuth)))
}

func (h *Handler) handleSignup(w http.ResponseWriter, r *http.Request) error {
	var input SignupInput
	if err := httpx.DecodeJSON(r, &input); err != nil {
		return apperror.Validation("invalid_json", "invalid JSON body", err)
	}

	result, err := h.service.Signup(input)
	if err != nil {
		return err
	}

	writeSessionCookie(w, result.Session)
	httpx.WriteJSON(w, http.StatusCreated, authResponse(result))
	return nil
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) error {
	var input LoginInput
	if err := httpx.DecodeJSON(r, &input); err != nil {
		return apperror.Validation("invalid_json", "invalid JSON body", err)
	}

	result, err := h.service.Login(input)
	if err != nil {
		return err
	}

	writeSessionCookie(w, result.Session)
	httpx.WriteJSON(w, http.StatusOK, authResponse(result))
	return nil
}

func (h *Handler) handleMe(w http.ResponseWriter, r *http.Request) error {
	result, ok := CurrentAuthResult(r)
	if !ok {
		return apperror.Unauthorized("not_authenticated", "not authenticated", nil)
	}

	httpx.WriteJSON(w, http.StatusOK, authResponse(result))
	return nil
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) error {
	token, err := readSessionToken(r)
	if err == nil {
		h.service.Logout(token)
	}

	clearSessionCookie(w)
	w.WriteHeader(http.StatusNoContent)
	return nil
}

type authPayload struct {
	User        users.User  `json:"user"`
	AuthAccount AuthAccount `json:"auth_account"`
	Session     SessionView `json:"session"`
}

type SessionView struct {
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func authResponse(result AuthResult) httpx.ResponseEnvelope {
	return httpx.ResponseEnvelope{
		Data: authPayload{
			User:        result.User,
			AuthAccount: result.AuthAccount,
			Session: SessionView{
				CreatedAt: result.Session.CreatedAt,
				ExpiresAt: result.Session.ExpiresAt,
			},
		},
	}
}

func readSessionToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func writeSessionCookie(w http.ResponseWriter, session Session) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    session.Token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		Expires:  session.ExpiresAt,
	})
}

func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
}
