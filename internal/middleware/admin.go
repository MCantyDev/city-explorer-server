package middleware

import (
	"net/http"
	"strings"

	"github.com/MCantyDev/city-explorer-server/internal/services"
	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing or invalid token",
			})
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims, err := services.ValidateJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}

		// Admin Specific Check
		status, _ := services.CheckAdminStatus(claims["id"].(float64))
		if !status {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Admin access required",
			})
			return
		}

		// Setting the ID in the current request context -> Will allow me to return proper data using a select statement
		c.Set("userID", claims["id"])
		c.Set("isAdmin", claims["isAdmin"])
		c.Next()
	}
}
