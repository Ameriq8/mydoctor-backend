package models

import "time"

type FacilityInsuranceProvider struct {
	FacilityID          int64          `json:"facility_id"`
	InsuranceProviderID int64          `json:"insurance_provider_id"`
	CoverageDetails     map[string]any `json:"coverage_details"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
}
