package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"sql-injection-go/internal/domain/models"
	"strconv"
	"html/template"

	"github.com/gin-gonic/gin"
)


type IjectionHandler struct {
	log *slog.Logger
	injectionProvider InjectionProvider
}

type InjectionProvider interface {
	GetStudentsSafe(ctx context.Context, id int) ([]models.Student, error)
	GetStudentInjection(ctx context.Context, id string) ([]models.Student, error)
}

func New(log *slog.Logger, injectionProvider InjectionProvider) *IjectionHandler {
	return &IjectionHandler{
		log: log,
		injectionProvider: injectionProvider,
	}
}


func (h* IjectionHandler) GetStudentsSafe(c *gin.Context) {
	const op = "get_student.safe"

	query := c.Query("query")
	id, err := strconv.Atoi(query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Invalid student number", "number": query})
	}

	log := h.log.With(
		slog.String("op", op), 
		slog.String("id", query),
	)

	log.Info("Attempting to retrieve student records")

	students, err := h.injectionProvider.GetStudentsSafe(c.Request.Context(), id)
	if (err != nil) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, students)
}

func (h* IjectionHandler) GetStudentInjection(c *gin.Context) {
	const op = "get_student.injection"

	query := c.Query("query")

	log := h.log.With(
		slog.String("op", op), 
		slog.String("id", query),
	)

	log.Info("Attempting to retrieve student records")

	students, err := h.injectionProvider.GetStudentInjection(c.Request.Context(), query)
	if (err != nil) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, students)
}


func (h *IjectionHandler) RenderSearch(c *gin.Context) {
	tmpl, err := template.ParseFiles("templates/search.html")
	
	// TODO: add slog
	if err != nil {
		// log.Printf("Ошибка парсинга шаблона: %v", err)
		c.String(500, "Ошибка шаблона")
		return
	}

	err = tmpl.Execute(c.Writer, nil)
	if err != nil {
		// log.Printf("Ошибка выполнения шаблона: %v", err)
		c.String(500, "Ошибка шаблона")
		return
	}
}