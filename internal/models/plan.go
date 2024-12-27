package models

type Plan struct {
	BaseModel
	Name         string         `json:"name"`
	MonthlyPrice float64        `json:"monthly_price"`
	YearlyPrice  float64        `json:"yearly_price"`
	Description  string         `json:"description"`
	Features     map[string]any `json:"features"`
}
