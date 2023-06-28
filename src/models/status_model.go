package models

import (
	"gorm.io/gorm"
)

type Status struct {
	gorm.Model
	Name string `json:"name" gorm:"not null"`
	IsPaid bool `json:"isPaid" gorm:"boolean,not null, default true"`
}