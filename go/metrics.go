package main

import (
	"AlternativeCryptoExporter/model"
	"fmt"
	"net/http"
	"strconv"
	"time"

	slog "github.com/StevenCyb/SimpleLogging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var ResponseTimeHistogram *prometheus.HistogramVec

// Middleware return a handler that sits before a handler to collect metrics
func MetricsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := model.StatusRecorder{ResponseWriter: w, StatusCode: 200}

		next.ServeHTTP(&rec, r)

		duration := time.Since(start)
		statusCode := strconv.Itoa(rec.StatusCode)
		route := r.URL.Path

		ResponseTimeHistogram.WithLabelValues(route, r.Method, statusCode).Observe(duration.Seconds())
	}
}

func StartMetricsServer() {
	mux := http.NewServeMux()

	slog.Debug(slog.EntryEvent("Setup collector..."))
	ResponseTimeHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_server_request_duration_seconds",
		Help:    "Histogram of response time for handler in seconds.",
		Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
	}, []string{"route", "method", "status_code"})
	prometheus.MustRegister(ResponseTimeHistogram)

	mux.HandleFunc("/metrics", MetricsMiddleware(HandlerToHandlerFuncWrapper(promhttp.Handler())))

	slog.Info(slog.EntryEvent(fmt.Sprintf("Listen on %s", listen)))
	if err := http.ListenAndServe(listen, mux); err != nil {
		slog.Fatal(slog.Entry{"event": "metrics server startup failed", "error": err.Error()})
	}
}

// HandlerToHandlerFuncWrapper wrap a hander to use it as handler func
func HandlerToHandlerFuncWrapper(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}
