package models

type City struct {
	BaseModel
	Name       string `json:"name"`
	Pupulation int    `json:"population"`
	ImageURL   string `json:"image_url"`
	Timezone   string `json:"timezone"`
}
