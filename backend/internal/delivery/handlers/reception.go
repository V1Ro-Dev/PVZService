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

type ReceptionUseCase interface {
	CreateReception(ctx context.Context, receptionForm forms.ReceptionForm) (models.Reception, error)
}

type ReceptionHandler struct {
	receptionUseCase ReceptionUseCase
}

func NewReceptionHandler(receptionUseCase ReceptionUseCase) *ReceptionHandler {
	return &ReceptionHandler{
		receptionUseCase: receptionUseCase,
	}
}

func (rc *ReceptionHandler) CreateReception(w http.ResponseWriter, r *http.Request) {
	logger.Info(r.Context(), "Got create reception request")

	var receptionForm forms.ReceptionForm
	err := json.NewDecoder(r.Body).Decode(&receptionForm)
	if err != nil {
		logger.Error(r.Context(), fmt.Sprintf("Error decoding json: %s", err.Error()))
		utils.WriteJsonError(w, "Error decoding json", http.StatusBadRequest)
		return
	}

	logger.Info(r.Context(), "Successfully parsed json")

	reception, err := rc.receptionUseCase.CreateReception(r.Context(), receptionForm)
	if err != nil {
		utils.WriteJsonError(w, "unclosed reception or non-existing pvzId", http.StatusBadRequest)
		return
	}

	utils.WriteJson(w, forms.ToReceptionFormOut(reception), http.StatusCreated)
}
