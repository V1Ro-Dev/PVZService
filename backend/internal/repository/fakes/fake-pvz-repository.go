package fakes

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"pvz/internal/delivery/forms"
	"pvz/internal/models"
)

type FakePvzRepository struct {
	fakeDB map[uuid.UUID]models.Pvz
}

func NewPostgresPvzRepository() *FakePvzRepository {
	return &FakePvzRepository{fakeDB: make(map[uuid.UUID]models.Pvz)}
}

func (p *FakePvzRepository) CreatePvz(ctx context.Context, pvzData models.Pvz) error {
	if _, ok := p.fakeDB[pvzData.Id]; ok {
		return errors.New("pvz already exists")
	}

	p.fakeDB[pvzData.Id] = pvzData
	return nil
}

func (p *FakePvzRepository) GetPvzInfo(ctx context.Context, form forms.GetPvzInfoForm) ([]models.PvzInfo, error) {
	return []models.PvzInfo{}, nil
}

func (p *FakePvzRepository) GetPvzList(ctx context.Context) ([]models.Pvz, error) {
	return []models.Pvz{}, nil
}
