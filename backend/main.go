package main

import (
	"log"
	"sync"

	"pvz/config"
	"pvz/internal"
	grpc "pvz/internal/grpc/server"
)

func main() {
	cfg, err := config.Parse("")
	if err != nil {
		log.Fatalf("failed to load PVZ configuration: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err = internal.Run(cfg); err != nil {
			log.Fatalf("failed to start PVZ HTTP Service: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err = grpc.RunGrpcServer(); err != nil {
			log.Fatalf("failed to start PVZ gRPC Service: %v", err)
		}
	}()
	
	wg.Wait()
}
