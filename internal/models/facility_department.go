package models

type FacilityDepartment struct {
	BaseModel
	FacilityID    int64  `json:"facility_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	FloorNumber   string `json:"floor_number"`
	HeadDoctorID  *int64 `json:"head_doctor_id,omitempty"`
	ContactNumber string `json:"contact_number"`
}
