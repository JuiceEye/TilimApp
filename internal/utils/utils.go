package utils

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"tilimauth/internal/dto/request"
	"tilimauth/internal/middleware"
	"unicode"
)

func ParseRequestBody(r *http.Request, body any) error {
	if r.ContentLength == 0 {
		return nil
	}

	err := json.NewDecoder(r.Body).Decode(body)

	return err
}

func ParseBodyAndValidate(r *http.Request, req request.Request) error {
	// log.Printf("[INFO] Parsing request %v", r.RequestURI)
	if err := ParseRequestBody(r, req); err != nil {
		return fmt.Errorf("ошибка при парсинге JSON тела")
	}

	// jsonBody, _ := json.MarshalIndent(req, "", "  ")
	// fmt.Println(string(jsonBody))

	if err := req.ValidateRequest(); err != nil {
		return err
	}

	return nil
}

func WriteJSONResponse(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// jsonPayload, _ := json.MarshalIndent(payload, "", "  ")
	// log.Printf("[INFO] Response: %v \n", string(jsonPayload))
	// log.Println("-----------------------------------------------------------------------------------------------")
	return json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	if status == 500 {
		log.Printf("\n\n [ERROR] %s \n\n", err.Error())
		err = fmt.Errorf("что-то пошло не так")
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

func GetUserID(r *http.Request) (int64, bool) {
	userIDstr, ok := r.Context().Value(middleware.UserIDKey).(string)
	userID, err := strconv.ParseInt(userIDstr, 10, 64)
	if err != nil {
		// handle error
	}
	return userID, ok
}
