package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/torbenconto/spooler/config"
	"github.com/torbenconto/spooler/models"
)

type CustomClaims struct {
	Email  string `json:"email"`
	Role   string `json:"role"`
	UserID uint   `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJWT(email string, role models.Role, ID uint) (string, error) {
	expiration := time.Now().Add(72 * time.Hour)

	claims := CustomClaims{
		Email:  email,
		Role:   string(role),
		UserID: ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := config.Cfg.SecretKey

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJWT(tokenString string) (*CustomClaims, error) {
	secret := config.Cfg.SecretKey

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
