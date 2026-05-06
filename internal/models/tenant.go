package models

import (
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	ID     string `gorm:"type:uuid;primaryKey"`
	UserID string `gorm:"type:uuid;not null"`
	Name   string `gorm:"not null"`
	CreatedAt time.Time
}

func (t *Tenant) GetID() uuid.UUID {
	id, _ := uuid.Parse(t.ID)
	return id
}

func (t *Tenant) GetUserID() uuid.UUID {
	id, _ := uuid.Parse(t.UserID)
	return id
}
