package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"pvz/config"
	"pvz/internal/delivery/forms"
	"pvz/internal/models"
	"pvz/internal/utils"
	"pvz/pkg/logger"
)

type PvzUseCase interface {
	CreatePvz(ctx context.Context, pvzForm forms.PvzForm) (models.Pvz, error)
	GetPvzInfo(ctx context.Context, form forms.GetPvzInfoForm) ([]models.PvzInfo, error)
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
	logger.Info(r.Context(), "Got Pvz creation request, trying to parse json")

	var pvzForm forms.PvzForm
	if err := json.NewDecoder(r.Body).Decode(&pvzForm); err != nil {
		logger.Error(r.Context(), fmt.Sprintf("Error decoding json: %s", err.Error()))
		utils.WriteJsonError(w, "failed to parse json", http.StatusBadRequest)
		return
	}

	logger.Info(r.Context(), "Successfully parsed json")

	if err := utils.ValidateCity(pvzForm.City); err != nil {
		logger.Error(r.Context(), fmt.Sprintf("City validation error: %s", err.Error()))
		utils.WriteJsonError(w, "Invalid city", http.StatusBadRequest)
		return
	}

	pvz, err := ph.pvzUseCase.CreatePvz(r.Context(), pvzForm)
	if err != nil {
		utils.WriteJsonError(w, "failed to create pvz", http.StatusBadRequest)
		return
	}

	logger.Info(r.Context(), fmt.Sprintf("Successfully created pvz with Id: %s", pvz.Id.String()))
	utils.WriteJson(w, forms.ToPvzForm(pvz), http.StatusCreated)
}

func (ph *PvzHandler) GetPvzInfo(w http.ResponseWriter, r *http.Request) {
	logger.Info(r.Context(), "Got Pvz info request, trying to parse query params")

	q := r.URL.Query()
	var page, limit int
	var start, end time.Time
	var err error

	if start, err = time.Parse(config.TimeStampLayout, q.Get("startDate")); err != nil {
		start = time.Time{}
	}

	if end, err = time.Parse(config.TimeStampLayout, q.Get("endDate")); err != nil {
		end = time.Now()
	}

	if page, err = strconv.Atoi(q.Get("page")); err != nil || page < 1 {
		page = 1
	}

	if limit, err = strconv.Atoi(q.Get("limit")); err != nil || limit < 1 {
		limit = 10
	}

	if err = utils.ValidateTime(start, end); err != nil {
		start, _ = time.Parse(config.TimeStampLayout, time.Time{}.String())
		end, _ = time.Parse(config.TimeStampLayout, time.Now().String())
	}

	pvzInfoForm := forms.GetPvzInfoForm{
		StartDate: start,
		EndDate:   end,
		Page:      page,
		Limit:     limit,
	}

	logger.Info(r.Context(), "Successfully parsed query params")
	logger.Info(r.Context(), pvzInfoForm)

	res, err := ph.pvzUseCase.GetPvzInfo(r.Context(), pvzInfoForm)
	if err != nil {
		utils.WriteJsonError(w, "unable to get info", http.StatusBadRequest)
		return
	}

	utils.WriteJson(w, forms.ToGetPvzInfoFormOut(res), http.StatusOK)

}
