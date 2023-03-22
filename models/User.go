package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name  string `gorm:"type:varchar(20);no null"`
	Admin bool
}
