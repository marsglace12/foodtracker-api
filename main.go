package main

import (
	"api/auth"
	"api/database"
	"api/routes"
	"log"
	"time"

	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // Permet l'envoi des cookies
		MaxAge:           24 * time.Hour,
	}))
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erreur loading .env")
	}
	auth.NewAuth()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "FoodTracker API is running!"})
	})
	database.InitDB()

	routes.SetupFoodRoutes(r)
	routes.SetupMealRoutes(r)
	routes.SetupAuthRoutes(r)
	r.RedirectTrailingSlash = false
	r.Run(":8080")
}
