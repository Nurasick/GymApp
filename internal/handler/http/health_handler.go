package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewHealthHandler() {
	// Health check is a simple handler, can be inline in routes
}

// Health godoc
// @Summary Health check
// @Description Check if API is running
// @ID health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "API is healthy"
// @Router /health [get]
func RegisterHealthRoutes(e *echo.Echo) {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
}
