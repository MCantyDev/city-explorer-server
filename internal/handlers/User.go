package handlers

import (
	"fmt"
	"net/http"

	"github.com/MCantyDev/city-explorer-server/internal/database"
	"github.com/MCantyDev/city-explorer-server/internal/models"
	"github.com/MCantyDev/city-explorer-server/internal/services"

	"github.com/gin-gonic/gin"
)

func SignUp(c *gin.Context) {
	var req models.SignupRequest

	// Bind incoming JSON request to the SignupRequest struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check if Email OR Username already exists
	var existingUser models.User
	query := database.NewQueryBuilder("SELECT").Table("users").Where("username = ? OR email = ?").Build()
	_, err := database.Execute(&existingUser, query, req.Username, req.Email)
	if err == nil {
		if existingUser.Username == req.Username {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Username already exists",
			})
			return
		}
		if existingUser.Email == req.Email {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Email already exists",
			})
			return
		}
	}

	// Hash the Password
	hashedPassword, err := services.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create User and Save to Database
	user := models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
	}
	query = database.NewQueryBuilder("INSERT").Table("users").Columns("first_name", "last_name", "username", "email", "password").Values(5).Build()
	_, err = database.Execute(nil, query, user.FirstName, user.LastName, user.Username, user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// Return a success message -> Need to Change to return a valid JWT Token
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
	})
}

func Login(c *gin.Context) {
	var req models.LoginRequest

	// Bind incoming JSON request to the LoginRequest struct
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Fetch the user from the database based on Username
	var user models.User
	query := database.NewQueryBuilder("SELECT").Table("users").Where("username = ?").Build()
	_, err = database.Execute(&user, query, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to query database",
		})
		return
	}

	// If no user is found
	if user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User does not exists",
		})
		return
	}

	// Compare the Provided password with the stored hashed password
	isPassword := services.CompareHashed(user.Password, req.Password)
	if !isPassword {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Password does not match",
		})
		return
	}

	// Generate a JWT Token
	token, err := services.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	// Validate the token immediately (for a Test as not connect to frontend for now)
	claims, err := services.ValidateJWT(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return a success message -> Need to change to return a Valid JWT Token
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("successfully logged in to City Explorer as %s", req.Username),
		"token":   token,
	})
}

func GetProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"user": "MarkC",
		"role": "admin",
	})
}
