package models

import "gorm.io/gorm"

type Food struct {
	gorm.Model
	Name     string  `json:"name"`
	Calories float64 `json:"calories"`
	Proteins float64 `json:"proteins"`
	Carbs    float64 `json:"carbs"`
	Fats     float64 `json:"fats"`
}
