package middleware

import (
	"net/http"
	"strconv"

	"github.com/MCantyDev/city-explorer-server/internal/services"
	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIdStr, exists := c.Get("userId")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "User ID not found in context",
			})
			return
		}

		// Use the User ID to check admin status
		// Convert to Float64...as im using that (interfaces still weird)
		idParsed, err := strconv.Atoi(userIdStr.(string))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID format",
			})
			return
		}

		isAdmin, err := services.CheckAdminStatus(uint(idParsed))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to check admin status",
			})
			return
		}

		if !isAdmin {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Admin status required",
			})
			return
		}

		c.Set("isAdmin", true)
		c.Next()
	}
}
