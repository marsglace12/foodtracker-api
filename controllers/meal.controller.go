package controllers

import (
	"api/database"
	"api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateMeals(c *gin.Context) {
	var meal models.Meal
	if err := c.ShouldBindJSON(&meal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Vérifier si l'utilisateur existe
	var user models.User
	if err := database.DB.First(&user, meal.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	database.DB.Create(&meal)
	c.JSON(http.StatusCreated, meal)
}

func GetMealsByUser(c *gin.Context) {
	userID := c.Param("user_id")
	var meals []models.Meal
	if err := database.DB.Where("user_id = ?", userID).Preload("Foods").Find(&meals).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meals not found"})
		return
	}

	c.JSON(http.StatusOK, meals)
}
