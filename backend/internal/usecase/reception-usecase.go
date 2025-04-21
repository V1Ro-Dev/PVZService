package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"pvz/config"
	"pvz/internal/delivery/forms"
	"pvz/internal/models"
	"pvz/pkg/logger"
)

var ReceptionNotOpened = errors.New("reception is not opened")

type ReceptionRepository interface {
	CreateReception(ctx context.Context, receptionForm models.Reception) error
	AddProduct(ctx context.Context, product models.Product) error
	GetOpenReception(ctx context.Context, pvzId uuid.UUID) (models.Reception, error)
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

func (rc *ReceptionService) AddProduct(ctx context.Context, productForm forms.ProductForm) (models.Product, error) {
	formattedStr := time.Now().Format(config.TimeStampLayout)
	dateTime, err := time.Parse(config.TimeStampLayout, formattedStr)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("Error parsing time: %s", err.Error()))
		return models.Product{}, err
	}

	product := models.Product{
		Id:          uuid.New(),
		DateTime:    dateTime,
		ProductType: productForm.Type,
		ReceptionId: uuid.UUID{},
	}

	reception, err := rc.receptionRepo.GetOpenReception(ctx, productForm.PvzId)
	if err != nil {
		return models.Product{}, err
	}

	product.ReceptionId = reception.Id

	err = rc.receptionRepo.AddProduct(ctx, product)
	if err != nil {
		return models.Product{}, err
	}

	return product, nil
}
