package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
	"strings"
	"tilimauth/internal/config"
	"tilimauth/internal/dto/response"
	"tilimauth/internal/utils"
	"time"
)

func GenerateTokenPair(w http.ResponseWriter, userID int64) (*response.AccessRefreshTokenPair, error) {
	accessToken, err := GenerateAccessToken(w, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return nil, err
	}

	refreshToken, err := GenerateRefreshToken(w, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return nil, err
	}

	return &response.AccessRefreshTokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// аксесс токен
func GenerateAccessToken(resp http.ResponseWriter, userID int64) (string, error) {
	secretKey := []byte(config.Envs.JWTSecret)
	expiration := time.Second * time.Duration(config.Envs.JWTAccessExpireSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    strconv.FormatInt(userID, 10),
		"expiredAt":  time.Now().Add(expiration).Unix(),
		"token_type": "access",
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		utils.WriteError(resp, http.StatusInternalServerError, err)
		return "", err
	}

	return tokenString, err
}

// рефреш токен
func GenerateRefreshToken(resp http.ResponseWriter, userID int64) (string, error) {
	secretKey := []byte(config.Envs.JWTSecret)
	expiration := time.Second * time.Duration(config.Envs.JWTRefreshExpireSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    strconv.FormatInt(userID, 10),
		"expiredAt":  time.Now().Add(expiration).Unix(),
		"token_type": "refresh",
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		utils.WriteError(resp, http.StatusInternalServerError, err)
		return "", err
	}

	return tokenString, err
}

// todo: add to the middleware layer later
func VerifyTokens(request *http.Request, tokenType string /*tokenString string*/) (int64, error) {
	var tokenString string

	if tokenType == "access" {
		authHeader := request.Header.Get("Authorization")
		if authHeader == "" {
			return 0, fmt.Errorf("access token header is required")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return 0, fmt.Errorf("invalid authorization header format")
		}

		tokenString = parts[1]
	} else if tokenType == "refresh" {
		refreshHeader := request.Header.Get("X-Refresh-Token")
		if refreshHeader == "" {
			return 0, fmt.Errorf("refresh token header is required")
		}

		tokenString = refreshHeader
	} else {
		return 0, fmt.Errorf("unknown token type")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret), nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}

	userIDStr := claims["user_id"].(string)
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid user_id in token")
	}

	return userID, nil

}
