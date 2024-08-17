package utils

import (
	"errors"
	"strings"
	"time"

	"github.com/nishanth-thoughtclan/student-api/config"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(userID string) (string, error) {
	cfg := config.LoadConfig()
	var secretKey = []byte(cfg.JWTSecretKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})
	return token.SignedString(secretKey)
}

func ValidateJWTToken(tokenString string) (jwt.MapClaims, error) {
	cfg := config.LoadConfig()
	var secretKey = []byte(cfg.JWTSecretKey)
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	if tokenString == "" {
		return nil, errors.New("empty token")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}
