package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
	"strings"
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

// todo in middleware layer later
func VerifyJWT(r *http.Request) (int, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("Authorization header is required")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0, fmt.Errorf("Invalid authorization header format")
	}

	tokenString := parts[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret), nil
	})

	if err != nil {
		return 0, fmt.Errorf("Failed to parse token: %w", err)
	}

	if !token.Valid {
		return 0, fmt.Errorf("Invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userIDStr := claims["user_id"].(string)
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			return 0, fmt.Errorf("Invalid user_id in token")
		}
		return userID, nil
	}

	return 0, fmt.Errorf("Invalid token claims")
}
