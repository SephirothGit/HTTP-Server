package server

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (

	// Total requests
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path"},
	)

	// Errors
	HttpErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of HTTP error responses",
		},
		[]string{"method", "path", "status"},
	)

	// Latency
	HttpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func init() {

	prometheus.MustRegister(HttpRequestsTotal)
	prometheus.MustRegister(HttpErrorsTotal)
	prometheus.MustRegister(HttpDuration)

}

func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

func MetricsMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		rec := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(rec, r)

		duration := time.Since(start).Seconds()

		path := r.URL.Path
		method := r.Method
		status := rec.status

		HttpRequestsTotal.WithLabelValues(method, path).Inc()

		HttpDuration.WithLabelValues(method, path).
			Observe(duration)

		if status >= 400 {

			HttpErrorsTotal.WithLabelValues(
				method,
				path,
				http.StatusText(status),
			).Inc()
		}
	})
}
