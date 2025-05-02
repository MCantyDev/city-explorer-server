package services

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/MCantyDev/city-explorer-server/internal/models"
	"github.com/golang-jwt/jwt"
)

var jwtSecretKey []byte
var jwtDurationInDays int

// Initialisation for JWT (ran in main.go during Server setup)
func InitJWT() error {
	// Initialise the JWT secret key from environment variables
	jwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
	if len(jwtSecretKey) == 0 {
		return fmt.Errorf("ensure \"JWT_SECRET_KEY\n is set in .env file")
	}
	// Initialise the JWT duration (in days) from the environment variables
	var err error
	jwtDurationInDays, err = strconv.Atoi(os.Getenv("JWT_DURATION_DAYS"))
	if err != nil || jwtDurationInDays <= 0 || jwtDurationInDays > 31 {
		return fmt.Errorf("ensure \"JWT_DURATION_DAYS\" is set in .env file and is a positive integer (1 - 31)")
	}

	return nil
}

// Generate JWT - Creates a new JWT Token with user info and expiration date (set with Environment Variable)
func GenerateJWT(user models.User) (string, error) {
	// Claims -> Variables OF the JWT Token
	claims := jwt.MapClaims{
		"id":   user.ID,
		"role": 1, // Temp Value
		"exp":  time.Now().Add(time.Hour * (24 * time.Duration(jwtDurationInDays))).Unix(),
	}

	// Test Cases for Time
	// time.Hour * (24 * time.Duration(jwtDurationInDays))
	// -24 * time.Hour

	// Create the Token using the Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the Token with the Secret Key
	signedToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is expected
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Check if the Token is Valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Return the Claims from the parsed token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("could not parse claims")
	}

	return claims, nil
}
