package middleware

import (
	"net/http"
	"strconv"

	"github.com/MCantyDev/city-explorer-server/internal/services"
	"github.com/gin-gonic/gin"
)

// New Authentication (Helps stop Session Hijacking as using HttpOnly Cookies)

func SessionAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("session_token")
		if err == nil && token != "" {
			claims, err := services.ValidateSessionToken(token)
			if err == nil {
				c.Set("userId", claims["sub"])
				c.Next()
				return
			}
		}

		refreshToken, err := c.Cookie("refresh_token")
		if err != nil || refreshToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid refresh token. Please login.",
			})
			return
		}

		refreshClaims, err := services.ValidateRefreshToken(refreshToken)
		if err != nil {
			services.DeleteCookie(c, "session_token")
			services.DeleteCookie(c, "refresh_token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired refresh token. Please login again.",
			})
			return
		}

		userIdStr, ok := refreshClaims["sub"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid user ID in refresh token",
			})
			return
		}

		userId, err := strconv.ParseUint(userIdStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid user ID format",
			})
			return
		}

		services.GenerateSessionCookie(c, uint(userId), "/")
		c.Set("userId", userIdStr)
		c.Next()
	}
}
