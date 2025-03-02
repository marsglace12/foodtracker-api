package models

import "gorm.io/gorm"

type Meal struct {
	gorm.Model
	Name   string `json:"name"`
	UserID uint   `json:"user_id"`
	User   User   `json:"-" gorm:"foreignKey:UserID"`
	Foods  []Food `json:"foods" gorm:"many2many:meel_foods"`
}
