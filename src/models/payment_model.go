package models

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	UserID uint `json:"userId" form:"userid" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;default null"`
	Amount float64 `json:"amount" gorm:"nut null"`
	Comment string `json:"comment" gorm:"text"`
}