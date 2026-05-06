package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           string    `gorm:"type:uuid;primaryKey"`
	Email        string    `gorm:"unique;not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time
}

func (u *User) GetID() uuid.UUID {
	id, _ := uuid.Parse(u.ID)
	return id
}

func (u *User) GetEmail() string {
	return u.Email
}
