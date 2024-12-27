package models

import "time"

type City struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Pupulation int       `json:"population"`
	ImageUrl   string    `json:"image_url"`
	Timezone   string    `json:"timezone"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
