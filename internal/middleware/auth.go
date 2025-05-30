package middleware

import (
	"context"
	"net/http"
	"tilimauth/internal/auth"
	"tilimauth/internal/utils"
)

type contextKey string

const UserIDKey = contextKey("userID")

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.VerifyTokens(r, "access")
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, err)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
