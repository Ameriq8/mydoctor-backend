package models

import "time"

type FacilityPlan struct {
	ID         int64      `json:"id"`
	FacilityID int64      `json:"facility_id"`
	PlanID     int64      `json:"plan_id"`
	StartDate  time.Time  `json:"start_date"`
	EndDate    *time.Time `json:"end_date,omitempty"`
	IsActive   bool       `json:"is_active"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
