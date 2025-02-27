package main

import (
	"api/database"
	"api/routes"

	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "FoodTracker API is running!"})
	})
	database.InitDB()
	routes.SetupFoodRoutes(r)
	r.RedirectTrailingSlash = false
	r.Run(":8080")
}
