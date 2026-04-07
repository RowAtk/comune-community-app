package auth

import (
	"comune/apps/api/internal/platform/apperror"
	"comune/apps/api/internal/platform/httpx"
	"context"
	"net/http"
)

type contextKey string

const (
	authResultKey contextKey = "auth_result"
)

type Middleware struct {
	service *Service
}

func NewMiddleware(service *Service) *Middleware {
	return &Middleware{service: service}
}

func (m *Middleware) RequireAuth(next httpx.HandlerFunc) httpx.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		token, err := readSessionToken(r)
		if err != nil {
			return apperror.Unauthorized("not_authenticated", "not authenticated", err)
		}

		result, err := m.service.GetSession(token)
		if err != nil {
			return err
		}

		ctx := context.WithValue(r.Context(), authResultKey, result)
		return next(w, r.WithContext(ctx))
	}
}

func CurrentAuthResult(r *http.Request) (AuthResult, bool) {
	result, ok := r.Context().Value(authResultKey).(AuthResult)
	return result, ok
}
