package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	"pvz/internal/delivery/forms"
	"pvz/internal/delivery/handlers"
	"pvz/internal/delivery/mocks"
	"pvz/internal/models"
)

func TestReceptionHandler_CreateReception(t *testing.T) {
	pvzId := uuid.New()
	receptionId := uuid.New()
	tests := []struct {
		name        string
		input       forms.ReceptionForm
		mockReturn  models.Reception
		mockError   error
		wantStatus  int
		wantBodyOut forms.ReceptionFormOut
	}{
		{
			name:        "ok",
			input:       forms.ReceptionForm{PvzId: pvzId},
			mockReturn:  models.Reception{Id: receptionId, PvzId: pvzId, Status: models.InProgress},
			mockError:   nil,
			wantStatus:  http.StatusCreated,
			wantBodyOut: forms.ReceptionFormOut{Id: receptionId, PvzId: pvzId, Status: string(models.InProgress)},
		},
		{
			name:        "invalid json",
			input:       forms.ReceptionForm{},
			mockError:   nil,
			wantStatus:  http.StatusBadRequest,
			wantBodyOut: forms.ReceptionFormOut{},
		},
		{
			name:        "usecase error",
			input:       forms.ReceptionForm{PvzId: pvzId},
			mockReturn:  models.Reception{},
			mockError:   errors.New("some error"),
			wantStatus:  http.StatusBadRequest,
			wantBodyOut: forms.ReceptionFormOut{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := mocks.NewMockReceptionUseCase(ctrl)
			h := handlers.NewReceptionHandler(mockUseCase)

			var body []byte
			var err error

			if tt.name == "invalid json" {
				body = []byte(`invalid`)
			} else {
				body, err = json.Marshal(tt.input)
				require.NoError(t, err)

				mockUseCase.EXPECT().
					CreateReception(gomock.Any(), tt.input).
					Return(tt.mockReturn, tt.mockError).
					Times(1)
			}

			req := httptest.NewRequest(http.MethodPost, "/reception", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			h.CreateReception(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			require.Equal(t, tt.wantStatus, res.StatusCode)

			var out forms.ReceptionFormOut
			err = json.NewDecoder(res.Body).Decode(&out)
			if tt.wantStatus == http.StatusCreated {
				require.NoError(t, err)
				require.Equal(t, tt.wantBodyOut, out)
			} else {
				require.NoError(t, err)
				require.Equal(t, forms.ReceptionFormOut{}, out)
			}
		})
	}
}

func TestReceptionHandler_AddProduct(t *testing.T) {
	pvzId := uuid.New()
	productId := uuid.New()
	receptionId := uuid.New()

	tests := []struct {
		name        string
		input       forms.ProductForm
		mockReturn  models.Product
		mockError   error
		wantStatus  int
		wantBodyOut forms.ProductFormOut
	}{
		{
			name:        "ok",
			input:       forms.ProductForm{PvzId: pvzId, Type: "обувь"},
			mockReturn:  models.Product{Id: productId, ReceptionId: receptionId, ProductType: "обувь", DateTime: time.Now().Truncate(time.Millisecond)},
			mockError:   nil,
			wantStatus:  http.StatusCreated,
			wantBodyOut: forms.ProductFormOut{Id: productId, ReceptionId: receptionId, ProductType: "обувь", DateTime: time.Now().Truncate(time.Millisecond)},
		},
		{
			name:        "invalid json",
			input:       forms.ProductForm{},
			mockError:   nil,
			wantStatus:  http.StatusBadRequest,
			wantBodyOut: forms.ProductFormOut{},
		},
		{
			name:        "usecase error",
			input:       forms.ProductForm{PvzId: uuid.New(), Type: "обувь"},
			mockReturn:  models.Product{},
			mockError:   errors.New("some error"),
			wantStatus:  http.StatusBadRequest,
			wantBodyOut: forms.ProductFormOut{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := mocks.NewMockReceptionUseCase(ctrl)
			h := handlers.NewReceptionHandler(mockUseCase)

			body, _ := json.Marshal(tt.input)

			if tt.name != "invalid json" {
				mockUseCase.EXPECT().
					AddProduct(gomock.Any(), tt.input).
					Return(tt.mockReturn, tt.mockError).
					Times(1)
			}

			req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			h.AddProduct(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			require.Equal(t, tt.wantStatus, res.StatusCode)

			var out forms.ProductFormOut
			err := json.NewDecoder(res.Body).Decode(&out)
			require.NoError(t, err)
			require.Equal(t, tt.wantBodyOut, out)
		})
	}
}

func TestReceptionHandler_RemoveProduct(t *testing.T) {
	tests := []struct {
		name       string
		pvzID      uuid.UUID
		mockError  error
		wantStatus int
	}{
		{
			name:       "ok",
			pvzID:      uuid.New(),
			mockError:  nil,
			wantStatus: http.StatusOK,
		},
		{
			name:       "usecase error",
			pvzID:      uuid.New(),
			mockError:  errors.New("some error"),
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := mocks.NewMockReceptionUseCase(ctrl)
			h := handlers.NewReceptionHandler(mockUseCase)

			mockUseCase.EXPECT().
				RemoveProduct(gomock.Any(), tt.pvzID).
				Return(tt.mockError).
				Times(1)

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/pvz/%s/delete_last_product", tt.pvzID.String()), nil)
			req = mux.SetURLVars(req, map[string]string{"pvzId": tt.pvzID.String()})
			rec := httptest.NewRecorder()

			h.RemoveProduct(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			require.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestReceptionHandler_CloseReception(t *testing.T) {
	receptionId := uuid.New()
	pvzId := uuid.New()
	tests := []struct {
		name        string
		pvzID       uuid.UUID
		mockReturn  models.Reception
		mockError   error
		wantStatus  int
		wantBodyOut forms.ReceptionFormOut
	}{
		{
			name:        "ok",
			pvzID:       pvzId,
			mockReturn:  models.Reception{Id: receptionId, PvzId: pvzId, Status: models.InProgress},
			mockError:   nil,
			wantStatus:  http.StatusOK,
			wantBodyOut: forms.ReceptionFormOut{Id: receptionId, PvzId: pvzId, Status: string(models.InProgress)},
		},
		{
			name:        "usecase error",
			pvzID:       pvzId,
			mockReturn:  models.Reception{},
			mockError:   errors.New("some error"),
			wantStatus:  http.StatusBadRequest,
			wantBodyOut: forms.ReceptionFormOut{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := mocks.NewMockReceptionUseCase(ctrl)
			h := handlers.NewReceptionHandler(mockUseCase)

			mockUseCase.EXPECT().
				CloseReception(gomock.Any(), tt.pvzID).
				Return(tt.mockReturn, tt.mockError).
				Times(1)

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/pvz/%s/close_last_reception", tt.pvzID.String()), nil)
			req = mux.SetURLVars(req, map[string]string{"pvzId": tt.pvzID.String()})
			rec := httptest.NewRecorder()

			h.CloseReception(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			require.Equal(t, tt.wantStatus, res.StatusCode)

			var out forms.ReceptionFormOut
			err := json.NewDecoder(res.Body).Decode(&out)
			require.NoError(t, err)
			require.Equal(t, tt.wantBodyOut, out)
		})
	}
}
