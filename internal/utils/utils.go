package utils

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"tilimauth/internal/dto/request"
	"unicode"
)

type contextKey string

const UserIDKey = contextKey("userID")

func ParseRequestBody(r *http.Request, body any) error {
	if r.ContentLength == 0 {
		return nil
	}

	err := json.NewDecoder(r.Body).Decode(body)

	return err
}

func ParseBodyAndValidate(r *http.Request, req request.Request) error {
	if err := ParseRequestBody(r, req); err != nil {
		return fmt.Errorf("ошибка при парсинге JSON тела")
	}

	if err := req.ValidateRequest(); err != nil {
		return err
	}

	return nil
}

func WriteJSONResponse(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	switch status {
	case http.StatusInternalServerError:
		log.Printf("\n\n [ERROR] %s \n\n", err.Error())
		err = fmt.Errorf("что-то пошло не так")
	case http.StatusUnauthorized:
		log.Printf("\n\n [ERROR] %s \n\n", err.Error())
		err = fmt.Errorf("пользователь не авторизован")
	}

	_ = WriteJSONResponse(w, status, map[string]string{"error": capitalizeFirst(err.Error())})
}

func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("не удалось захешировать пароль: %w", err)
	}
	return string(hashedBytes), nil
}

func ComparePassword(hashedPassword, userPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userPassword))
}

func GetUserID(r *http.Request) (int64, error) {
	userIDStr, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		panic("userID not found in context — middleware missing?")
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
