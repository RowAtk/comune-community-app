package server

import (
	"comune/apps/api/internal/modules/auth"
	"comune/apps/api/internal/modules/communities"
	"comune/apps/api/internal/modules/organizations"
	"comune/apps/api/internal/platform/httpx"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func NewMux(logger *zap.Logger, authService *auth.Service, organizationService *organizations.Service, communityService *communities.Service) *http.ServeMux {
	mux := http.NewServeMux()
	authMiddleware := auth.NewMiddleware(authService)
	authHandler := auth.NewHandler(authService)
	organizationHandler := organizations.NewHandler(organizationService)
	communityHandler := communities.NewHandler(communityService)

	mux.Handle("GET /healthz", httpx.Adapt(logger, func(w http.ResponseWriter, r *http.Request) error {
		httpx.WriteJSON(w, http.StatusOK, httpx.ResponseEnvelope{
			Data: map[string]string{"status": "ok"},
		})
		return nil
	}))

	authHandler.Register(mux, logger, authMiddleware.RequireAuth)
	organizationHandler.Register(mux, logger, authMiddleware.RequireAuth)
	communityHandler.Register(mux, logger, authMiddleware.RequireAuth)

	return mux
}

func NewHTTPServer(addr string, logger *zap.Logger, authService *auth.Service, organizationService *organizations.Service, communityService *communities.Service) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           logRequests(logger, NewMux(logger, authService, organizationService, communityService)),
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func logRequests(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		next.ServeHTTP(w, r)
		logger.Info(
			"http request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
			zap.Duration("duration", time.Since(started).Round(time.Millisecond)),
		)
	})
}
