package models

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
	BaseModel
	Name             string       `json:"name" db:"name"`
	Type             FacilityType `json:"type" db:"type"`
	CategoryID       *int64       `json:"category_id" db:"category_id"`
	CityID           *int64       `json:"city_id" db:"city_id"`
	Location         string       `json:"location" db:"location"`
	Coordinates      *string      `json:"coordinates" db:"coordinates"` // Use *string for nullable column
	Phone            string       `json:"phone" db:"phone"`
	EmergencyPhone   *string      `json:"emergency_phone" db:"emergency_phone"`
	Email            string       `json:"email" db:"email"`
	Website          *string      `json:"website" db:"website"`
	Rating           *float32     `json:"rating" db:"rating"`
	BedCapacity      int          `json:"bed_capacity" db:"bed_capacity"`
	Is24Hours        bool         `json:"is_24_hours" db:"is_24_hours"`
	HasEmergency     bool         `json:"has_emergency" db:"has_emergency"`
	HasParking       bool         `json:"has_parking" db:"has_parking"`
	HasAmbulance     bool         `json:"has_ambulance" db:"has_ambulance"`
	AcceptsInsurance bool         `json:"accepts_insurance" db:"accepts_insurance"`
	Description      *string      `json:"description" db:"description"`
	ImageURL         *string      `json:"image_url" db:"image_url"`
	Amenities        *string      `json:"amenities" db:"amenities"`
	Accreditations   *string      `json:"accreditations" db:"accreditations"`
	MetaData         *string      `json:"meta_data" db:"meta_data"`
}
