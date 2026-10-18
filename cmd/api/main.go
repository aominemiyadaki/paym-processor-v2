package main

import (
	"context"
	"log"

	"paym-processor/internal/app"
)

func main() {
	ctx := context.Background()

	application, err := app.New(ctx)
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}
	defer application.Close()

	if err := application.Run(ctx); err != nil {
		log.Fatalf("application runtime error: %v", err)
	}
}
