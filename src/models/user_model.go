package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Firstname string `json:"firstname,omitempty" gorm:"text, not null, defualt null"`
	Lastname  string `json:"lastname,omitempty" gorm:"text, not null, default null"`
	Username  string `json:"username,omitempty" gorm:"unique;not null"`
	Email     string `json:"email,omitempty" gorm:"text, default null"`
	Password  string `json:"password,omitempty" gorm:"text,not null,default null"`
	IsAdmin   bool   `json:"isAdmin,omitempty" gorm:"default false"`
	ImgUrl    string `json:"imageUrl" gorm:"not null,default:/avatar.png"`
}
