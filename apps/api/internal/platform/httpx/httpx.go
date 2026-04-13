package httpx

import (
	"comune/apps/api/internal/platform/apperror"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type ResponseEnvelope struct {
	Data  any            `json:"data,omitempty"`
	Error *ErrorResponse `json:"error,omitempty"`
}

type ErrorResponse struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

type HandlerFunc func(http.ResponseWriter, *http.Request) error

type Middleware func(HandlerFunc) HandlerFunc

func Chain(next HandlerFunc, middlewares ...Middleware) HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		next = middlewares[i](next)
	}

	return next
}

func DecodeJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}

func Adapt(logger *zap.Logger, handler HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			WriteAppError(logger, w, r, err)
		}
	})
}

func WriteJSON(w http.ResponseWriter, status int, payload ResponseEnvelope) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, status int, code string, message string) {
	WriteJSON(w, status, ResponseEnvelope{
		Error: &ErrorResponse{
			Code:    code,
			Message: message,
		},
	})
}

func WriteAppError(logger *zap.Logger, w http.ResponseWriter, r *http.Request, err error) {
	appErr, ok := apperror.As(err)
	if !ok {
		appErr = apperror.Internal("internal_error", "internal server error", err)
	}

	status := http.StatusInternalServerError
	switch appErr.Kind {
	case apperror.KindValidation:
		status = http.StatusBadRequest
	case apperror.KindUnauthorized:
		status = http.StatusUnauthorized
	case apperror.KindNotFound:
		status = http.StatusNotFound
	case apperror.KindConflict:
		status = http.StatusConflict
	}

	if logger != nil {
		fields := []zap.Field{
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("error_kind", string(appErr.Kind)),
			zap.String("error_code", appErr.Code),
			zap.String("error_message", appErr.Message),
			zap.String("error_chain", formatErrorChain(err)),
		}
		if appErr.Err != nil {
			fields = append(fields, zap.NamedError("cause", appErr.Err))
		}

		logger.Error(
			"http request failed",
			fields...,
		)
	}

	WriteError(w, status, appErr.Code, appErr.Message)
}

func formatErrorChain(err error) string {
	if err == nil {
		return ""
	}

	var chain []string
	for current := err; current != nil; current = errors.Unwrap(current) {
		chain = append(chain, current.Error())
	}

	return strings.Join(chain, " | caused by: ")
}
