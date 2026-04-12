package main

import (
	"comune/apps/api/internal/modules/auth"
	"comune/apps/api/internal/modules/communities"
	"comune/apps/api/internal/modules/organizations"
	"comune/apps/api/internal/platform/config"
	"comune/apps/api/internal/platform/server"
	"comune/apps/api/internal/platform/telemetry"
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

	_, shutdownTelemetry, err := telemetry.Setup(context.Background(), cfg, logger)
	if err != nil {
		logger.Fatal("failed to initialize telemetry", zap.Error(err))
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdownTelemetry(shutdownCtx); err != nil {
			logger.Error("failed to shutdown telemetry", zap.Error(err))
		}
	}()

	authService := auth.NewService(dbPool, cfg.SessionSecret, cfg.SessionDuration)
	organizationService := organizations.NewService(dbPool)
	communityService := communities.NewService(dbPool)
	server := server.NewHTTPServer(cfg.APIAddr, logger, authService, organizationService, communityService)

	logger.Info(
		"api listening",
		zap.String("addr", cfg.APIAddr),
		zap.String("app_env", cfg.AppEnv),
		zap.String("database_host", dbPool.Config().ConnConfig.Host),
		zap.Duration("session_duration", cfg.SessionDuration),
		zap.String("service_name", cfg.ServiceName),
		zap.Bool("otel_exporter_configured", cfg.OTelEndpoint != ""),
	)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("server stopped unexpectedly", zap.Error(err))
	}
}

func newLogger(cfg config.Config) (*zap.Logger, error) {
	//if cfg.IsDevelopment() {
	//return zap.NewDevelopment()
	//}

	return zap.NewProduction()
}
