package server

import (
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(rec, r)

		duration := time.Since(start)
		requestID := GetRequestID(r.Context())

		fields := []interface{}{
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration_ms", duration.Milliseconds(),
			"request_id", requestID,
		}

		if rec.status >= 500 {
			Log.Error("request completed", zapFields(fields)...)
		} else if rec.status >= 400 {
			Log.Warn("request completed", zapFields(fields)...)
		} else {
			Log.Info("request completed", zapFields(fields)...)
		}
	})
}
