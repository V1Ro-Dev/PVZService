package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"pvz/internal/delivery/forms"
	"pvz/internal/models"
	"pvz/internal/usecase"
	"pvz/internal/usecase/mocks"
	"pvz/internal/utils"
)

func TestAuthService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := usecase.NewAuthService(mockRepo)

	tests := []struct {
		name    string
		input   forms.SignUpFormIn
		mock    func()
		wantErr bool
		want    *models.User
	}{
		{
			name: "success",
			input: forms.SignUpFormIn{
				Email:    "test@example.com",
				Password: "secure",
				Role:     "admin",
			},
			mock: func() {
				mockRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
			want:    &models.User{Email: "test@example.com", Role: "admin"},
		},
		{
			name: "repository error",
			input: forms.SignUpFormIn{
				Email:    "fail@example.com",
				Password: "password",
				Role:     "user",
			},
			mock: func() {
				mockRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: true,
			want:    nil,
		},
	}

	for _, tt := range tests {
		tt.mock()
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.CreateUser(context.Background(), tt.input)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.Equal(t, tt.want.Email, got.Email)
				assert.Equal(t, tt.want.Role, got.Role)
			}
		})
	}
}

func TestAuthService_IsUserExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := usecase.NewAuthService(mockRepo)

	tests := []struct {
		name    string
		email   string
		mock    func()
		want    bool
		wantErr bool
	}{
		{
			name:  "user exists",
			email: "user@example.com",
			mock: func() {
				mockRepo.EXPECT().IsUserExist(gomock.Any(), "user@example.com").Return(true, nil)
			},
			want:    true,
			wantErr: false,
		},
		{
			name:  "user does not exist",
			email: "nouser@example.com",
			mock: func() {
				mockRepo.EXPECT().IsUserExist(gomock.Any(), "nouser@example.com").Return(false, nil)
			},
			want:    false,
			wantErr: false,
		},
		{
			name:  "repository error",
			email: "error@example.com",
			mock: func() {
				mockRepo.EXPECT().IsUserExist(gomock.Any(), "error@example.com").Return(false, errors.New("error"))
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt.mock()
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.IsUserExist(context.Background(), tt.email)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestAuthService_LogInUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := usecase.NewAuthService(mockRepo)

	rawPassword := "password123"
	salt := utils.GenSalt()
	hashed := utils.HashPassword(rawPassword, salt)

	tests := []struct {
		name    string
		input   forms.LogInFormIn
		mock    func()
		want    string
		wantErr bool
	}{
		{
			name: "success login",
			input: forms.LogInFormIn{
				Email:    "login@example.com",
				Password: rawPassword,
			},
			mock: func() {
				mockRepo.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(models.User{
					Email:    "login@example.com",
					Password: hashed,
					Salt:     salt,
					Role:     "user",
				}, nil)
			},
			want:    "user",
			wantErr: false,
		},
		{
			name: "wrong password",
			input: forms.LogInFormIn{
				Email:    "login@example.com",
				Password: "wrongpass",
			},
			mock: func() {
				mockRepo.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(models.User{
					Email:    "login@example.com",
					Password: hashed,
					Salt:     salt,
					Role:     "user",
				}, nil)
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "repo error",
			input: forms.LogInFormIn{
				Email:    "fail@example.com",
				Password: rawPassword,
			},
			mock: func() {
				mockRepo.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(models.User{}, errors.New("db fail"))
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt.mock()
		t.Run(tt.name, func(t *testing.T) {
			role, err := service.LogInUser(context.Background(), tt.input)
			assert.Equal(t, tt.want, role)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestAuthService_DummyLogin(t *testing.T) {
	service := usecase.NewAuthService(nil)

	tests := []struct {
		name    string
		role    string
		wantErr bool
	}{
		{
			name:    "valid token gen.bat",
			role:    "employee",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := service.DummyLogin(context.Background(), tt.role)
			if tt.wantErr {
				assert.Empty(t, token)
			} else {
				assert.NotEmpty(t, token)
			}
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
