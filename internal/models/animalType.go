package models

import "time"

type Animal struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Species     string    `json:"species"`
	Type        string    `json:"type"`
	Habitat     string    `json:"habitat"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	DietType    string    `json:"dietType"`
	Lifespan    string    `json:"lifespan"`
	FunFact     string    `json:"funFact"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
