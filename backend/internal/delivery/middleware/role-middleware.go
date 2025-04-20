package middleware

import (
	"net/http"
	"pvz/pkg/logger"
	"strings"

	"github.com/gorilla/mux"

	"pvz/internal/utils"
)

func RoleMiddleware(allowedTypes ...string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := utils.SetRequestId(r.Context())

			headerValue := r.Header.Get("Authorization")
			tokenParts := strings.Split(headerValue, " ")

			if len(tokenParts) != 2 {
				logger.Error(ctx, "Invalid Authorization header params count")
				utils.WriteJsonError(w, "invalid token", http.StatusBadRequest)
				return
			}

			role, err := utils.GetRole(tokenParts[1])
			if err != nil {
				logger.Error(ctx, "Incorrect role or wrong token format")
				utils.WriteJsonError(w, "Incorrect role or wrong token format", http.StatusForbidden)
				return
			}

			for _, allowed := range allowedTypes {
				if role == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}
			logger.Error(ctx, "Role doesn't have permission")
			utils.WriteJsonError(w, "You don't have permission to use this endpoint", http.StatusForbidden)
		})
	}
}
