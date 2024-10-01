package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var (
	responseTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "http_response_time_seconds",
		Help: "Histogram of response time for HTTP requests.",
	})
	statusCodes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_status_codes",
			Help: "HTTP status codes.",
		},
		[]string{"code"},
	)
	logger = logrus.New()
)

func init() {
	prometheus.MustRegister(responseTime)
	prometheus.MustRegister(statusCodes)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
}

func monitor(url string, logger *logrus.Logger) {
	const maxRetries = 3
	const retryDelay = 2 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		start := time.Now()
		resp, err := http.Get(url)
		if err != nil {
			logger.WithError(err).Errorf("Attempt %d: Error making request to %s", attempt, url)

			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt) * retryDelay) // Exponential backoff.
				continue
			} else {
				return // Give up after max retries.
			}
		}
		defer resp.Body.Close()

		elapsed := time.Since(start).Seconds()
		responseTime.Observe(elapsed)
		statusCode := fmt.Sprintf("%d", resp.StatusCode)
		statusCodes.WithLabelValues(statusCode).Inc()

		logger.WithFields(logrus.Fields{
			"url":           url,
			"method":        "GET", // Add the HTTP method
			"status_code":   statusCode,
			"response_time": elapsed,
		}).Info("Request completed")

		return // Success, so return.
	}
}

func monitorHandler(w http.ResponseWriter, r *http.Request) {
	url := os.Getenv("MONITOR_URL")
	monitor(url, logger)
	w.WriteHeader(http.StatusOK)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	url := os.Getenv("MONITOR_URL")
	if url == "" {
		logger.Fatal("MONITOR_URL environment variable not set")
	}

	intervalStr := os.Getenv("MONITOR_INTERVAL")
	if intervalStr == "" {
		intervalStr = "10s"
	}

	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		logger.WithError(err).Fatal("Invalid MONITOR_INTERVAL")
	}

	http.HandleFunc("/monitor", monitorHandler)
	http.Handle("/metrics", promhttp.Handler())

	// Graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		cancel()
		logger.Info("Shutting down metrics exporter...")
	}()

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return // Exit the loop when the context is canceled
			default:
				monitor(url, logger)
				time.Sleep(interval)
			}
		}
	}(ctx)

	logger.Info("Starting HTTP server to expose metrics and monitoring endpoint on :8080")
	http.ListenAndServe(":8080", nil)
}
