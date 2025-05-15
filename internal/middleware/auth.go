package middleware

import (
	"net/http"
	"strings"

	"github.com/MCantyDev/city-explorer-server/internal/services"
	"github.com/gin-gonic/gin"
)

// JWT Middleware -> Validates the Authorisation Token sent within the Request in the "Authorization" header (Standardised naming convention thus using "z")
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Getting the Value within the Authorization header
		tokenString := c.GetHeader("Authorization")
		// Checking IF token string has a value and leads with "Bearer "
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing or invalid token",
			})
			return
		}

		// Trimming away the "Bearer " prefix to be left with only the token string
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Validating the JWT Token using authorisation service function
		claims, err := services.ValidateJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}

		// Storing the ID in the context -> May need to also set a Role value in the future
		id, ok := claims["id"]
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to extract ID from token claims",
			})
			return
		}

		isAdmin, ok := claims["isAdmin"]
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to extract information from token claims",
			})
			return
		}

		// Setting the ID in the current request context -> Will allow me to return proper data using a select statement
		c.Set("userID", id)
		c.Set("role", isAdmin)
		c.Next()
	}
}
