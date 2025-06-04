package middleware

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
	"strings"
	"tilimauth/internal/config"
	"tilimauth/internal/utils"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := VerifyTokens(r, "access")
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, err)
			return
		}

		userIDStr := strconv.FormatInt(userID, 10)
		ctx := context.WithValue(r.Context(), utils.UserIDKey, userIDStr)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func VerifyTokens(request *http.Request, tokenType string /*tokenString string*/) (int64, error) {
	var tokenString string

	if tokenType == "access" {
		authHeader := request.Header.Get("Authorization")
		if authHeader == "" {
			return 0, fmt.Errorf("требуется заголовок c access token")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return 0, fmt.Errorf("недопустимый формат заголовка авторизации")
		}

		tokenString = parts[1]
	} else if tokenType == "refresh" {
		refreshHeader := request.Header.Get("Refresh-Token")
		if refreshHeader == "" {
			return 0, fmt.Errorf("требуется заголовок (header) от refresh token")
		}

		tokenString = refreshHeader
	} else {
		return 0, fmt.Errorf("неизвестный тип токена")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод входа: %v", token.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret), nil
	})

	if err != nil {
		return 0, fmt.Errorf("не удалось спарсить токен: %w", err)
	}

	if !token.Valid {
		return 0, fmt.Errorf("недопустимый токен")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("недопустимые требования к токенам")
	}

	userIDStr := claims["user_id"].(string)
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("недопустимый user_id в токене")
	}

	return userID, nil

}
