package models

import "time"

type Review struct {
	ID         int64     `json:"id"`
	EntityType string    `json:"entity_type"`
	EntityID   int64     `json:"entity_id"`
	UserID     *int64    `json:"user_id,omitempty"`
	Rating     float64   `json:"rating"`
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
