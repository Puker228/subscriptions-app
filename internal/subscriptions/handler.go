package subscriptions

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type Handler struct{}

func (h *Handler) Hello(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World!"})
}
