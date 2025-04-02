package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"sql-injection-go/internal/app"
	"sql-injection-go/internal/config"
	"sql-injection-go/internal/handlers"
	storage "sql-injection-go/internal/storage/postgres"
	srv "sql-injection-go/internal/app/server"

	"github.com/gin-gonic/gin"
)



func main() {
	config := config.MustLoad()
	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)

	store, err := storage.New(context.Background(), config.StorageConfig.DatabaseUrl)
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	injectionHandler := handlers.New(
		log, 
		store,
	)

	server, err := srv.New(
		gin.Default(),
		srv.WithHost("0.0.0.0"),
		srv.WithPort(8080))


	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't setup server: %v\n", err)
		os.Exit(1)
	}

	app, err := app.New(
		log,
		injectionHandler,
		server,
		*config)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable setup app: %v\n", err)
		os.Exit(1)
	}

	app.Run()

	// TODO: Graceful shutdown
}