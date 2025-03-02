package auth

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func NewAuth() {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	clientCallbackURL := os.Getenv("CLIENT_CALLBACK_URL")
	if clientID == "" || clientSecret == "" || clientCallbackURL == "" {
		log.Fatal("Environment variables (CLIENT_ID, CLIENT_SECRET, CLIENT_CALLBACK_URL) are required")
	}
	secret := os.Getenv("SESSION_KEY")
	var store = sessions.NewCookieStore([]byte(secret))
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Mets true si HTTPS
		SameSite: http.SameSiteLaxMode,
	}
	gothic.Store = store
	goth.UseProviders(
		google.New(clientID, clientSecret, clientCallbackURL),
	)
}

func GetAuthCallbackFunction(c *gin.Context) {
	provider := c.Param("provider")
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", provider))
	session, err := gothic.Store.Get(c.Request, "gothic-session")
	if err != nil {
		log.Println("Error retrieving session in AuthCallback:", err)
	} else {
		log.Println("Session retrieved in AuthCallback, session ID:", session.ID)
	}
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		log.Println("❌ Error completing authentication:", err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := generateJWT(user.Email, user.UserID, time.Now().Add(time.Hour*24).Unix())
	if err != nil {
		log.Println("❌ Erreur de génération du token:", err)
		http.Error(c.Writer, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := generateJWT(user.Email, user.UserID, time.Now().Add(time.Hour*24*7).Unix())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Stocker le Refresh Token en cookie sécurisé
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Mettre true en prod avec HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	// Stocker le token en cookie sécurisé
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,  // Empêche l’accès via JavaScript
		Secure:   false, // Passe à true si HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	log.Println("✅ Authenticated user:", user)
	c.Redirect(http.StatusFound, "http://localhost:5173/")
}

func HandleAuthRequest(c *gin.Context) {
	provider := c.Param("provider")
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", provider))

	// Démarrer l'authentification avec le provider
	if _, err := gothic.CompleteUserAuth(c.Writer, c.Request); err == nil {
		return
	} else {
		gothic.BeginAuthHandler(c.Writer, c.Request)
	}
}

func Logout(c *gin.Context) {
	gothic.Logout(c.Writer, c.Request)
	c.Writer.Header().Set("Location", "/")
	c.Writer.WriteHeader(http.StatusTemporaryRedirect)
}

func generateJWT(userEmail string, userID string, exp int64) (string, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	claims := jwt.MapClaims{
		"id":    userID,
		"email": userEmail,
		"exp":   exp, // Expiration dans 24h
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GetUser(c *gin.Context) {
	// Récupérer l'email stocké dans le middleware
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	// Retourner l'email de l'utilisateur
	c.JSON(http.StatusOK, gin.H{"email": email})
}

func RefreshToken(c *gin.Context) {
	cookie, err := c.Request.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No refresh token"})
		return
	}

	refreshToken := cookie.Value
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Générer un nouveau JWT (exp: court)
	accessToken, err := generateJWT(claims["id"].(string), claims["email"].(string), time.Now().Add(time.Hour*24).Unix())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new access token"})
		return
	}

	// Mettre à jour le cookie auth_token
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "auth_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed"})
}
