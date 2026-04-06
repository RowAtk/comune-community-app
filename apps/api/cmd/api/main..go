package main

import (
	"comune/apps/api/internal/modules/auth"
	"comune/apps/api/internal/platform/config"
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger, err := newLogger(cfg)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("failed to create postgres pool", zap.Error(err))
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		logger.Fatal("failed to connect to postgres", zap.Error(err))
	} else {
		logger.Info("Successfully connected to database")
	}

	authService := auth.NewService(dbPool, cfg.SessionSecret, cfg.SessionDuration)
	authHandler := auth.NewHandler(authService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	authHandler.Register(mux)

	server := &http.Server{
		Addr:              cfg.APIAddr,
		Handler:           logRequests(logger, mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Info(
		"api listening",
		zap.String("addr", cfg.APIAddr),
		zap.String("app_env", cfg.AppEnv),
		zap.String("database_host", dbPool.Config().ConnConfig.Host),
		zap.Duration("session_duration", cfg.SessionDuration),
	)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("server stopped unexpectedly", zap.Error(err))
	}
}

func newLogger(cfg config.Config) (*zap.Logger, error) {
	if cfg.IsDevelopment() {
		return zap.NewDevelopment()
	}

	return zap.NewProduction()
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
