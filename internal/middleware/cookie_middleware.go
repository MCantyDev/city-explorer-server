package middleware

import (
	"net/http"

	"github.com/MCantyDev/city-explorer-server/internal/services"
	"github.com/gin-gonic/gin"
)

// New Authentication (Helps stop Session Hijacking as using HttpOnly Cookies)
func CookieMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Cookie
		token, err := c.Cookie("session_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			return
		}

		// Validate the token
		claims, err := services.ValidateSessionToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Extract out the Subject claim
		sub, ok := claims["sub"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})
			return
		}

		c.Set("userId", sub)
		c.Next()
	}
}
