package routes

import (
	"api/auth"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.Engine) {
	r.GET("/auth/:provider", auth.HandleAuthRequest)
	r.GET("/auth/:provider/callback", auth.GetAuthCallbackFunction)
	r.GET("/auth/:provider/logout", auth.Logout)
	r.GET("/auth/me", auth.AuthMiddleware, auth.GetUser)

}
