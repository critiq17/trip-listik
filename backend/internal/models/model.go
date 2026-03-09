package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	TelegramID int64 `gorm:"uniqueIndex;not null"`
	FirstName  string
	UserName   string
	Places     []Place `gorm:"foreignKey:UserID"`
}

type Place struct {
	gorm.Model
	Name   string
	UserID uint `gorm:"not null"`
}
