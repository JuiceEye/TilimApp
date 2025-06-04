package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
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
		"expiredAt":  time.Now().UTC().Add(expiration).Unix(),
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
		"expiredAt":  time.Now().UTC().Add(expiration).Unix(),
		"token_type": "refresh",
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		utils.WriteError(resp, http.StatusInternalServerError, err)
		return "", err
	}

	return tokenString, err
}
