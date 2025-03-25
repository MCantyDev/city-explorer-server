package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SignUp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully signed up to City Explorer",
	})
}

func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully logged in to City Explorer",
	})
}

func GetProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"user": "MarkC",
		"role": "admin",
	})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully logged out of City Explorer",
	})
}
