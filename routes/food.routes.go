package routes

import (
	"api/controllers"

	"github.com/gin-gonic/gin"
)

// SetupFoodRoutes définit les routes pour les aliments
func SetupFoodRoutes(router *gin.Engine) {
	// Route sans "/" à la fin
	router.POST("/foods", controllers.CreateFood)
	router.GET("/foods", controllers.GetFoods)

	// Route avec "/" à la fin
	router.POST("/foods/", controllers.CreateFood)
	router.GET("/foods/", controllers.GetFoods)
}
