package usecase

import (
	"context"

	"pvz/internal/delivery/forms"
	"pvz/internal/models"
)

type PvzRepository interface {
	CreatePvz(ctx context.Context, pvzData models.Pvz) error
	GetPvzInfo(ctx context.Context, form forms.GetPvzInfoForm) ([]models.PvzInfo, error)
	GetPvzList(ctx context.Context) ([]models.Pvz, error)
}

type PvzService struct {
	pvzRepo PvzRepository
}

func NewPvzService(pvzRepo PvzRepository) *PvzService {
	return &PvzService{
		pvzRepo: pvzRepo,
	}
}

func (p *PvzService) CreatePvz(ctx context.Context, pvzForm forms.PvzForm) (models.Pvz, error) {
	pvzData := models.Pvz{
		Id:               pvzForm.Id,
		RegistrationDate: pvzForm.RegistrationDate,
		City:             pvzForm.City,
	}

	err := p.pvzRepo.CreatePvz(ctx, pvzData)
	if err != nil {
		return models.Pvz{}, err
	}

	return pvzData, nil
}

func (p *PvzService) GetPvzInfo(ctx context.Context, form forms.GetPvzInfoForm) ([]models.PvzInfo, error) {
	res, err := p.pvzRepo.GetPvzInfo(ctx, form)
	if err != nil {
		return []models.PvzInfo{}, err
	}

	return res, nil
}

func (p *PvzService) GetPvzList(ctx context.Context) ([]models.Pvz, error) {
	res, err := p.pvzRepo.GetPvzList(ctx)
	if err != nil {
		return []models.Pvz{}, err
	}

	return res, nil
}
