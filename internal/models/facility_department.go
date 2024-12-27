package models

import "time"

type FacilityDepartment struct {
	ID            int64     `json:"id"`
	FacilityID    int64     `json:"facility_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	FloorNumber   string    `json:"floor_number"`
	HeadDoctorID  *int64    `json:"head_doctor_id,omitempty"`
	ContactNumber string    `json:"contact_number"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
