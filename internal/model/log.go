package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Log struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Level     string    `gorm:"index" json:"level"`
	Service   string    `gorm:"index" json:"service"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

func (l*Log) BeforeCreate(tx *gorm.DB) (err error){
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return
}

// defined a Log model with UUID primary key, indexed fields for querying, and a BeforeCreate hook to ensure IDs are generated automatically