package fakes

import (
	"context"
	"pvz/internal/models"
)

type FakeUserRepository struct {
}

func NewFakeUserRepository() *FakeUserRepository {
	return &FakeUserRepository{}
}

func (f *FakeUserRepository) CreateUser(ctx context.Context, user models.User) error {
	return nil
}

func (f *FakeUserRepository) IsUserExist(ctx context.Context, email string) (bool, error) {
	return true, nil
}
func (f *FakeUserRepository) GetUserByEmail(ctx context.Context, logInData models.LoginData) (models.User, error) {
	return models.User{}, nil
}
