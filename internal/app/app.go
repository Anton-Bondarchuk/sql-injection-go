package app

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"log/slog"
	srv "sql-injection-go/internal/app/server"
	"sql-injection-go/internal/config"
	"sql-injection-go/internal/lib/logger/handlers/slogpretty"

	"github.com/gin-gonic/gin"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type App struct {
	log *slog.Logger
	handler Handler
	server srv.HttpServer
}

// NOTE: does't extandable entity
type Handler interface {
	GetStudentInjection(c *gin.Context)
	GetStudentsSafe(c *gin.Context)
	RenderSearch(c *gin.Context)
}


func New(
	log *slog.Logger, 
	hander Handler,
	server *srv.HttpServer,
	cnf config.Config) (*App, error) {
	
	server.Srv.LoadHTMLGlob("templates/*")
	server.Srv.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusPermanentRedirect, "/search")
	})
	server.Srv.GET("/search", hander.RenderSearch)

	server.Srv.GET("/students", hander.GetStudentInjection)
	server.Srv.GET("/students_safe", hander.GetStudentsSafe)

	log = setupLogger(cnf.Env)

	return &App{
		log: log,
		handler: hander,
		server: *server,
	}, nil
}

func (a *App) Run() {
	portStr := strconv.Itoa(a.server.Port)
	addr := fmt.Sprintf("%s:%s", a.server.Host, portStr)

	a.server.Srv.Run(addr)
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