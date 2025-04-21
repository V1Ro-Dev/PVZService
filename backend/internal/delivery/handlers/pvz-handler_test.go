package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"pvz/config"
	"pvz/internal/delivery/forms"
	"pvz/internal/delivery/handlers"
	"pvz/internal/delivery/mocks"
	"pvz/internal/models"
)

func TestCreatePvz(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockPvzUseCase(ctrl)
	handler := handlers.NewPvzHandler(mockUC)

	tests := []struct {
		name         string
		input        interface{}
		mockError    error
		expectStatus int
		expectedBody interface{}
	}{
		{
			name: "valid create",
			input: forms.PvzForm{
				City: "Москва",
			},
			mockError:    nil,
			expectStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{"city": "Москва", "id": "00000000-0000-0000-0000-000000000000", "registrationDate": "0001-01-01T00:00:00Z"},
		},
		{
			name: "invalid city",
			input: forms.PvzForm{
				City: "InvalidCity",
			},
			mockError:    errors.New("invalid city"),
			expectStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{"message": "Invalid city"},
		},
		{
			name:         "invalid json",
			input:        `{{}`,
			mockError:    errors.New("invalid json"),
			expectStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{"message": "failed to parse json"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockError == nil {
				mockUC.EXPECT().CreatePvz(gomock.Any(), gomock.Eq(tt.input)).Return(models.Pvz{City: tt.input.(forms.PvzForm).City}, nil)
			}

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			handler.CreatePvz(rec, req)

			assert.Equal(t, tt.expectStatus, rec.Code)

			var responseBody map[string]interface{}
			if err := json.NewDecoder(rec.Body).Decode(&responseBody); err != nil {
				t.Fatalf("Error decoding response body: %s", err)
			}

			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}

func TestGetPvzInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockPvzUseCase(ctrl)
	handler := handlers.NewPvzHandler(mockUC)

	pvzId := uuid.New()
	receptionId := uuid.New()
	productId := uuid.New()

	pvz := models.Pvz{
		Id:               pvzId,
		RegistrationDate: time.Now().UTC().Truncate(time.Millisecond),
		City:             "Москва",
	}

	reception := models.Reception{
		Id:       receptionId,
		DateTime: time.Now().UTC().Truncate(time.Millisecond),
		PvzId:    pvzId,
		Status:   models.InProgress,
	}

	product := models.Product{
		Id:          productId,
		DateTime:    time.Now().UTC().Truncate(time.Millisecond),
		ProductType: "обувь",
		ReceptionId: receptionId,
	}

	receptionProducts := models.ReceptionProducts{
		Reception: reception,
		Products:  []models.Product{product},
	}

	tests := []struct {
		name         string
		queryParams  map[string]string
		mockResponse []models.PvzInfo
		mockError    error
		expectStatus int
		expectedBody interface{}
		expectedForm forms.GetPvzInfoForm
	}{
		{
			name: "valid request with page and limit",
			queryParams: map[string]string{
				"startDate": "2006-01-02T15:04:05.000Z",
				"endDate":   "2006-01-02T15:04:05.000Z",
				"page":      "1",
				"limit":     "10",
			},
			mockResponse: []models.PvzInfo{
				{
					Pvz: pvz,
					Receptions: []models.ReceptionProducts{
						receptionProducts,
					},
				},
			},
			mockError:    nil,
			expectStatus: http.StatusOK,
			expectedBody: []map[string]interface{}{
				{
					"pvz": map[string]interface{}{
						"city":             pvz.City,
						"id":               pvz.Id.String(),
						"registrationDate": pvz.RegistrationDate.Format(config.TimeStampLayout),
					},
					"receptions": []interface{}{
						map[string]interface{}{
							"reception": map[string]interface{}{
								"dateTime": reception.DateTime.Format(config.TimeStampLayout),
								"id":       reception.Id.String(),
								"pvzId":    reception.PvzId.String(),
								"status":   string(reception.Status),
							},
							"products": []interface{}{
								map[string]interface{}{
									"dateTime":    product.DateTime.Format(config.TimeStampLayout),
									"id":          product.Id.String(),
									"productType": product.ProductType,
									"receptionId": product.ReceptionId.String(),
								},
							},
						},
					},
				},
			},

			expectedForm: forms.GetPvzInfoForm{
				StartDate: mustParseTime("2006-01-02T15:04:05.000Z"),
				EndDate:   mustParseTime("2006-01-02T15:04:05.000Z"),
				Page:      1,
				Limit:     10,
			},
		},
		{
			name: "invalid date format",
			queryParams: map[string]string{
				"startDate": "invalid-date",
				"endDate":   "invalid-date",
			},
			mockResponse: []models.PvzInfo{},
			mockError:    nil,
			expectStatus: http.StatusOK,
			expectedBody: map[string]interface{}(nil),
			expectedForm: forms.GetPvzInfoForm{},
		},
		{
			name: "server error",
			queryParams: map[string]string{
				"startDate": "2006-01-02T15:04:05.000Z",
				"endDate":   "2006-01-02T15:04:05.000Z",
			},
			mockResponse: nil,
			mockError:    errors.New("unable to get pvz info"),
			expectStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{"message": "unable to get info"},
			expectedForm: forms.GetPvzInfoForm{
				StartDate: mustParseTime("2006-01-02T15:04:05.000Z"),
				EndDate:   mustParseTime("2006-01-02T15:04:05.000Z"),
				Page:      1,
				Limit:     10,
			},
		},
		{
			name:         "default values when params not provided",
			queryParams:  map[string]string{},
			mockResponse: []models.PvzInfo{},
			mockError:    nil,
			expectStatus: http.StatusOK,
			expectedBody: map[string]interface{}(nil),
			expectedForm: forms.GetPvzInfoForm{
				StartDate: time.Time{},
				EndDate:   time.Now(),
				Page:      1,
				Limit:     10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC.EXPECT().GetPvzInfo(
				gomock.Any(),
				gomock.Any(),
			).DoAndReturn(func(ctx context.Context, form forms.GetPvzInfoForm) ([]models.PvzInfo, error) {
				return tt.mockResponse, tt.mockError
			}).Times(1)

			req := httptest.NewRequest(http.MethodGet, "/pvz", nil)
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()

			rec := httptest.NewRecorder()
			handler.GetPvzInfo(rec, req)

			assert.Equal(t, tt.expectStatus, rec.Code)

			if tt.name == "valid request with page and limit" {
				var responseBody []map[string]interface{}
				if err := json.NewDecoder(rec.Body).Decode(&responseBody); err != nil {
					t.Fatalf("Error decoding response body: %s", err)
				}

				assert.Equal(t, tt.expectedBody, responseBody)
			} else {
				var responseBody map[string]interface{}
				if err := json.NewDecoder(rec.Body).Decode(&responseBody); err != nil {
					t.Fatalf("Error decoding response body: %s", err)
				}

				assert.Equal(t, tt.expectedBody, responseBody)
			}
		})
	}
}

func mustParseTime(timeStr string) time.Time {
	t, err := time.Parse(config.TimeStampLayout, timeStr)
	if err != nil {
		panic(err)
	}
	return t
}
