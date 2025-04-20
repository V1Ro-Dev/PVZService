package main

import (
	"log"

	"pvz/config"
	"pvz/internal"
)

func main() {
	cfg, err := config.Parse("")
	if err != nil {
		log.Fatalf("failed to load PVZ configuration: %v", err)
	}

	if err = internal.Run(cfg); err != nil {
		log.Fatalf("failed to start QuickFlow: %v", err)
	}
}
