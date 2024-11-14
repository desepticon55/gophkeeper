package auth

import (
	"github.com/desepticon55/gophkeeper/internal/model"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func createJWTToken(username string, key string, expirationMinutes int) (string, error) {
	expirationTime := time.Now().Add(time.Duration(expirationMinutes) * time.Minute)
	claims := &model.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}
