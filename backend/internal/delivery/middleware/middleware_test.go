package middleware_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"pvz/internal/delivery/middleware"
	"pvz/internal/models"
	"pvz/internal/utils"
	"pvz/pkg/logger"
)

func TestRequestIDMiddleware(t *testing.T) {
	var called bool
	var requestIDFromContext string

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		// Извлекаем request_id из контекста
		id := GetRequestId(r.Context())
		requestIDFromContext = id

		w.WriteHeader(http.StatusOK)
	})

	handlerToTest := middleware.RequestIDMiddleware(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handlerToTest.ServeHTTP(rec, req)

	assert.True(t, called, "next handler should be called")
	assert.NotEmpty(t, requestIDFromContext, "request ID should be set in context")
	assert.Equal(t, http.StatusOK, rec.Code)
}

func GetRequestId(ctx context.Context) string {
	if id, ok := ctx.Value(logger.RequestID).(string); ok {
		return id
	}
	return ""
}

func TestRoleMiddleware(t *testing.T) {
	originalGetRole := utils.GetRole
	defer func() { utils.GetRole = originalGetRole }() // восстановим после теста

	type testCase struct {
		name           string
		authHeader     string
		mockGetRole    func(token string) (string, error)
		allowedRoles   []models.Role
		expectStatus   int
		expectResponse string
	}

	tests := []testCase{
		{
			name:         "valid token and role allowed",
			authHeader:   "Bearer token",
			mockGetRole:  func(token string) (string, error) { return "client", nil },
			allowedRoles: []models.Role{models.Client},
			expectStatus: http.StatusOK,
		},
		{
			name:           "invalid token format",
			authHeader:     "Bearer",
			mockGetRole:    nil, // не будет вызова
			allowedRoles:   []models.Role{models.Client},
			expectStatus:   http.StatusBadRequest,
			expectResponse: `{"message":"invalid token"}`,
		},
		{
			name:           "role not allowed",
			authHeader:     "Bearer token",
			mockGetRole:    func(token string) (string, error) { return "admin", nil },
			allowedRoles:   []models.Role{models.Client},
			expectStatus:   http.StatusForbidden,
			expectResponse: `{"message":"You don't have permission to use this endpoint"}`,
		},
		{
			name:           "token parse error",
			authHeader:     "Bearer token",
			mockGetRole:    func(token string) (string, error) { return "", errors.New("bad token") },
			allowedRoles:   []models.Role{models.Client},
			expectStatus:   http.StatusForbidden,
			expectResponse: `{"message":"Incorrect role or wrong token format"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockGetRole != nil {
				utils.GetRole = tt.mockGetRole
			}

			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			})

			middleware := middleware.RoleMiddleware(tt.allowedRoles...)(next)

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", tt.authHeader)
			rec := httptest.NewRecorder()

			middleware.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectStatus, rec.Code)
			if tt.expectStatus != http.StatusOK {
				assert.JSONEq(t, tt.expectResponse, rec.Body.String())
			} else {
				assert.True(t, nextCalled, "Next handler should be called")
			}
		})
	}
}
