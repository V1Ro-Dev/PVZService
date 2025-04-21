//go:build integration
// +build integration

package integration_tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	"pvz/internal/delivery/forms"
	"pvz/internal/delivery/handlers"
	"pvz/internal/delivery/middleware"
	"pvz/internal/models"
	"pvz/internal/repository/fakes"
	"pvz/internal/usecase"
)

func SetupTest() *mux.Router {
	fakePvzRepo := fakes.NewPostgresPvzRepository()
	fakeReceptionRepo := fakes.NewFakeReceptionRepository()
	fakeUserRepo := fakes.NewFakeUserRepository()

	newAuthService := usecase.NewAuthService(fakeUserRepo)
	newPvzService := usecase.NewPvzService(fakePvzRepo)
	newReceptionService := usecase.NewReceptionService(fakeReceptionRepo)

	newAuthHandler := handlers.NewAuthHandler(newAuthService)
	newPvzHandler := handlers.NewPvzHandler(newPvzService)
	newReceptionHandler := handlers.NewReceptionHandler(newReceptionService)

	r := mux.NewRouter()

	r.Use(middleware.RequestIDMiddleware)
	r.HandleFunc("/dummyLogin", newAuthHandler.DummyLogin).Methods("POST")

	// endpoints for moderators only
	protectedModer := r.PathPrefix("/").Subrouter()
	protectedModer.Use(middleware.RoleMiddleware(models.Moderator))
	protectedModer.HandleFunc("/pvz", newPvzHandler.CreatePvz).Methods("POST")

	//endpoints for employees only
	protectedEmp := r.PathPrefix("/").Subrouter()
	protectedEmp.Use(middleware.RoleMiddleware(models.Employee))
	protectedEmp.HandleFunc("/receptions", newReceptionHandler.CreateReception).Methods("POST")
	protectedEmp.HandleFunc("/products", newReceptionHandler.AddProduct).Methods("POST")
	protectedEmp.HandleFunc("/pvz/{pvzId:[0-9a-fA-F-]{36}}/close_last_reception", newReceptionHandler.CloseReception).Methods("POST")

	return r
}

func TestBasicFlow(t *testing.T) {
	router := SetupTest()
	server := httptest.NewServer(router)
	defer server.Close()

	client := server.Client()

	// Шаг 1: DummyLogin для получения куки с ролью
	reqBody := `{"role": "moderator"}`
	resp, err := client.Post(server.URL+"/dummyLogin", "application/json", strings.NewReader(reqBody))

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var modToken string
	if err = json.NewDecoder(resp.Body).Decode(&modToken); err != nil {
		t.Fatal(err)
	}
	require.NotEmpty(t, modToken)

	// Шаг 2: Создание ПВЗ
	pvzPayload := `{
	  "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
	  "registrationDate": "2025-04-23T01:52:52.102Z",
	  "city": "Москва"
	}`
	req, _ := http.NewRequest("POST", server.URL+"/pvz", strings.NewReader(pvzPayload))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", modToken))
	resp, err = client.Do(req)

	var response forms.PvzForm
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	require.NotEmpty(t, response)

	// Шаг 3: DummyLogin как employee
	reqBody = `{"role": "employee"}`
	resp, err = client.Post(server.URL+"/dummyLogin", "application/json", strings.NewReader(reqBody))

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var empToken string
	if err = json.NewDecoder(resp.Body).Decode(&empToken); err != nil {
		t.Fatal(err)
	}
	require.NotEmpty(t, empToken)

	// Шаг 4: Создание приёмки
	receptionPayload := fmt.Sprintf(`{"pvzId": "%s"}`, response.Id)

	req, _ = http.NewRequest("POST", server.URL+"/receptions", strings.NewReader(receptionPayload))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", empToken))
	resp, err = client.Do(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var reception forms.ReceptionFormOut
	if err = json.NewDecoder(resp.Body).Decode(&reception); err != nil {
		t.Fatal(err)
	}
	require.NotEmpty(t, reception)

	// Шаг 5: Добавление 50 товаров
	prType := "обувь"
	for i := 0; i < 50; i++ {
		productPayload := fmt.Sprintf(`{
			"type": "%s",
			"pvzId": "%s"
		}`, prType, response.Id.String())
		req, err = http.NewRequest("POST", server.URL+"/products", strings.NewReader(productPayload))
		require.NoError(t, err)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", empToken))
		resp, err = client.Do(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var product forms.ProductFormOut
		if err = json.NewDecoder(resp.Body).Decode(&product); err != nil {
			t.Fatal(err)
		}
		require.NotEmpty(t, product)
	}

	// Шаг 6: Закрытие приёмки
	closeURL := fmt.Sprintf("%s/pvz/%s/close_last_reception", server.URL, response.Id)
	req, _ = http.NewRequest("POST", closeURL, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", empToken))
	resp, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var closedReception forms.ReceptionFormOut
	if err = json.NewDecoder(resp.Body).Decode(&closedReception); err != nil {
		t.Fatal(err)
	}
	require.NotEmpty(t, closedReception)
}
