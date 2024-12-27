package models

import "time"

type FacilityCertification struct {
	BaseModel
	FacilityID       int64      `json:"facility_id"`
	Name             string     `json:"name"`
	IssuingAuthority string     `json:"issuing_authority"`
	IssueDate        *time.Time `json:"issue_date,omitempty"`
	ExpiryDate       *time.Time `json:"expiry_date,omitempty"`
	Status           string     `json:"status"`
	DocumentURL      string     `json:"document_url"`
}
