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
	AddProduct(ctx context.Context, productForm forms.ProductForm) (models.Product, error)
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

func (rc *ReceptionHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	logger.Info(r.Context(), "Got add product request")

	var productForm forms.ProductForm
	err := json.NewDecoder(r.Body).Decode(&productForm)
	if err != nil {
		logger.Error(r.Context(), fmt.Sprintf("Error decoding json: %s", err.Error()))
		utils.WriteJsonError(w, "Error decoding json", http.StatusBadRequest)
		return
	}

	logger.Info(r.Context(), "Successfully parsed json")

	if err = utils.ValidateProductType(productForm.Type); err != nil {
		logger.Error(r.Context(), fmt.Sprintf("Error validating product type: %s", err.Error()))
		utils.WriteJsonError(w, "Product type not allowed", http.StatusBadRequest)
		return
	}

	product, err := rc.receptionUseCase.AddProduct(r.Context(), productForm)
	if err != nil {
		utils.WriteJsonError(w, "There are no active receptions or non-existing pvzId", http.StatusBadRequest)
		return
	}

	utils.WriteJson(w, forms.ToProductFormOut(product), http.StatusCreated)
}
