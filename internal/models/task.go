package models

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	TenantID  string `gorm:"type:uuid;not null"`
	Title     string `gorm:"not null"`
	Completed bool   `gorm:"default:false"`
	CreatedAt time.Time
}

func (t *Task) GetID() uuid.UUID {
	id, _ := uuid.Parse(t.ID)
	return id
}

func (t *Task) GetTenantID() uuid.UUID {
	id, _ := uuid.Parse(t.TenantID)
	return id
}
