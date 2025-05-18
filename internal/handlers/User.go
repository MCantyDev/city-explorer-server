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

	// initialise an Error Values
	userCreationError := models.NewUserCreationError()

	// Bind incoming JSON request to the SignupRequest struct
	if err := c.ShouldBindJSON(&req); err != nil {
		userCreationError.Error = "Internal Server Error (Try again later)"
		c.JSON(http.StatusInternalServerError, userCreationError)
		return
	}

	// Check if Email OR Username already exists
	var existingUser models.User
	query := database.NewQueryBuilder("SELECT").Table("users").Where("username = ? OR email = ?").Build()
	_, err := database.Execute(&existingUser, query, req.Username, req.Email)
	if err == nil {
		if existingUser.Username == req.Username {
			userCreationError.Error = "Username Already Taken"
			c.JSON(http.StatusBadRequest, userCreationError)
			return
		}
		if existingUser.Email == req.Email {
			userCreationError.Error = "Email Already Taken"
			c.JSON(http.StatusBadRequest, userCreationError)
			return
		}
	}

	// Validate the Inputs based on the user creation requirements (see models request_errors.go)
	// Cool way of doing it, creating a slice of structs that take in a value and a function(string) (string, bool)
	// Giving it 3 different pairs
	validations := []struct {
		value    string
		validate func(string) (string, bool)
	}{
		{req.FirstName, services.ValidateFirstName},
		{req.LastName, services.ValidateLastName},
		{req.Password, services.ValidatePassword},
	}

	// Then Loop through the validations
	for _, v := range validations {
		errMsg, valid := v.validate(v.value)
		if !valid {
			userCreationError.Error = errMsg
			c.JSON(http.StatusInternalServerError, userCreationError)
			return
		}
	}

	// Check if Passwords are the same
	if req.Password != req.ConfirmPassword {
		userCreationError.Error = "Passwords Don't Match"
		c.JSON(http.StatusInternalServerError, userCreationError)
		return
	}

	// Hash the Password
	hashedPassword, err := services.HashPassword(req.Password)
	if err != nil {
		userCreationError.Error = "Internal Server Error (Try again later)"
		c.JSON(http.StatusInternalServerError, userCreationError)
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
	query = database.NewQueryBuilder("INSERT").Table("users").Columns("first_name", "last_name", "username", "email", "password", "is_admin").Values(6).Build()
	_, err = database.Execute(nil, query, user.FirstName, user.LastName, user.Username, user.Email, user.Password, false)
	if err != nil {
		userCreationError.Error = "Internal Server Error (Try again later)"
		c.JSON(http.StatusInternalServerError, userCreationError)
		return
	}

	query = database.NewQueryBuilder("SELECT").Table("users").Where("username = ? AND email = ?").Build()
	_, err = database.Execute(&user, query, user.Username, user.Email)
	if err != nil {
		userCreationError.Error = "Internal Server Error (Try again later)"
		c.JSON(http.StatusInternalServerError, userCreationError)
		return
	}

	// Generate Cookies
	path := "/sign-up"
	err = services.GenerateCookies(c, user.Id, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return a success message -> Need to Change to return a valid JWT Token
	c.JSON(http.StatusOK, gin.H{
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"username":  user.Username,
		"email":     user.Email,
		"isAdmin":   user.IsAdmin,
	})
}

func Login(c *gin.Context) {
	var req models.LoginRequest

	// initialise an Error Values
	userLoginError := models.UserLoginError{
		Error: "",
	}

	// Bind incoming JSON request to the LoginRequest struct
	err := c.ShouldBindJSON(&req)
	if err != nil {
		userLoginError.Error = "Internal Server Error (Try again later)"
		c.JSON(http.StatusInternalServerError, userLoginError)
		return
	}

	// Fetch the user from the database based on Username

	var user models.User
	query := database.NewQueryBuilder("SELECT").Table("users").Where("username = ?").Build()
	_, err = database.Execute(&user, query, req.Username)
	if err != nil {
		userLoginError.Error = "Internal Server Error (Try again later)"
		c.JSON(http.StatusInternalServerError, userLoginError)
		return
	}

	// If no user is found
	if user.Username == "" {
		userLoginError.Error = fmt.Sprintf("No User found with Username: %s", req.Username)
		c.JSON(http.StatusUnauthorized, userLoginError)
		return
	}

	// Compare the Provided password with the stored hashed password
	isPassword := services.CompareHashed(user.Password, req.Password)
	if !isPassword {
		userLoginError.Error = "Passwords do not match"
		c.JSON(http.StatusUnauthorized, userLoginError)
		return
	}

	// Generate Cookies
	path := "/"
	err = services.GenerateCookies(c, user.Id, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return the User's Information
	c.JSON(http.StatusOK, gin.H{
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"username":  user.Username,
		"email":     user.Email,
		"isAdmin":   user.IsAdmin,
	})
}

func Logout(c *gin.Context) {
	services.DeleteCookie(c, "session_token")
	services.DeleteCookie(c, "refresh_token")

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}

func GetProfile(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "User ID missing in context",
		})
	}

	var user models.User
	query := database.NewQueryBuilder("SELECT").Table("users").Columns("first_name", "last_name", "username", "email", "is_admin").Where("id = ?").Build()

	_, err := database.Execute(&user, query, userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to retreive User data from Server",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"username":  user.Username,
		"email":     user.Email,
		"isAdmin":   user.IsAdmin,
	})
}
