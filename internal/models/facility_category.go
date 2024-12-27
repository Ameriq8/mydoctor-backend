package models

type FacilityCategory struct {
	BaseModel
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    *int64 `json:"parent_id,omitempty"`
}
