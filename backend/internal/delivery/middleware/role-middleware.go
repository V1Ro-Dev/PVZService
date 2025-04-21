package middleware

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"

	"pvz/internal/models"
	"pvz/internal/utils"
	"pvz/pkg/logger"
)

func RoleMiddleware(allowedTypes ...models.Role) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headerValue := r.Header.Get("Authorization")
			tokenParts := strings.Split(headerValue, " ")

			if len(tokenParts) != 2 {
				logger.Error(r.Context(), "Invalid Authorization header params count")
				utils.WriteJsonError(w, "invalid token", http.StatusBadRequest)
				return
			}

			if tokenParts[0] != "Bearer" {
				logger.Error(r.Context(), fmt.Sprintf("Invalid first param: %s", tokenParts[0]))
			}

			role, err := utils.GetRole(tokenParts[1])
			if err != nil {
				logger.Error(r.Context(), fmt.Sprintf("Incorrect role or wrong token format: %s", err.Error()))
				utils.WriteJsonError(w, "Incorrect role or wrong token format", http.StatusForbidden)
				return
			}

			for _, allowed := range allowedTypes {
				if role == string(allowed) {
					next.ServeHTTP(w, r)
					return
				}
			}
			logger.Error(r.Context(), "Role doesn't have permission")
			utils.WriteJsonError(w, "You don't have permission to use this endpoint", http.StatusForbidden)
		})
	}
}
