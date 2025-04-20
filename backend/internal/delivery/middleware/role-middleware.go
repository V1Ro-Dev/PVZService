package middleware

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"pvz/internal/utils"
)

func RoleMiddleware(allowedTypes ...string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headerValue := r.Header.Get("Authorization")
			tokenParts := strings.Split(headerValue, " ")

			if len(tokenParts) != 2 {
				utils.WriteJsonError(w, "invalid token", http.StatusBadRequest)
				return
			}

			role, err := utils.GetRole(tokenParts[1])
			if err != nil {
				utils.WriteJsonError(w, "Incorrect role or wrong token format", http.StatusForbidden)
				return
			}

			for _, allowed := range allowedTypes {
				if role == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}
			utils.WriteJsonError(w, "You don't have permission to use this endpoint", http.StatusForbidden)
		})
	}
}
