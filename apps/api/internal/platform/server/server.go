package server

import (
	"bytes"
	"comune/apps/api/internal/modules/auth"
	"comune/apps/api/internal/modules/communities"
	"comune/apps/api/internal/modules/invoicing"
	"comune/apps/api/internal/modules/organizations"
	"comune/apps/api/internal/platform/httpx"
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const requestIDHeader = "X-Request-Id"

func NewMux(
	logger *zap.Logger,
	authService *auth.Service,
	organizationService *organizations.Service,
	communityService *communities.Service,
	invoicingService *invoicing.Service,
) *http.ServeMux {
	mux := http.NewServeMux()
	authMiddleware := auth.NewMiddleware(authService)
	authHandler := auth.NewHandler(authService)
	organizationHandler := organizations.NewHandler(organizationService)
	communityHandler := communities.NewHandler(communityService)
	invoicingHandler := invoicing.NewHandler(invoicingService)

	mux.Handle("GET /healthz", httpx.Adapt(logger, func(w http.ResponseWriter, r *http.Request) error {
		httpx.WriteJSON(w, http.StatusOK, httpx.ResponseEnvelope{
			Data: map[string]string{"status": "ok"},
		})
		return nil
	}))

	authHandler.Register(mux, logger, authMiddleware.RequireAuth)
	organizationHandler.Register(mux, logger, authMiddleware.RequireAuth)
	communityHandler.Register(mux, logger, authMiddleware.RequireAuth)
	invoicingHandler.Register(mux, logger, authMiddleware.RequireAuth)

	return mux
}

func NewHTTPServer(
	addr string,
	logger *zap.Logger,
	authService *auth.Service,
	organizationService *organizations.Service,
	communityService *communities.Service,
	invoicingService *invoicing.Service,
) *http.Server {
	var handler http.Handler = NewMux(logger, authService, organizationService, communityService, invoicingService)
	handler = withTracing(handler)
	handler = logRequests(logger, handler)

	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func logRequests(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := requestIDFromRequest(r)
		w.Header().Set(requestIDHeader, requestID)
		r = r.WithContext(context.WithValue(r.Context(), requestIDContextKey{}, requestID))

		logger.Info(
			"http request started",
			zap.String("request_id", requestID),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
		)

		started := time.Now()
		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(recorder, r)

		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
			zap.Int("status", recorder.statusCode),
			zap.Duration("duration", time.Since(started).Round(time.Millisecond)),
		}
		if body := recorder.body.String(); body != "" {
			fields = append(fields, zap.String("response_body", body))
		}
		if recorder.truncated {
			fields = append(fields, zap.Bool("response_body_truncated", true))
		}

		logger.Info(
			"http request completed",
			fields...,
		)
	})
}

const maxLoggedResponseBodyBytes = 8192

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
	truncated  bool
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(data []byte) (int, error) {
	if r.body.Len() < maxLoggedResponseBodyBytes {
		remaining := maxLoggedResponseBodyBytes - r.body.Len()
		if len(data) > remaining {
			_, _ = r.body.Write(data[:remaining])
			r.truncated = true
		} else {
			_, _ = r.body.Write(data)
		}
	} else {
		r.truncated = true
	}

	return r.ResponseWriter.Write(data)
}

func withTracing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := requestIDFromRequest(r)
		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagationHeaderCarrier(r.Header))
		tracer := otel.Tracer("comune/apps/api/http")
		spanName := r.Method + " " + r.URL.Path
		ctx, span := tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		span.SetAttributes(
			semconv.HTTPRequestMethodKey.String(r.Method),
			attribute.String("url.path", r.URL.Path),
			attribute.String("http.route", r.URL.Path),
			attribute.String("request.id", requestID),
			attribute.String("net.peer.address", r.RemoteAddr),
		)

		recorder := &traceResponseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(recorder, r.WithContext(context.WithValue(ctx, requestIDContextKey{}, requestID)))

		span.SetAttributes(semconv.HTTPResponseStatusCodeKey.Int(recorder.statusCode))
		if recorder.statusCode >= http.StatusInternalServerError {
			span.SetStatus(codes.Error, http.StatusText(recorder.statusCode))
		}
	})
}

type traceResponseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *traceResponseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func requestIDFromRequest(r *http.Request) string {
	if existing := r.Header.Get(requestIDHeader); existing != "" {
		return existing
	}

	if existing, ok := requestIDFromContext(r.Context()); ok {
		return existing
	}

	var buf [16]byte
	if _, err := rand.Read(buf[:]); err == nil {
		return hex.EncodeToString(buf[:])
	}

	return time.Now().UTC().Format("20060102150405.000000000")
}

type requestIDContextKey struct{}

func requestIDFromContext(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(requestIDContextKey{}).(string)
	return requestID, ok
}

type propagationHeaderCarrier http.Header

func (c propagationHeaderCarrier) Get(key string) string {
	return http.Header(c).Get(key)
}

func (c propagationHeaderCarrier) Set(key string, value string) {
	http.Header(c).Set(key, value)
}

func (c propagationHeaderCarrier) Keys() []string {
	headers := http.Header(c)
	keys := make([]string, 0, len(headers))
	for key := range headers {
		keys = append(keys, key)
	}
	return keys
}
