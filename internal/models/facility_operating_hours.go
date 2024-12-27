package models

import "time"

type FacilityOperatingHours struct {
	ID           int64      `json:"id"`
	FacilityID   int64      `json:"facility_id"`
	DepartmentID *int64     `json:"department_id,omitempty"`
	DayOfWeek    int        `json:"day_of_week"`
	StartTime    *time.Time `json:"start_time,omitempty"`
	EndTime      *time.Time `json:"end_time,omitempty"`
	IsClosed     bool       `json:"is_closed"`
	CreatedAt    time.Time  `json:"created_at"`
}
