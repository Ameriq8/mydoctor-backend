package models

import "time"

type FacilityCertification struct {
	ID               int64      `json:"id"`
	FacilityID       int64      `json:"facility_id"`
	Name             string     `json:"name"`
	IssuingAuthority string     `json:"issuing_authority"`
	IssueDate        *time.Time `json:"issue_date,omitempty"`
	ExpiryDate       *time.Time `json:"expiry_date,omitempty"`
	Status           string     `json:"status"`
	DocumentURL      string     `json:"document_url"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}
