package models

import "time"

type BaseModel struct {
	ID        int64     `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// SQLRepository provides a SQL-based implementation of the Repository interface.
// I'm using "database/sql" & "github.com/lib/pq"
// Each Repository have this methods
// Find
// FindMany
// Create
// CreateMany
// Update
// UpdateMany
// Delete
// DeleteMany
