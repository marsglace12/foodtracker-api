package controllers

import (
	"api/database"
	"api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ajouter un aliment
func CreateFood(c *gin.Context) {
	var food models.Food
	if err := c.ShouldBindJSON(&food); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Create(&food)
	c.JSON(http.StatusCreated, food)
}

// Obtenir tous les aliments
func GetFoods(c *gin.Context) {
	var foods []models.Food
	database.DB.Find(&foods)
	c.JSON(http.StatusOK, foods)
}
