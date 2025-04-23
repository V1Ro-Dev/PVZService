package server

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	pvz "pvz/internal/grpc/pvz"
	"pvz/internal/repository"
	"pvz/internal/usecase"
	"pvz/pkg/logger"
)

func RunGrpcServer() error {
	ctx := context.Background()
	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("failed to listen: %v", err))
		return err
	}

	server := grpc.NewServer()

	newPvzRepo := repository.NewPostgresPvzRepository()
	defer newPvzRepo.Close()

	newPvzService := usecase.NewPvzService(newPvzRepo)

	pvz.RegisterPVZServiceServer(server, NewPvzManager(newPvzService))

	logger.Info(ctx, fmt.Sprintf("starting grpc server at %s", lis.Addr().String()))
	if err = server.Serve(lis); err != nil {
		logger.Error(ctx, fmt.Sprintf("failed to serve: %v", err))
		return err
	}

	return nil
}
