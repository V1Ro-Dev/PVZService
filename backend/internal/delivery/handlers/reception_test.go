package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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
	productType := "обувь"
	receptionId := uuid.New()
	now := time.Now().Truncate(time.Millisecond)

	tests := []struct {
		name        string
		body        io.Reader
		expectCall  bool
		input       forms.ProductForm
		mockReturn  models.Product
		mockError   error
		wantStatus  int
		wantBodyOut forms.ProductFormOut
	}{
		{
			name:        "ok",
			body:        toJSONBody(forms.ProductForm{PvzId: pvzId, Type: productType}),
			expectCall:  true,
			input:       forms.ProductForm{PvzId: pvzId, Type: productType},
			mockReturn:  models.Product{Id: productId, ReceptionId: receptionId, ProductType: productType, DateTime: now},
			mockError:   nil,
			wantStatus:  http.StatusCreated,
			wantBodyOut: forms.ProductFormOut{Id: productId, ReceptionId: receptionId, ProductType: productType, DateTime: now},
		},
		{
			name:        "invalid json",
			body:        strings.NewReader("{invalid json"),
			expectCall:  false,
			wantStatus:  http.StatusBadRequest,
			wantBodyOut: forms.ProductFormOut{},
		},
		{
			name:        "invalid product type",
			body:        toJSONBody(forms.ProductForm{PvzId: pvzId, Type: "рандомный тип"}),
			expectCall:  false,
			wantStatus:  http.StatusBadRequest,
			wantBodyOut: forms.ProductFormOut{},
		},
		{
			name:        "usecase error",
			body:        toJSONBody(forms.ProductForm{PvzId: pvzId, Type: productType}),
			expectCall:  true,
			input:       forms.ProductForm{PvzId: pvzId, Type: productType},
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

			if tt.expectCall {
				mockUseCase.EXPECT().
					AddProduct(gomock.Any(), tt.input).
					Return(tt.mockReturn, tt.mockError).
					Times(1)
			}

			req := httptest.NewRequest(http.MethodPost, "/products", tt.body)
			rec := httptest.NewRecorder()

			h.AddProduct(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			require.Equal(t, tt.wantStatus, res.StatusCode)

			if res.StatusCode == http.StatusCreated {
				var out forms.ProductFormOut
				err := json.NewDecoder(res.Body).Decode(&out)
				require.NoError(t, err)
				require.Equal(t, tt.wantBodyOut, out)
			}
		})
	}
}

func toJSONBody(v any) io.Reader {
	b, _ := json.Marshal(v)
	return bytes.NewReader(b)
}

func TestReceptionHandler_RemoveProduct(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		setVars    map[string]string
		mockExpect bool
		mockPvzID  uuid.UUID
		mockError  error
		wantStatus int
	}{
		{
			name:       "ok",
			url:        "/pvz/ok-id/delete_last_product",
			setVars:    map[string]string{"pvzId": uuid.New().String()},
			mockExpect: true,
			mockError:  nil,
			wantStatus: http.StatusOK,
		},
		{
			name:       "usecase error",
			url:        "/pvz/error-id/delete_last_product",
			setVars:    map[string]string{"pvzId": uuid.New().String()},
			mockExpect: true,
			mockError:  errors.New("some error"),
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing pvzId",
			url:        "/pvz//delete_last_product",
			setVars:    map[string]string{}, // No pvzId
			mockExpect: false,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid pvzId",
			url:        "/pvz/invalid-uuid/delete_last_product",
			setVars:    map[string]string{"pvzId": "invalid-uuid"},
			mockExpect: false,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := mocks.NewMockReceptionUseCase(ctrl)
			h := handlers.NewReceptionHandler(mockUseCase)

			var expectedPvzID uuid.UUID
			if tt.mockExpect {
				var err error
				expectedPvzID, err = uuid.Parse(tt.setVars["pvzId"])
				require.NoError(t, err)
				mockUseCase.EXPECT().
					RemoveProduct(gomock.Any(), expectedPvzID).
					Return(tt.mockError).
					Times(1)
			}

			req := httptest.NewRequest(http.MethodDelete, tt.url, nil)
			req = mux.SetURLVars(req, tt.setVars)
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
		url         string
		setVars     map[string]string
		mockExpect  bool
		mockReturn  models.Reception
		mockError   error
		wantStatus  int
		wantBodyOut forms.ReceptionFormOut
	}{
		{
			name:        "ok",
			url:         fmt.Sprintf("/pvz/%s/close_last_reception", pvzId.String()),
			setVars:     map[string]string{"pvzId": pvzId.String()},
			mockExpect:  true,
			mockReturn:  models.Reception{Id: receptionId, PvzId: pvzId, Status: models.InProgress},
			mockError:   nil,
			wantStatus:  http.StatusOK,
			wantBodyOut: forms.ReceptionFormOut{Id: receptionId, PvzId: pvzId, Status: string(models.InProgress)},
		},
		{
			name:        "usecase error",
			url:         fmt.Sprintf("/pvz/%s/close_last_reception", pvzId.String()),
			setVars:     map[string]string{"pvzId": pvzId.String()},
			mockExpect:  true,
			mockReturn:  models.Reception{},
			mockError:   errors.New("some error"),
			wantStatus:  http.StatusBadRequest,
			wantBodyOut: forms.ReceptionFormOut{},
		},
		{
			name:        "missing pvzId",
			url:         "/pvz//close_last_reception",
			setVars:     map[string]string{}, // no pvzId
			mockExpect:  false,
			wantStatus:  http.StatusBadRequest,
			wantBodyOut: forms.ReceptionFormOut{},
		},
		{
			name:        "invalid pvzId",
			url:         "/pvz/invalid-uuid/close_last_reception",
			setVars:     map[string]string{"pvzId": "invalid-uuid"},
			mockExpect:  false,
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

			if tt.mockExpect {
				parsedPvzID, err := uuid.Parse(tt.setVars["pvzId"])
				require.NoError(t, err)
				mockUseCase.EXPECT().
					CloseReception(gomock.Any(), parsedPvzID).
					Return(tt.mockReturn, tt.mockError).
					Times(1)
			}

			req := httptest.NewRequest(http.MethodPost, tt.url, nil)
			req = mux.SetURLVars(req, tt.setVars)
			rec := httptest.NewRecorder()

			h.CloseReception(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			require.Equal(t, tt.wantStatus, res.StatusCode)

			if res.StatusCode == http.StatusOK {
				var out forms.ReceptionFormOut
				err := json.NewDecoder(res.Body).Decode(&out)
				require.NoError(t, err)
				require.Equal(t, tt.wantBodyOut, out)
			}
		})
	}
}
