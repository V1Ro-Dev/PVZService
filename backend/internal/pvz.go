package internal

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"pvz/internal/delivery/middleware"

	"github.com/gorilla/mux"

	"pvz/config"
	"pvz/internal/delivery/handlers"
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

	newAuthService := usecase.NewAuthService(newUserRepo)
	newPvzService := usecase.NewPvzService(newPvzRepo)

	newAuthHandler := handlers.NewAuthHandler(newAuthService)
	newPvzHandler := handlers.NewPvzHandler(newPvzService)

	defer newUserRepo.Close()

	r := mux.NewRouter()

	r.HandleFunc("/dummyLogin", newAuthHandler.DummyLogin).Methods("POST")
	r.HandleFunc("/register", newAuthHandler.Register).Methods("POST")
	r.HandleFunc("/login", newAuthHandler.Login).Methods("POST")

	protected := r.PathPrefix("/").Subrouter()
	protected.Use(middleware.RoleMiddleware("moderator"))
	protected.HandleFunc("/pvz", newPvzHandler.CreatePvz).Methods("POST")

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
