package models

type Doctor struct {
	BaseModel
	Name              string `json:"name"`
	Specialty         string `json:"specialty"`
	PrimaryFacilityID *int64 `json:"primary_facility_id,omitempty"`
	ContactNumber     string `json:"contact_number"`
	Email             string `json:"email"`
}
