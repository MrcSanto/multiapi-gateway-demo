package model

import (
	"time"
)

type User struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	Email     string    `gorm:"size:255;uniqueIndex;not null" json:"email"`
	User      string    `gorm:"column:username;size:255;uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"password,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName especifica o nome da tabela
func (User) TableName() string {
	return "users"
}
