package models

import "time"

type Doctor struct {
	ID                int64     `json:"id"`
	Name              string    `json:"name"`
	Specialty         string    `json:"specialty"`
	PrimaryFacilityID *int64    `json:"primary_facility_id,omitempty"`
	ContactNumber     string    `json:"contact_number"`
	Email             string    `json:"email"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
