package models

import (
	"time"

	"gorm.io/gorm"
)

type Shift struct {
	gorm.Model
	UserID    uint      `json:"userId" form:"userid" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;default null"`
	StartTime time.Time `json:"startTime" form:"startTime" gorm:"default:current_timestamp"`
	EndTime   time.Time `json:"endTime" form:"endTime" gorm:"default:current_timestamp"`
	Hurs      float64   `json:"hours" gorm:"not null;defualt null"`
	Comment   string    `json:"comment" form:"comment" gorm:"text"`
	Rate      float64   `json:"rate" form:"rate" gorm:"not null"`
	StatusID  uint      `json:"statusId" form:"statusID" gorm:"foreignKey:StatusID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
