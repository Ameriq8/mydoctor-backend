package repositories

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	repoQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "repository_query_duration_seconds",
			Help:    "Duration of repository queries in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "table", "status"}, // Add "table" label
	)

	repoQueryTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "repository_query_total",
			Help: "Total number of repository queries executed",
		},
		[]string{"operation", "table", "status"}, // Add "table" label
	)
)

func init() {
	prometheus.MustRegister(repoQueryDuration)
	prometheus.MustRegister(repoQueryTotal)
}

// trackMetrics is a helper function to measure and record metrics.
func trackMetrics(operation string, table string, start time.Time, err error) {
	duration := time.Since(start).Seconds()
	status := "success"
	if err != nil {
		status = "failure"
	}

	repoQueryDuration.WithLabelValues(operation, table, status).Observe(duration)
	repoQueryTotal.WithLabelValues(operation, table, status).Inc()
}
