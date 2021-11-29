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
var CryptoCurrencyPriceGauge *prometheus.GaugeVec
var CryptoCurrencyVolume24hGauge *prometheus.GaugeVec
var CryptoCurrencyMarketCapGauge *prometheus.GaugeVec
var CryptoCurrencyPriceChangePercentageGauge *prometheus.GaugeVec
var CryptoCurrencyCirculatingSupplyGauge *prometheus.GaugeVec
var CryptoCurrencyTotalSupplyGauge *prometheus.GaugeVec
var CryptoCurrencyMaxSupplyGauge *prometheus.GaugeVec

func init() {
	ResponseTimeHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_server_request_duration_seconds",
		Help:    "Histogram of response time for handler in seconds.",
		Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
	}, []string{"route", "method", "status_code"})
	prometheus.MustRegister(ResponseTimeHistogram)

	CryptoCurrencyCirculatingSupplyGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "crypto_currency_circulating_supply",
		Help: "Last market cap of crypto currency.",
	}, []string{"id", "name", "symbol", "quotes"})
	prometheus.MustRegister(CryptoCurrencyCirculatingSupplyGauge)

	CryptoCurrencyTotalSupplyGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "crypto_currency_total_supply",
		Help: "Last market cap of crypto currency.",
	}, []string{"id", "name", "symbol", "quotes"})
	prometheus.MustRegister(CryptoCurrencyTotalSupplyGauge)

	CryptoCurrencyMaxSupplyGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "crypto_currency_max_supply",
		Help: "Last market cap of crypto currency.",
	}, []string{"id", "name", "symbol", "quotes"})
	prometheus.MustRegister(CryptoCurrencyMaxSupplyGauge)

	CryptoCurrencyPriceGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "crypto_currency_price",
		Help: "Last price of crypto currency.",
	}, []string{"id", "name", "symbol", "quotes"})
	prometheus.MustRegister(CryptoCurrencyPriceGauge)

	CryptoCurrencyVolume24hGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "crypto_currency_volume_24h",
		Help: "Last volume of crypto currency in last 24h.",
	}, []string{"id", "name", "symbol", "quotes"})
	prometheus.MustRegister(CryptoCurrencyVolume24hGauge)

	CryptoCurrencyMarketCapGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "crypto_currency_market_cap",
		Help: "Last market cap of crypto currency.",
	}, []string{"id", "name", "symbol", "quotes"})
	prometheus.MustRegister(CryptoCurrencyMarketCapGauge)

	CryptoCurrencyPriceChangePercentageGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "crypto_currency_price_change_percentage",
		Help: "Last price percentage change of crypto currency.",
	}, []string{"id", "name", "symbol", "quotes", "duration"})
	prometheus.MustRegister(CryptoCurrencyPriceChangePercentageGauge)
}

func metricsMiddleware(next http.HandlerFunc) http.HandlerFunc {
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

func startMetricsServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/metrics", metricsMiddleware(HandlerToHandlerFuncWrapper(promhttp.Handler())))

	slog.Info(slog.EntryEvent(fmt.Sprintf("Listen on %s", listen)))
	if err := http.ListenAndServe(listen, mux); err != nil {
		slog.Fatal(slog.Entry{"event": "metrics server startup failed", "error": err.Error()})
	}
}

func HandlerToHandlerFuncWrapper(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}
