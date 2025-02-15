package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"tilimauth/internal/config"

	"time"
)

func GenerateJWT(secretKey []byte, userID int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpireSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   strconv.Itoa(userID),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, err
}
