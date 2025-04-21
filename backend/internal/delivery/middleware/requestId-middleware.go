package middleware

import (
	"net/http"

	"pvz/internal/utils"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := utils.SetRequestId(r.Context())
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
