package middleware

import (
	"context"
	"net/http"
	"strconv"
	"tilimauth/internal/auth"
	"tilimauth/internal/utils"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.VerifyTokens(r, "access")
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, err)
			return
		}

		userIDStr := strconv.FormatInt(userID, 10)
		ctx := context.WithValue(r.Context(), utils.UserIDKey, userIDStr)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
