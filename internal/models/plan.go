package models

import "time"

type Plan struct {
	ID           int64          `json:"id"`
	Name         string         `json:"name"`
	MonthlyPrice float64        `json:"monthly_price"`
	YearlyPrice  float64        `json:"yearly_price"`
	Description  string         `json:"description"`
	Features     map[string]any `json:"features"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}
