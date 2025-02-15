package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
	"tilimauth/internal/config"
	"tilimauth/internal/utils"

	"time"
)

func GenerateJWT(w http.ResponseWriter, userID int) (string, error) {
	secretKey := []byte(config.Envs.JWTSecret)
	expiration := time.Second * time.Duration(config.Envs.JWTExpireSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   strconv.Itoa(userID),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return "", err
	}

	return tokenString, err
}
