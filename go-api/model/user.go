package model

import "time"

type User struct {
	ID        int        `json:"id_user"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	User      string     `json:"username"`
	Password  string     `json:"password,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
