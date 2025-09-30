package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AppClaims struct {
	UserID  int64  `json:"user_id"`
	Email   string `json:"email,omitempty"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(userID int64, secretKey string, isAdmin bool) (string, error) {
	claims := AppClaims{
		UserID:  userID,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
