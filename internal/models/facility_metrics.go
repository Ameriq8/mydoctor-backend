package models

import "time"

type FacilityMetrics struct {
	FacilityID  int64     `json:"facility_id"`
	MetricType  string    `json:"metric_type"`
	MetricValue float64   `json:"metric_value"`
	MeasuredAt  time.Time `json:"measured_at"`
	CreatedAt   time.Time `json:"created_at"`
}
