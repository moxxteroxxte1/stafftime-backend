package models

import (
	"gorm.io/gorm"
	"time"
)

type Contract struct {
	gorm.Model
	UserID uint `json:"userId" form:"userid" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;default null"`
	StartDate time.Time `json:"startDate" gorm:"default:NOW()"`
	EndDate time.Time `json:"endDate" gorm:"default:current_timestamp"`
	Comment string `json:"comment" gorm:"text"`
	Rate  float64 `json:"rate" gorm:"not null,default 0"`
	Hours float64 `json:"hours" gorm:"not null, default 0"`
}