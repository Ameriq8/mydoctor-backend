package models

type City struct {
	BaseModel
	Name       string `json:"name" db:"name"`
	Pupulation int    `json:"population" db:"population"`
	ImageURL   string `json:"image_url" db:"image_url"`
	Timezone   string `json:"timezone"`
}
