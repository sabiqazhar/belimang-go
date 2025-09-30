package token

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type AppClaims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email,omitempty"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(userID int64, secretKey string) (string, error) {
	claims := AppClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
