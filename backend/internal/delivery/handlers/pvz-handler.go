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

type PvzUseCase interface {
	CreatePvz(ctx context.Context, pvzForm forms.PvzForm) (models.Pvz, error)
}

type PvzHandler struct {
	pvzUseCase PvzUseCase
}

func NewPvzHandler(pvzUseCase PvzUseCase) *PvzHandler {
	return &PvzHandler{
		pvzUseCase: pvzUseCase,
	}
}

func (ph *PvzHandler) CreatePvz(w http.ResponseWriter, r *http.Request) {
	ctx := utils.SetRequestId(r.Context())
	logger.Info(ctx, "Got Pvz creation request, trying to parse json")

	var pvzForm forms.PvzForm
	if err := json.NewDecoder(r.Body).Decode(&pvzForm); err != nil {
		logger.Error(ctx, fmt.Sprintf("Error decoding json: %s", err.Error()))
		utils.WriteJsonError(w, "failed to parse json", http.StatusBadRequest)
		return
	}

	if err := utils.ValidateCity(pvzForm.City); err != nil {
		logger.Error(ctx, fmt.Sprintf("City validation error: %s", err.Error()))
		utils.WriteJsonError(w, "Invalid city", http.StatusBadRequest)
		return
	}

	pvz, err := ph.pvzUseCase.CreatePvz(ctx, pvzForm)
	if err != nil {
		utils.WriteJsonError(w, "failed to create pvz", http.StatusBadRequest)
		return
	}

	logger.Info(ctx, fmt.Sprintf("Successfully created pvz with Id: %s", pvz.Id.String()))
	utils.WriteJson(w, forms.ToPvzForm(pvz), http.StatusCreated)
}
