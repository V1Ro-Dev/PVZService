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
	LogInUser(ctx context.Context, logInForm forms.LogInFormIn) (string, error)
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
	logger.Info(r.Context(), "Got dummy login request, trying to parse json")

	var dummyLoginForm forms.DummyLoginForm
	if err := json.NewDecoder(r.Body).Decode(&dummyLoginForm); err != nil {
		logger.Error(r.Context(), fmt.Sprintf("Error decoding json: %s", err.Error()))
		utils.WriteJsonError(w, "failed to parse json", http.StatusBadRequest)
		return
	}

	logger.Info(r.Context(), "Successfully parsed json")

	if utils.ValidateRole(dummyLoginForm.Role) == false {
		logger.Error(r.Context(), fmt.Sprintf("Role %s is not valid", dummyLoginForm.Role))
		utils.WriteJsonError(w, "Incorrect role was given", http.StatusBadRequest)
		return
	}

	token, err := a.authUseCase.DummyLogin(r.Context(), dummyLoginForm.Role)
	if err != nil {
		utils.WriteJsonError(w, "failed to gen.bat token", http.StatusBadRequest)
		return
	}

	utils.WriteJson(w, token, http.StatusOK)

	logger.Info(r.Context(), "Successfully processed dummy login request")
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	logger.Info(r.Context(), "Got register request, trying to parse json")

	var signUpForm forms.SignUpFormIn
	if err := json.NewDecoder(r.Body).Decode(&signUpForm); err != nil {
		logger.Error(r.Context(), fmt.Sprintf("Error decoding json: %s", err.Error()))
		utils.WriteJsonError(w, "failed to parse json", http.StatusBadRequest)
		return
	}

	logger.Info(r.Context(), "Successfully parsed json")

	if err := utils.ValidateAll(signUpForm.Email, signUpForm.Role); err != nil {
		logger.Error(r.Context(), err.Error())
		utils.WriteJsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	isExists, err := a.authUseCase.IsUserExist(r.Context(), signUpForm.Email)
	if err != nil {
		utils.WriteJsonError(w, "failed to check user exists", http.StatusBadRequest)
		return
	}

	if isExists {
		utils.WriteJsonError(w, "user already exists", http.StatusBadRequest)
		return
	}

	user, err := a.authUseCase.CreateUser(r.Context(), signUpForm)
	if err != nil {
		utils.WriteJsonError(w, "failed to create user", http.StatusBadRequest)
		return
	}

	logger.Info(r.Context(), fmt.Sprintf("Successfully created user with Id: %s", user.Id))
	utils.WriteJson(w, forms.ToSignUpOut(user), http.StatusCreated)
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	logger.Info(r.Context(), "Got login request")

	var logInForm forms.LogInFormIn
	if err := json.NewDecoder(r.Body).Decode(&logInForm); err != nil {
		logger.Error(r.Context(), fmt.Sprintf("Error decoding json: %s", err.Error()))
		utils.WriteJsonError(w, "failed to parse json", http.StatusBadRequest)
		return
	}

	role, err := a.authUseCase.LogInUser(r.Context(), logInForm)
	if err != nil {
		utils.WriteJsonError(w, "Wrong auth data", http.StatusUnauthorized)
		return
	}

	token, err := a.authUseCase.DummyLogin(r.Context(), role)
	if err != nil {
		utils.WriteJsonError(w, "failed to gen.bat token", http.StatusUnauthorized)
		return
	}

	utils.WriteJson(w, token, http.StatusOK)

	logger.Info(r.Context(), "Successfully processed login request")
}
