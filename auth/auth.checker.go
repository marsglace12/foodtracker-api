package auth

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(c *gin.Context) {

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	cookie, err := c.Request.Cookie("auth_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort() // Interrompt le traitement de la requête
		return
	}

	tokenStr := cookie.Value
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	fmt.Println("token", token)
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	// Ajoute l'email dans le contexte Gin pour qu'il soit accessible dans les handlers suivants
	c.Set("email", claims["email"])
	c.Set("id", claims["id"])
	c.Next() // Appelle le handler suivant

}
