package repositories

import (
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// MetricsRegistry stores Prometheus metrics for each table.
type MetricsRegistry struct {
	QueryDuration   *prometheus.HistogramVec
	QueryTotal      *prometheus.CounterVec
	InFlightQueries *prometheus.GaugeVec
}

var (
	metricsMap = make(map[string]*MetricsRegistry) // Registry for table-specific metrics
	mu         sync.RWMutex                        // Mutex to handle concurrent access
)

// RegisterMetricsForTable registers Prometheus metrics for a specific table.
func RegisterMetricsForTable(table string) *MetricsRegistry {
	mu.Lock()
	defer mu.Unlock()

	// If metrics already exist, return them
	if metrics, exists := metricsMap[table]; exists {
		return metrics
	}

	// Define custom buckets for query durations (e.g., 1ms, 5ms, 10ms, 100ms, 500ms, 1s, 5s, 10s)
	customBuckets := []float64{0.001, 0.005, 0.01, 0.1, 0.5, 1, 5, 10}

	// Create and register new metrics for the table
	queryDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    fmt.Sprintf("database_%s_query_duration_seconds", table),
			Help:    fmt.Sprintf("Duration of queries on the %s table in seconds", table),
			Buckets: customBuckets,
		},
		[]string{"operation", "status"},
	)

	queryTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("database_%s_query_total", table),
			Help: fmt.Sprintf("Total number of queries executed on the %s table", table),
		},
		[]string{"operation", "status"},
	)

	inFlightQueries := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("database_%s_in_flight_queries", table),
			Help: fmt.Sprintf("Number of in-flight queries on the %s table", table),
		},
		[]string{"operation"},
	)

	// Register metrics with Prometheus
	prometheus.MustRegister(queryDuration)
	prometheus.MustRegister(queryTotal)
	prometheus.MustRegister(inFlightQueries)

	// Store in the registry and return
	metrics := &MetricsRegistry{
		QueryDuration:   queryDuration,
		QueryTotal:      queryTotal,
		InFlightQueries: inFlightQueries,
	}
	metricsMap[table] = metrics
	return metrics
}

// trackMetrics records query metrics for a specific table and operation.
func trackMetrics(operation, table string, start time.Time, err error) {
	metrics := RegisterMetricsForTable(table)

	// Increment in-flight queries
	metrics.InFlightQueries.WithLabelValues(operation).Inc()
	defer metrics.InFlightQueries.WithLabelValues(operation).Dec()

	// Calculate duration and determine status
	duration := time.Since(start).Seconds()
	status := "success"
	if err != nil {
		status = "failure"
	}

	// Update Prometheus metrics
	metrics.QueryDuration.WithLabelValues(operation, status).Observe(duration)
	metrics.QueryTotal.WithLabelValues(operation, status).Inc()
}
