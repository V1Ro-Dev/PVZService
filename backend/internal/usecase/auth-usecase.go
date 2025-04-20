package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"pvz/internal/delivery/forms"
	"pvz/internal/models"
	"pvz/internal/utils"
	"pvz/pkg/logger"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) error
	IsUserExist(ctx context.Context, email string) (bool, error)
	GetUserByEmail(ctx context.Context, logInData models.LoginData) (models.User, error)
}

type AuthService struct {
	userRepo UserRepository
}

func NewAuthService(userRepo UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (a *AuthService) DummyLogin(ctx context.Context, role string) (string, error) {
	logger.Info(ctx, "Trying to gen token")

	token, err := utils.GenerateToken(role)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("Error generating token: %s", err.Error()))
		return "", err
	}

	logger.Info(ctx, "Token was successfully generated")

	return token, nil
}

func (a *AuthService) CreateUser(ctx context.Context, in forms.SignUpFormIn) (models.User, error) {
	logger.Info(ctx, "Trying to create user")

	salt := utils.GenSalt()
	hashedPass := utils.HashPassword(in.Password, salt)
	user := models.User{
		Id:       uuid.NewString(),
		Email:    in.Email,
		Password: hashedPass,
		Salt:     salt,
		Role:     in.Role,
	}

	if err := a.userRepo.CreateUser(ctx, user); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (a *AuthService) IsUserExist(ctx context.Context, email string) (bool, error) {
	isExists, err := a.userRepo.IsUserExist(ctx, email)
	if err != nil {
		return false, err
	}

	return isExists, nil
}

func (a *AuthService) LogInUser(ctx context.Context, logInForm forms.LogInFormIn) (string, error) {
	loginData := models.LoginData{
		Email:    logInForm.Email,
		Password: logInForm.Password,
	}

	user, err := a.userRepo.GetUserByEmail(ctx, loginData)
	if err != nil {
		return "", err
	}

	if !utils.CheckPassword(loginData.Password, user.Password, user.Salt) {
		logger.Error(ctx, "Passwords don't match")
		return "", errors.New("wrong auth data")
	}

	return user.Role, nil
}
