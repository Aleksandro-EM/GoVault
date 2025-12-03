package models

import (
	"gorm.io/gorm"
)

type Password struct {
	gorm.Model
	Website  string `gorm:"not null"`
	Username string `gorm:"not null"`
	Password string `gorm:"not null"`
	Notes    string `gorm:"not null"`
	UserID   string `gorm:"not null"`
}

type User struct {
	gorm.Model
	Email          string     `gorm:"not null;UniqueIndex"`
	MasterPassword string     `gorm:"not null"`
	Password       []Password `gorm:"foreignKey:UserID"`
}
