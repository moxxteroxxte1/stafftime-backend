package models

import (
	"gorm.io/gorm"
	"time"
)

type Key struct {
	gorm.Model
	Key         string    `json:"key" gorm:"not null;default:null"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description,omitempty" gorm:"text"`
	ExpiresAt   time.Time `json:"expires-at" gorm:"not null;default:current_timestamp"`
}
