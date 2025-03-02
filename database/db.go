package database

import (
	"api/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=localhost user=admin password=marsglace12 dbname=foodtracker port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Erreur de connexion à la base de données : ", err)
	} else {
		DB = db
		log.Default().Println("Connexion à la base de données réussie")
		DB.AutoMigrate(&models.Food{}, &models.Meal{}, &models.User{})
	}

}
