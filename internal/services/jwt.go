package services

import (
	"fmt"
	"time"

	"github.com/MCantyDev/city-explorer-server/internal/config"
	"github.com/golang-jwt/jwt"
)

// Generate JWT - Generates a JWT token with given claims
func generateJWT(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(config.Cfg.JWT.SecretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
func GenerateRefreshJWT(id uint) (string, error) {
	claims := jwt.MapClaims{
		"sub":     fmt.Sprint(id),
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
		"refresh": true,
	}
	return generateJWT(claims)
}
func GenerateSessionJWT(id uint) (string, error) {
	claims := jwt.MapClaims{
		"sub":     fmt.Sprint(id),
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
		"session": true,
	}
	return generateJWT(claims)
}

// Validate JWT - Validates if the JWT token is valid or not
func validateJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return config.Cfg.JWT.SecretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("could not parse claims")
	}

	return claims, nil
}

// Public Functions for Token Validations
func ValidateRefreshToken(tokenStr string) (jwt.MapClaims, error) {
	claims, err := validateJWT(tokenStr)
	if err != nil {
		return nil, err
	}

	if isRefresh, ok := claims["refresh"].(bool); !ok || !isRefresh {
		return nil, fmt.Errorf("not a refresh token")
	}
	return claims, nil
}
func ValidateSessionToken(tokenStr string) (jwt.MapClaims, error) {
	claims, err := validateJWT(tokenStr)
	if err != nil {
		return nil, err
	}

	if isSession, ok := claims["session"].(bool); !ok || !isSession {
		return nil, fmt.Errorf("not a session token")
	}
	return claims, nil
}
