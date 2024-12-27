package models

import "time"

type FacilityEquipment struct {
	ID                  int64      `json:"id"`
	FacilityID          int64      `json:"facility_id"`
	DepartmentID        *int64     `json:"department_id,omitempty"`
	Name                string     `json:"name"`
	Model               string     `json:"model"`
	Manufacturer        string     `json:"manufacturer"`
	PurchaseDate        *time.Time `json:"purchase_date,omitempty"`
	LastMaintenanceDate *time.Time `json:"last_maintenance_date,omitempty"`
	NextMaintenanceDate *time.Time `json:"next_maintenance_date,omitempty"`
	Status              string     `json:"status"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}
