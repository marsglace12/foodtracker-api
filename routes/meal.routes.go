package routes

import (
	"api/controllers"

	"github.com/gin-gonic/gin"
)

// SetupFoodRoutes définit les routes pour les aliments
func SetupMealRoutes(router *gin.Engine) {
	// Route sans "/" à la fin
	router.POST("/meal", controllers.CreateMeals)
	router.GET("/meals", controllers.GetMealsByUser)

	// Route avec "/" à la fin
	router.POST("/meal/", controllers.CreateFood)
	router.GET("/meals/", controllers.GetMealsByUser)
}
