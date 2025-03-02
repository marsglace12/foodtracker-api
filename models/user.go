package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"` // Ne pas exposer le mot de passe dans les réponses JSON
}
