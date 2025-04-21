package fakes

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"pvz/internal/models"
)

type FakeReceptionRepository struct {
	fakeDB         map[uuid.UUID]models.Reception
	fakeProductsDB map[uuid.UUID]models.Product
}

func NewFakeReceptionRepository() *FakeReceptionRepository {
	return &FakeReceptionRepository{
		fakeDB:         make(map[uuid.UUID]models.Reception),
		fakeProductsDB: make(map[uuid.UUID]models.Product),
	}
}

func (p *FakeReceptionRepository) CreateReception(ctx context.Context, reception models.Reception) error {
	if _, ok := p.fakeDB[reception.Id]; ok {
		return errors.New("reception already exists")
	}

	p.fakeDB[reception.Id] = reception
	return nil
}

func (p *FakeReceptionRepository) GetOpenReception(ctx context.Context, pvzId uuid.UUID) (models.Reception, error) {
	for _, reception := range p.fakeDB {
		if reception.PvzId == pvzId && reception.Status == models.InProgress {
			return reception, nil
		}
	}

	return models.Reception{}, errors.New("there is no open reception")
}

func (p *FakeReceptionRepository) AddProduct(ctx context.Context, product models.Product) error {
	p.fakeProductsDB[product.Id] = product
	return nil
}

func (p *FakeReceptionRepository) RemoveProduct(ctx context.Context, receptionId uuid.UUID) error {
	return nil
}

func (p *FakeReceptionRepository) CloseReception(ctx context.Context, receptionData models.Reception) error {
	receptionData.Status = models.Closed
	p.fakeDB[receptionData.Id] = receptionData

	return nil
}
