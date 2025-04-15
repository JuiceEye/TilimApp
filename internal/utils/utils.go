package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"tilimauth/internal/dto/request"
	"unicode"
)

func ParseRequestBody(r *http.Request, body any) error {
	if r.ContentLength == 0 {
		return fmt.Errorf("missing body")
	}

	err := json.NewDecoder(r.Body).Decode(body)

	jsonBody, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[INFO] Parsing request %v", string(jsonBody))

	return err
}

func ParseBodyAndValidate(r *http.Request, req request.Request) error {
	if err := ParseRequestBody(r, req); err != nil {
		return errors.New("invalid JSON request")
	}

	if err := req.ValidateRequest(); err != nil {
		return err
	}

	return nil
}

func WriteJSONResponse(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	jsonPayload, _ := json.MarshalIndent(payload, "", "  ")
	log.Printf("[INFO] Response: %v \n", string(jsonPayload))
	fmt.Println("-----------------------------------------------------------------------------------------------")
	return json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	if status == 500 {
		log.Printf("\n\n [ERROR] %s \n\n", err.Error())
		err = errors.New("что-то пошло не так")
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
