package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"pvz/config"
	"pvz/internal/delivery/forms"
	"pvz/internal/models"
	"pvz/pkg/logger"
)

type ReceptionRepository interface {
	CreateReception(ctx context.Context, receptionForm models.Reception) error
}

type ReceptionService struct {
	receptionRepo ReceptionRepository
}

func NewReceptionService(receptionRepo ReceptionRepository) *ReceptionService {
	return &ReceptionService{
		receptionRepo: receptionRepo,
	}
}

func (rc *ReceptionService) CreateReception(ctx context.Context, receptionForm forms.ReceptionForm) (models.Reception, error) {
	formattedStr := time.Now().Format(config.TimeStampLayout)
	dateTime, err := time.Parse(config.TimeStampLayout, formattedStr)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("Error parsing time: %s", err.Error()))
		return models.Reception{}, err
	}

	reception := models.Reception{
		Id:       uuid.New(),
		DateTime: dateTime,
		PvzId:    receptionForm.PvzId,
		Status:   models.InProgress,
	}

	err = rc.receptionRepo.CreateReception(ctx, reception)
	if err != nil {
		return models.Reception{}, err
	}

	return reception, nil
}
