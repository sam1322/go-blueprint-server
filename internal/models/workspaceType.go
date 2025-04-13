package models

import "time"

type Workspace struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	JoinCode  string    `json:"join_code"`
	UserId    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
