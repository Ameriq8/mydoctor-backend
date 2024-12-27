package models

type Review struct {
	BaseModel
	EntityType string  `json:"entity_type"`
	EntityID   int64   `json:"entity_id"`
	UserID     *int64  `json:"user_id,omitempty"`
	Rating     float64 `json:"rating"`
	Comment    string  `json:"comment"`
}
