package models

import "time"

type FacilityType string

const (
	PublicHospital       FacilityType = "Public Hospital"
	TeachingHospital     FacilityType = "Teaching Hospital"
	PrivateHospital      FacilityType = "Private Hospital"
	RehabilitationCenter FacilityType = "Rehabilitation Center"
	MedicalComplex       FacilityType = "Medical Complex"
	Clinic               FacilityType = "Clinic"
	Pharmacy             FacilityType = "Pharmacy"
	Laboratory           FacilityType = "Laboratory"
	ImagingCenter        FacilityType = "Imaging Center"
)

type Facility struct {
	ID               int64        `json:"id"`
	Name             string       `json:"name"`
	Type             FacilityType `json:"type"`
	CategoryID       *int64       `json:"category_id,omitempty"`
	CityID           *int64       `json:"city_id,omitempty"`
	Location         string       `json:"location"`
	Coordinates      string       `json:"coordinates"`
	Phone            string       `json:"phone"`
	EmergencyPhone   string       `json:"emergency_phone"`
	Email            string       `json:"email"`
	Website          string       `json:"website"`
	Rating           float32      `json:"rating"`
	BedCapacity      int          `json:"bed_capacity"`
	Is24Hours        bool         `json:"is_24_hours"`
	HasEmergency     bool         `json:"has_emergency"`
	HasParking       bool         `json:"has_parking"`
	HasAmbulance     bool         `json:"has_ambulance"`
	AcceptsInsurance bool         `json:"accepts_insurance"`
	Description      string       `json:"description"`
	ImageURL         string       `json:"image_url"`
	Amenities        string       `json:"amenities"`
	Accreditations   string       `json:"accreditations"`
	MetaData         string       `json:"meta_data"`
	CreatedAt        time.Time    `json:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at"`
}
