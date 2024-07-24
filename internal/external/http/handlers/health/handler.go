package health

import (
	"context"
	"net/http"
	"time"

	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/domain/dto/health_dto"
	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/external/persistence"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	database *persistence.DbService
}

func NewHandler(database *persistence.DbService) *Handler {
	return &Handler{
		database: database,
	}
}

func (h *Handler) Handle(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	db, err := h.database.Instance.DB()
	if err != nil {
		status := health_dto.New(err)
		return c.JSON(http.StatusInternalServerError, status)
	}

	if err := db.PingContext(ctx); err != nil {
		status := health_dto.New(err)
		return c.JSON(http.StatusInternalServerError, status)
	}

	return c.JSON(http.StatusOK, health_dto.New(nil))
}
