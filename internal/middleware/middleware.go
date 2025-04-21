package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func CreateChain(mws ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(mws) - 1; i >= 0; i-- {
			next = mws[i](next)
		}

		return next
	}
}
