package http_response

import (
	"net/http"

	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/infrastructure/shared/app_error"
	"github.com/labstack/echo/v4"
)

func OK(c echo.Context, body interface{}) error {
	return c.JSON(http.StatusOK, body)
}

func Created(c echo.Context, body interface{}) error {
	return c.JSON(http.StatusCreated, body)
}

func NoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func BadRequest(c echo.Context, message string, err error) error {
	return c.JSON(http.StatusBadRequest, app_error.New(http.StatusBadRequest, message, err))
}

func UnprocessableEntity(c echo.Context, message string, err error) error {
	return c.JSON(http.StatusUnprocessableEntity, app_error.New(http.StatusUnprocessableEntity, message, err))
}

func NotFound(c echo.Context) error {
	return c.JSON(http.StatusNotFound, app_error.New(http.StatusNotFound, "resource not found"))
}

func Forbidden(c echo.Context, message string) error {
	return c.JSON(http.StatusForbidden, app_error.New(http.StatusForbidden, message))
}

func HandleErr(c echo.Context, err error) error {
	switch t := err.(type) {
	case *app_error.AppError:
		return c.JSON(t.Code, t)
	default:
		return c.JSON(http.StatusInternalServerError,
			app_error.New(http.StatusInternalServerError, "Internal server error", err))
	}
}
