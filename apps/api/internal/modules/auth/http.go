package auth

import (
	"comune/apps/api/internal/modules/users"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

const sessionCookieName = "comune_session"

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/auth/signup", h.handleSignup)
	mux.HandleFunc("POST /v1/auth/login", h.handleLogin)
	mux.HandleFunc("GET /v1/auth/me", h.handleMe)
	mux.HandleFunc("POST /v1/auth/logout", h.handleLogout)
}

func (h *Handler) handleSignup(w http.ResponseWriter, r *http.Request) {
	var input SignupInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	result, err := h.service.Signup(input)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, ErrEmailTaken) {
			status = http.StatusConflict
		}

		writeError(w, status, err.Error())
		return
	}

	writeSessionCookie(w, result.Session)
	writeJSON(w, http.StatusCreated, authResponse(result))
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var input LoginInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	result, err := h.service.Login(input)
	if err != nil {
		status := http.StatusUnauthorized
		if errors.Is(err, ErrMissingCredentials) {
			status = http.StatusBadRequest
		}
		if errors.Is(err, ErrInactiveUser) {
			status = http.StatusForbidden
		}

		writeError(w, status, err.Error())
		return
	}

	writeSessionCookie(w, result.Session)
	writeJSON(w, http.StatusOK, authResponse(result))
}

func (h *Handler) handleMe(w http.ResponseWriter, r *http.Request) {
	token, err := readSessionToken(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	result, err := h.service.GetSession(token)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	writeJSON(w, http.StatusOK, authResponse(result))
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	token, err := readSessionToken(r)
	if err == nil {
		h.service.Logout(token)
	}

	clearSessionCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

type responseEnvelope struct {
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

type authPayload struct {
	User        users.User        `json:"user"`
	AuthAccount AuthAccount `json:"auth_account"`
	Session     SessionView `json:"session"`
}

type SessionView struct {
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func authResponse(result AuthResult) responseEnvelope {
	return responseEnvelope{
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

func decodeJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}

func writeJSON(w http.ResponseWriter, status int, payload responseEnvelope) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, responseEnvelope{Error: message})
}
