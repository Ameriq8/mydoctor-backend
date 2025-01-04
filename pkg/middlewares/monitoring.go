package middlewares

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var httpRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "server_http_requests_total",
		Help: "Total HTTP requests processed by the server",
	},
	[]string{"method", "endpoint", "status"},
)

var httpRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "server_http_request_duration_seconds",
		Help:    "Duration of HTTP requests in seconds",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"method", "endpoint", "status"},
)

var httpRequestSize = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "server_http_request_size_bytes",
		Help:    "Size of HTTP requests in bytes",
		Buckets: prometheus.ExponentialBuckets(100, 10, 6), // 100B to ~1MB
	},
	[]string{"method", "endpoint"},
)

var httpResponseSize = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "server_http_response_size_bytes",
		Help:    "Size of HTTP responses in bytes",
		Buckets: prometheus.ExponentialBuckets(100, 10, 6), // 100B to ~1MB
	},
	[]string{"method", "endpoint", "status"},
)

var httpFailedRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "server_http_failed_requests_total",
		Help: "Total HTTP requests that failed",
	},
	[]string{"method", "endpoint", "status"},
)

func init() {
	// Register Prometheus metrics
	prometheus.MustRegister(httpRequests)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(httpRequestSize)
	prometheus.MustRegister(httpResponseSize)
	prometheus.MustRegister(httpFailedRequests)
}

func MonitoringMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Capture request size
		reqSize := computeRequestSize(c.Request)

		// Process the request
		c.Next()

		// Measure time taken and capture metrics
		duration := time.Since(startTime).Seconds()
		status := c.Writer.Status()
		method := c.Request.Method
		endpoint := c.FullPath() // Use FullPath to get the defined route pattern (e.g., "/api/v1/resource")
		respSize := float64(c.Writer.Size())

		// Update Prometheus metrics
		httpRequests.WithLabelValues(method, endpoint, strconv.Itoa(status)).Inc()
		httpRequestDuration.WithLabelValues(method, endpoint, strconv.Itoa(status)).Observe(duration)
		httpRequestSize.WithLabelValues(method, endpoint).Observe(reqSize)
		httpResponseSize.WithLabelValues(method, endpoint, strconv.Itoa(status)).Observe(respSize)

		// Count failed requests (status >= 400)
		if status >= 400 {
			httpFailedRequests.WithLabelValues(method, endpoint, strconv.Itoa(status)).Inc()
		}
	}
}

// Helper function to compute request size
func computeRequestSize(r *http.Request) float64 {
	size := 0
	if r.URL != nil {
		size += len(r.URL.String())
	}

	size += len(r.Method)
	size += len(r.Proto)

	for name, values := range r.Header {
		size += len(name)
		for _, value := range values {
			size += len(value)
		}
	}

	if r.ContentLength > 0 {
		size += int(r.ContentLength)
	}

	return float64(size)
}
