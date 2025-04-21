package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"pvz/internal/delivery/forms"
	"pvz/internal/delivery/handlers"
	"pvz/internal/delivery/mocks"
	"pvz/internal/models"
)

func TestDummyLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockAuthUseCase(ctrl)
	handler := handlers.NewAuthHandler(mockUC)

	tests := []struct {
		name         string
		input        forms.DummyLoginForm
		mockToken    string
		mockErr      error
		expectStatus int
		expectBody   string
	}{
		{
			name:         "valid moderator role",
			input:        forms.DummyLoginForm{Role: string(models.Moderator)},
			mockToken:    "token123",
			mockErr:      nil,
			expectStatus: http.StatusOK,
			expectBody:   `"token123"`,
		},
		{
			name:         "valid client role",
			input:        forms.DummyLoginForm{Role: string(models.Client)},
			mockToken:    "token123",
			mockErr:      nil,
			expectStatus: http.StatusOK,
			expectBody:   `"token123"`,
		},
		{
			name:         "valid employee role",
			input:        forms.DummyLoginForm{Role: string(models.Employee)},
			mockToken:    "token123",
			mockErr:      nil,
			expectStatus: http.StatusOK,
			expectBody:   `"token123"`,
		},
		{
			name:         "invalid role",
			input:        forms.DummyLoginForm{Role: "invalid_role"},
			expectStatus: http.StatusBadRequest,
			expectBody:   `{"message":"Incorrect role was given"}`,
		},
		{
			name:         "usecase error",
			input:        forms.DummyLoginForm{Role: string(models.Client)},
			mockErr:      errors.New("fail"),
			expectStatus: http.StatusBadRequest,
			expectBody:   `{"message":"failed to gen token"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockErr == nil && tt.expectStatus == http.StatusOK {
				mockUC.EXPECT().DummyLogin(gomock.Any(), tt.input.Role).Return(tt.mockToken, nil)
			} else if tt.mockErr != nil && tt.input.Role != "invalid_role" {
				mockUC.EXPECT().DummyLogin(gomock.Any(), tt.input.Role).Return("", tt.mockErr)
			}

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/dummy-login", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			handler.DummyLogin(rec, req)

			assert.Equal(t, tt.expectStatus, rec.Code)

			respBody := rec.Body.String()
			assert.JSONEq(t, tt.expectBody, respBody)
		})
	}
}

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockAuthUseCase(ctrl)
	handler := handlers.NewAuthHandler(mockUC)

	validUser := models.User{
		Id:    "123",
		Email: "test@test.ru",
		Role:  string(models.Client),
	}

	tests := []struct {
		name         string
		input        forms.SignUpFormIn
		isExist      bool
		createErr    error
		mockErr      error
		expectStatus int
		expectBody   string
	}{
		{
			name: "valid signup",
			input: forms.SignUpFormIn{
				Email:    "test@test.ru",
				Password: "123",
				Role:     string(models.Client),
			},
			isExist:      false,
			createErr:    nil,
			expectStatus: http.StatusCreated,
			expectBody:   `{"id":"123","email":"test@test.ru","role":"client"}`,
		},
		{
			name: "User already exists",
			input: forms.SignUpFormIn{
				Email:    "test@test.ru",
				Password: "123",
				Role:     string(models.Client),
			},
			isExist:      true,
			expectStatus: http.StatusBadRequest,
			expectBody:   `{"message":"user already exists"}`,
		},
		{
			name: "invalid email",
			input: forms.SignUpFormIn{
				Email:    "badmail",
				Password: "123",
				Role:     string(models.Client),
			},
			expectStatus: http.StatusBadRequest,
			expectBody:   `{"message":"email badmail is not valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectStatus == http.StatusCreated {
				mockUC.EXPECT().IsUserExist(gomock.Any(), tt.input.Email).Return(false, nil)
				mockUC.EXPECT().CreateUser(gomock.Any(), tt.input).Return(validUser, nil)
			} else if tt.isExist {
				mockUC.EXPECT().IsUserExist(gomock.Any(), tt.input.Email).Return(true, nil)
			} else if tt.expectStatus == http.StatusBadRequest && tt.input.Email != "badmail" {
				mockUC.EXPECT().IsUserExist(gomock.Any(), tt.input.Email).Return(false, nil)
				mockUC.EXPECT().CreateUser(gomock.Any(), tt.input).Return(models.User{}, errors.New("fail"))
			}

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			handler.Register(rec, req)

			assert.Equal(t, tt.expectStatus, rec.Code)

			respBody := rec.Body.String()
			assert.JSONEq(t, tt.expectBody, respBody)
		})
	}
}

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockAuthUseCase(ctrl)
	handler := handlers.NewAuthHandler(mockUC)

	tests := []struct {
		name         string
		input        forms.LogInFormIn
		role         string
		loginErr     error
		token        string
		tokenErr     error
		expectStatus int
		expectBody   string
	}{
		{
			name: "valid login",
			input: forms.LogInFormIn{
				Email:    "email@test.com",
				Password: "pass",
			},
			role:         string(models.Moderator),
			token:        "token",
			expectStatus: http.StatusOK,
			expectBody:   `"token"`,
		},
		{
			name: "invalid credentials",
			input: forms.LogInFormIn{
				Email:    "email@test.com",
				Password: "wrong",
			},
			loginErr:     errors.New("wrong"),
			expectStatus: http.StatusUnauthorized,
			expectBody:   `{"message":"Wrong auth data"}`,
		},
		{
			name: "token generation error",
			input: forms.LogInFormIn{
				Email:    "email@test.com",
				Password: "pass",
			},
			role:         string(models.Client),
			tokenErr:     errors.New("token gen fail"),
			expectStatus: http.StatusUnauthorized,
			expectBody:   `{"message":"failed to gen token"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.loginErr == nil {
				mockUC.EXPECT().LogInUser(gomock.Any(), tt.input).Return(tt.role, nil)
			} else {
				mockUC.EXPECT().LogInUser(gomock.Any(), tt.input).Return("", tt.loginErr)
			}

			if tt.tokenErr == nil && tt.expectStatus == http.StatusOK {
				mockUC.EXPECT().DummyLogin(gomock.Any(), tt.role).Return(tt.token, nil)
			} else if tt.tokenErr != nil {
				mockUC.EXPECT().DummyLogin(gomock.Any(), tt.role).Return("", tt.tokenErr)
			}

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			handler.Login(rec, req)

			assert.Equal(t, tt.expectStatus, rec.Code)

			respBody := rec.Body.String()
			assert.JSONEq(t, tt.expectBody, respBody)
		})
	}
}
