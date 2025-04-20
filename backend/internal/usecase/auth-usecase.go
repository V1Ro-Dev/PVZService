package usecase

import (
	"context"
	"fmt"

	"pvz/internal/utils"
	"pvz/pkg/logger"
)

type UserRepository interface {
}

type AuthService struct {
	userRepo UserRepository
}

func NewAuthService(userRepo UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (a *AuthService) dummyLogin(ctx context.Context, role string) (string, error) {
	logger.Info(ctx, "Trying to gen token")

	token, err := utils.GenerateToken(role)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("Error generating token: %s", err.Error()))
		return "", err
	}

	logger.Info(ctx, "Token was successfully generated")

	return token, nil
}
