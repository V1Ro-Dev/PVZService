package internal

import (
	"context"
	"errors"
	"fmt"
	"net/http"

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

	newAuthService := usecase.NewAuthService(newUserRepo)

	newAuthHandler := handlers.NewAuthHandler(newAuthService)

	defer newUserRepo.Close()

	r := mux.NewRouter()

	r.HandleFunc("/dummyLogin", newAuthHandler.DummyLogin)
	r.HandleFunc("/register", newAuthHandler.Register).Methods("POST")

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
