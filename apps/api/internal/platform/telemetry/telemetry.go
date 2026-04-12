package telemetry

import (
	"comune/apps/api/internal/platform/config"
	"context"
	"errors"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Provider struct {
	TracerProvider *sdktrace.TracerProvider
	Tracer         trace.Tracer
}

func Setup(ctx context.Context, cfg config.Config, logger *zap.Logger) (Provider, func(context.Context) error, error) {
	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			attribute.String("deployment.environment", cfg.AppEnv),
		),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithProcess(),
	)
	if err != nil {
		return Provider{}, nil, err
	}

	options := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	}

	endpoint := strings.TrimSpace(cfg.OTelEndpoint)
	if endpoint != "" {
		exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(endpoint))
		if err != nil {
			return Provider{}, nil, err
		}
		options = append(options, sdktrace.WithBatcher(exporter))
		if logger != nil {
			logger.Info("opentelemetry exporter configured", zap.String("endpoint", endpoint))
		}
	} else if logger != nil {
		logger.Info("opentelemetry exporter not configured; spans will not be exported")
	}

	tracerProvider := sdktrace.NewTracerProvider(options...)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	provider := Provider{
		TracerProvider: tracerProvider,
		Tracer:         tracerProvider.Tracer(cfg.ServiceName),
	}

	shutdown := func(ctx context.Context) error {
		var shutdownErr error
		if err := tracerProvider.Shutdown(ctx); err != nil {
			shutdownErr = errors.Join(shutdownErr, err)
		}
		return shutdownErr
	}

	return provider, shutdown, nil
}
