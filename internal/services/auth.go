package services

/* Authorisation Services

- Password Hasher
- JWT Token Validator
*/

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(unhashedPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(unhashedPassword), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CompareHashed(hashedString string, unhashedString string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedString), []byte(unhashedString))
	return err == nil
}

func GenerateCookies(c *gin.Context, id uint, path string) error {
	// Generate Cookies
	sessionToken, err := GenerateSessionJWT(id)
	if err != nil {
		return fmt.Errorf("failed to generate session token")
	}

	refreshToken, err := GenerateRefreshJWT(id)
	if err != nil {
		return fmt.Errorf("failed to generate refresh token")
	}

	// Set Cookies
	// Short Lived Session Token lasting 15 mins (15 * 60)
	c.SetCookie("session_token", sessionToken, 15*60, path, "", true, true)
	// Long Lived Refresh Token lasting 7 days (7 * 24 * 60 * 60)
	c.SetCookie("refresh_token", refreshToken, 7*24*60*60, path, "", true, true)
	return nil
}

func GenerateSessionCookie(c *gin.Context, id uint, path string) error {
	sessionToken, err := GenerateSessionJWT(id)
	if err != nil {
		return fmt.Errorf("failed to generate session token")
	}
	c.SetCookie("session_token", sessionToken, 15*60, path, "", true, true)
	return nil
}

func DeleteCookie(c *gin.Context, name string) {
	c.SetCookie(name, "", -1, "/", "", true, true)
}
