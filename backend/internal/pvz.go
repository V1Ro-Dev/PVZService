package internal

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"pvz/config"
	"pvz/internal/delivery/handlers"
	"pvz/internal/delivery/middleware"
	"pvz/internal/models"
	"pvz/internal/repository"
	"pvz/internal/usecase"
	"pvz/pkg/logger"
)

func Run(cfg *config.Config) error {
	if cfg == nil {
		return errors.New("config is nil")
	}

	ctx := context.Background()

	newUserRepo := repository.NewPostgresUserRepository()
	newPvzRepo := repository.NewPostgresPvzRepository()
	newReceptionRepo := repository.NewPostgresReceptionRepository()

	newAuthService := usecase.NewAuthService(newUserRepo)
	newPvzService := usecase.NewPvzService(newPvzRepo)
	newReceptionService := usecase.NewReceptionService(newReceptionRepo)

	newAuthHandler := handlers.NewAuthHandler(newAuthService)
	newPvzHandler := handlers.NewPvzHandler(newPvzService)
	newReceptionHandler := handlers.NewReceptionHandler(newReceptionService)

	defer newUserRepo.Close()

	r := mux.NewRouter()

	r.Use(middleware.RequestIDMiddleware)
	r.HandleFunc("/dummyLogin", newAuthHandler.DummyLogin).Methods("POST")
	r.HandleFunc("/register", newAuthHandler.Register).Methods("POST")
	r.HandleFunc("/login", newAuthHandler.Login).Methods("POST")

	// endpoints for moderators only
	protectedModer := r.PathPrefix("/").Subrouter()
	protectedModer.Use(middleware.RoleMiddleware(models.Moderator))
	protectedModer.HandleFunc("/pvz", newPvzHandler.CreatePvz).Methods("POST")

	// endpoints for moderators and employees
	protectedModerEmp := r.PathPrefix("/").Subrouter()
	protectedModerEmp.Use(middleware.RoleMiddleware(models.Moderator, models.Employee))
	protectedModerEmp.HandleFunc("/pvz", newPvzHandler.GetPvzInfo).Methods("GET")

	//endpoints for employees only
	protectedEmp := r.PathPrefix("/").Subrouter()
	protectedEmp.Use(middleware.RoleMiddleware(models.Employee))
	protectedEmp.HandleFunc("/receptions", newReceptionHandler.CreateReception).Methods("POST")
	protectedEmp.HandleFunc("/products", newReceptionHandler.AddProduct).Methods("POST")
	protectedEmp.HandleFunc("/pvz/{pvzId:[0-9a-fA-F-]{36}}/delete_last_product", newReceptionHandler.RemoveProduct).Methods("POST")
	protectedEmp.HandleFunc("/pvz/{pvzId:[0-9a-fA-F-]{36}}/close_last_reception", newReceptionHandler.CloseReception).Methods("POST")

	server := http.Server{
		Addr:         cfg.Addr,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	logger.Info(ctx, fmt.Sprintf("starting server at %s\n", cfg.Addr))
	err := server.ListenAndServe()
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("failed to start server: %v", err))
		return err
	}

	return nil
}
