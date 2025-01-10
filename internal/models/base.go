package models

import "time"

type BaseModel struct {
	ID        int64     `json:"id,omitempty" db:"id" sqlx:"primary_key"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}
