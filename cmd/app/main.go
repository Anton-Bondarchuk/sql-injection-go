package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"sql-injection-go/internal/config"
	"sql-injection-go/internal/handlers"
	"sql-injection-go/internal/lib/logger/handlers/slogpretty"
	storage "sql-injection-go/internal/storage/postgres"

	"github.com/gin-gonic/gin"
)


const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	config := config.MustLoad()
	log := setupLogger(config.Env)

	store, err := storage.New(context.Background(), config.StorageConfig.DatabaseUrl)
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	injectionHandler := handlers.New(
		log, 
		store,
	)

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusPermanentRedirect, "/search")
	})
	router.GET("/search", injectionHandler.RenderSearch)

	router.GET("/students", injectionHandler.GetStudentInjection)
	router.GET("/students_safe", injectionHandler.GetStudentsSafe)

	// go func() {
	// 	router.Run("0.0.0.0:8080")
	// }()
	router.Run("0.0.0.0:8080")


	// TODO: Graceful shutdown 
}


func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}


func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}