package usecase_test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"pvz/internal/delivery/forms"
	"pvz/internal/models"
	"pvz/internal/usecase"
	"pvz/internal/usecase/mocks"
)

func TestPvzService_CreatePvz(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzId := uuid.New()
	regDate := time.Now().Truncate(time.Millisecond)

	mockRepo := mocks.NewMockPvzRepository(ctrl)
	service := usecase.NewPvzService(mockRepo)

	tests := []struct {
		name    string
		input   forms.PvzForm
		mock    func()
		want    models.Pvz
		wantErr bool
	}{
		{
			name: "success",
			input: forms.PvzForm{
				Id:               pvzId,
				RegistrationDate: regDate,
				City:             "Москва",
			},
			mock: func() {
				mockRepo.EXPECT().CreatePvz(gomock.Any(), gomock.Any()).Return(nil)
			},
			want: models.Pvz{
				Id:               pvzId,
				RegistrationDate: regDate,
				City:             "Москва",
			},
			wantErr: false,
		},
		{
			name: "repository error",
			input: forms.PvzForm{
				Id:               pvzId,
				RegistrationDate: regDate,
				City:             "Казань",
			},
			mock: func() {
				mockRepo.EXPECT().CreatePvz(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
			want:    models.Pvz{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt.mock()
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.CreatePvz(context.Background(), tt.input)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestPvzService_GetPvzInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPvzRepository(ctrl)
	service := usecase.NewPvzService(mockRepo)

	startDate := time.Now().Truncate(time.Millisecond)
	endDate := time.Now().Truncate(time.Millisecond)

	pvzId := uuid.New()
	receptionId := uuid.New()
	productId := uuid.New()

	form := forms.GetPvzInfoForm{
		StartDate: startDate,
		EndDate:   endDate,
		Page:      1,
		Limit:     10,
	}

	tests := []struct {
		name    string
		mock    func()
		want    []models.PvzInfo
		wantErr bool
	}{
		{
			name: "success",
			mock: func() {
				mockRepo.EXPECT().GetPvzInfo(gomock.Any(), form).Return([]models.PvzInfo{
					{
						Pvz: models.Pvz{
							Id:               pvzId,
							RegistrationDate: time.Now().Truncate(time.Millisecond),
							City:             "Москва",
						},
						Receptions: []models.ReceptionProducts{
							{
								Reception: models.Reception{
									Id:       receptionId,
									DateTime: time.Now().Truncate(time.Millisecond),
									PvzId:    pvzId,
									Status:   models.InProgress,
								},
								Products: []models.Product{
									{
										Id:          productId,
										DateTime:    time.Now().Truncate(time.Millisecond),
										ProductType: "электроника",
										ReceptionId: receptionId,
									},
								},
							},
						},
					},
				}, nil)
			},
			want:    []models.PvzInfo{{}},
			wantErr: false,
		},
		{
			name: "repository error",
			mock: func() {
				mockRepo.EXPECT().GetPvzInfo(gomock.Any(), form).Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt.mock()
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetPvzInfo(context.Background(), form)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestPvzService_GetPvzList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPvzRepository(ctrl)
	service := usecase.NewPvzService(mockRepo)

	pvzIdFirst := uuid.New()
	pvzIdSecond := uuid.New()

	tests := []struct {
		name    string
		mock    func()
		want    []models.Pvz
		wantErr bool
	}{
		{
			name: "success",
			mock: func() {
				mockRepo.EXPECT().GetPvzList(gomock.Any()).Return([]models.Pvz{
					{Id: pvzIdFirst, City: "Москва"},
					{Id: pvzIdSecond, City: "Казань"},
				}, nil)
			},
			want: []models.Pvz{
				{Id: pvzIdFirst, City: "Москва"},
				{Id: pvzIdSecond, City: "Казань"},
			},
			wantErr: false,
		},
		{
			name: "repository error",
			mock: func() {
				mockRepo.EXPECT().GetPvzList(gomock.Any()).Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt.mock()
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetPvzList(context.Background())
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
