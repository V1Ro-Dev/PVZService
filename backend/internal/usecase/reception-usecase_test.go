package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"pvz/internal/delivery/forms"
	"pvz/internal/models"
	"pvz/internal/usecase"
	"pvz/internal/usecase/mocks"
)

func TestReceptionService_CreateReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockReceptionRepository(ctrl)
	service := usecase.NewReceptionService(mockRepo)

	form := forms.ReceptionForm{
		PvzId: uuid.New(),
	}

	tests := []struct {
		name    string
		mock    func()
		want    models.Reception
		wantErr bool
	}{
		{
			name: "success",
			mock: func() {
				mockRepo.EXPECT().CreateReception(gomock.Any(), gomock.Any()).Return(nil)
			},
			want:    models.Reception{Status: models.InProgress},
			wantErr: false,
		},
		{
			name: "repository error",
			mock: func() {
				mockRepo.EXPECT().CreateReception(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
			want:    models.Reception{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt.mock()
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.CreateReception(context.Background(), form)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want.Status, got.Status)
		})
	}
}

func TestReceptionService_AddProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockReceptionRepository(ctrl)
	service := usecase.NewReceptionService(mockRepo)

	form := forms.ProductForm{
		PvzId: uuid.New(),
		Type:  "Electronics",
	}

	tests := []struct {
		name    string
		mock    func()
		want    models.Product
		wantErr bool
	}{
		{
			name: "success",
			mock: func() {
				mockRepo.EXPECT().GetOpenReception(gomock.Any(), gomock.Any()).Return(models.Reception{
					Id:       uuid.New(),
					Status:   models.InProgress,
					PvzId:    uuid.New(),
					DateTime: time.Now(),
				}, nil)
				mockRepo.EXPECT().AddProduct(gomock.Any(), gomock.Any()).Return(nil)
			},
			want:    models.Product{ProductType: "Electronics"},
			wantErr: false,
		},
		{
			name: "reception not opened",
			mock: func() {
				mockRepo.EXPECT().GetOpenReception(gomock.Any(), gomock.Any()).Return(models.Reception{}, nil)
			},
			want:    models.Product{},
			wantErr: true,
		},
		{
			name: "repository error",
			mock: func() {
				mockRepo.EXPECT().GetOpenReception(gomock.Any(), gomock.Any()).Return(models.Reception{
					Id: uuid.New(),
				}, nil)
				mockRepo.EXPECT().AddProduct(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
			want:    models.Product{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt.mock()
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.AddProduct(context.Background(), form)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want.ProductType, got.ProductType)
		})
	}
}

func TestReceptionService_RemoveProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockReceptionRepository(ctrl)
	service := usecase.NewReceptionService(mockRepo)

	pvzId := uuid.New()

	tests := []struct {
		name    string
		mock    func()
		wantErr bool
	}{
		{
			name: "success",
			mock: func() {
				mockRepo.EXPECT().GetOpenReception(gomock.Any(), pvzId).Return(models.Reception{
					Id: uuid.New(),
				}, nil)
				mockRepo.EXPECT().RemoveProduct(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "reception not opened",
			mock: func() {
				mockRepo.EXPECT().GetOpenReception(gomock.Any(), pvzId).Return(models.Reception{}, nil)
			},
			wantErr: true,
		},
		{
			name: "repository error",
			mock: func() {
				mockRepo.EXPECT().GetOpenReception(gomock.Any(), pvzId).Return(models.Reception{
					Id: uuid.New(),
				}, nil)
				mockRepo.EXPECT().RemoveProduct(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt.mock()
		t.Run(tt.name, func(t *testing.T) {
			err := service.RemoveProduct(context.Background(), pvzId)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestReceptionService_CloseReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockReceptionRepository(ctrl)
	service := usecase.NewReceptionService(mockRepo)

	receptionId := uuid.New()
	dateTime := time.Now().Truncate(time.Millisecond)

	pvzId := uuid.New()

	tests := []struct {
		name    string
		mock    func()
		want    models.Reception
		wantErr bool
	}{
		{
			name: "success",
			mock: func() {
				mockRepo.EXPECT().GetOpenReception(gomock.Any(), pvzId).Return(models.Reception{
					Id:       receptionId,
					DateTime: dateTime,
					PvzId:    pvzId,
					Status:   models.InProgress,
				}, nil)
				mockRepo.EXPECT().CloseReception(gomock.Any(), gomock.Any()).Return(nil)
			},
			want: models.Reception{
				Id:       receptionId,
				DateTime: dateTime,
				PvzId:    pvzId,
				Status:   models.InProgress,
			},
			wantErr: false,
		},
		{
			name: "reception not opened",
			mock: func() {
				mockRepo.EXPECT().GetOpenReception(gomock.Any(), pvzId).Return(models.Reception{}, nil)
			},
			want:    models.Reception{},
			wantErr: true,
		},
		{
			name: "repository error",
			mock: func() {
				mockRepo.EXPECT().GetOpenReception(gomock.Any(), pvzId).Return(models.Reception{
					Id: uuid.New(),
				}, nil)
				mockRepo.EXPECT().CloseReception(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
			want:    models.Reception{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt.mock()
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.CloseReception(context.Background(), pvzId)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
