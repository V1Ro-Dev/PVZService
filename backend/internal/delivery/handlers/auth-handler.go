package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"pvz/internal/delivery/forms"
	"pvz/internal/utils"
	"pvz/pkg/logger"
)

type AuthUseCase interface {
	dummyLogin(ctx context.Context, role string) (string, error)
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

	token, err := a.authUseCase.dummyLogin(ctx, dummyLoginForm.Role)
	if err != nil {
		utils.WriteJsonError(w, "failed to gen token", http.StatusUnauthorized)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
	})

	utils.WriteJson(w, token, http.StatusOK)

	logger.Info(ctx, "Successfully processed dummy login request")
}
