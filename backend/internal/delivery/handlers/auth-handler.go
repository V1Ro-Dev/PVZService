package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"pvz/internal/delivery/forms"
	"pvz/internal/models"
	"pvz/internal/utils"
	"pvz/pkg/logger"
)

type AuthUseCase interface {
	DummyLogin(ctx context.Context, role string) (string, error)
	CreateUser(ctx context.Context, signUpForm forms.SignUpFormIn) (models.User, error)
	IsUserExist(ctx context.Context, email string) (bool, error)
}

type AuthHandler struct {
	authUseCase AuthUseCase
}

func NewAuthHandler(authUseCase AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

func (a *AuthHandler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	ctx := utils.SetRequestId(r.Context())

	logger.Info(ctx, "Got dummy login request, trying to parse json")

	var dummyLoginForm forms.DummyLoginForm
	if err := json.NewDecoder(r.Body).Decode(&dummyLoginForm); err != nil {
		logger.Error(ctx, fmt.Sprintf("Error decoding json: %s", err.Error()))
		utils.WriteJsonError(w, "failed to parse json", http.StatusBadRequest)
		return
	}

	logger.Info(ctx, "Successfully parsed json")

	if utils.ValidateRole(dummyLoginForm.Role) == false {
		logger.Error(ctx, fmt.Sprintf("Role %s is not valid", dummyLoginForm.Role))
		utils.WriteJsonError(w, "Incorrect role was given", http.StatusBadRequest)
		return
	}

	token, err := a.authUseCase.DummyLogin(ctx, dummyLoginForm.Role)
	if err != nil {
		utils.WriteJsonError(w, "failed to gen token", http.StatusUnauthorized)
		return
	}

	utils.WriteJson(w, token, http.StatusOK)

	logger.Info(ctx, "Successfully processed dummy login request")
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := utils.SetRequestId(r.Context())

	logger.Info(ctx, "Got register request")

	var signUpForm forms.SignUpFormIn
	if err := json.NewDecoder(r.Body).Decode(&signUpForm); err != nil {
		logger.Error(ctx, fmt.Sprintf("Error decoding json: %s", err.Error()))
		utils.WriteJsonError(w, "failed to parse json", http.StatusBadRequest)
		return
	}

	logger.Info(ctx, "Successfully parsed json")

	if utils.ValidateRole(signUpForm.Role) == false {
		logger.Error(ctx, fmt.Sprintf("Role %s is not valid", signUpForm.Role))
		utils.WriteJsonError(w, "Incorrect role was given", http.StatusBadRequest)
		return
	}

	isExists, err := a.authUseCase.IsUserExist(ctx, signUpForm.Email)
	if err != nil {
		utils.WriteJsonError(w, "failed to check user exists", http.StatusBadRequest)
		return
	}
	
	if isExists {
		utils.WriteJsonError(w, "user already exists", http.StatusBadRequest)
		return
	}

	user, err := a.authUseCase.CreateUser(ctx, signUpForm)
	if err != nil {
		utils.WriteJsonError(w, "failed to create user", http.StatusBadRequest)
		return
	}

	logger.Info(ctx, fmt.Sprintf("Successfully created user with Id: %s", user.Id))
	utils.WriteJson(w, forms.ToSignUpOut(user), http.StatusCreated)
}
