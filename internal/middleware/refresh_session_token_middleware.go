package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/MCantyDev/city-explorer-server/internal/services"
	"github.com/gin-gonic/gin"
)

func RefreshSessionTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Both Cookies as they are BOTH going to be used here
		sessionToken, err := c.Cookie("session_token")
		if err != nil {
			c.Next()
			return
		}
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			c.Next()
			return
		}

		// Validate Session Token
		sessionClaims, err := services.ValidateSessionToken(sessionToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Validate Refresh Token
		refreshClaims, err := services.ValidateRefreshToken(refreshToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Check Expiry of Session Token
		sessionExpiry := int64(sessionClaims["exp"].(float64))
		sessionExpiration := time.Unix(sessionExpiry, 0)

		if time.Until(sessionExpiration) < 2*time.Minute {
			// Check Refresh Token
			refreshExpiry := int64(refreshClaims["exp"].(float64))
			refreshExpiration := time.Unix(refreshExpiry, 0)

			if time.Until(refreshExpiration) > 0 {
				userIdStr, ok := sessionClaims["sub"].(string)
				if !ok {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"error": "Invalid user ID in token",
					})
					return
				}

				// Convert string to uint
				userId, err := strconv.ParseUint(userIdStr, 10, 64)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"error": "Invalid user ID format",
					})
					return
				}

				services.GenerateSessionCookie(c, uint(userId), "/")
			} else {
				services.DeleteCookie(c, "session_token")
				services.DeleteCookie(c, "refresh_token")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Session expired. Please login again.",
				})
				return
			}
		}

		c.Next()
	}
}
