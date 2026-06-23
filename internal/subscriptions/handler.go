package subscriptions

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	s *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) Create(c *echo.Context) error {
	var sub Subscription

	if err := c.Bind(&sub); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	created, err := h.s.Create(c.Request().Context(), sub)
	if err != nil {
		return h.handleError(c, err, "failed to create subscription")
	}

	return c.JSON(http.StatusCreated, created)
}

func (h *Handler) GetOneByID(c *echo.Context) error {
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid id",
		})
	}

	sub, err := h.s.GetOneByID(c.Request().Context(), ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "subscription not found",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to get subscription",
		})
	}

	return c.JSON(http.StatusOK, sub)
}

func (h *Handler) Update(c *echo.Context) error {
	var sub Subscription

	if err := c.Bind(&sub); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if err := h.s.Update(c.Request().Context(), sub); err != nil {
		return h.handleError(c, err, "failed to update subscription")
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) Delete(c *echo.Context) error {
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid id",
		})
	}

	if err := h.s.Delete(c.Request().Context(), ID); err != nil {
		return h.handleError(c, err, "failed to delete subscription")
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) handleError(c *echo.Context, err error, message string) error {
	if errors.Is(err, ErrEmptyServiceName) || errors.Is(err, ErrInvalidPrice) || errors.Is(err, ErrInvalidUserID) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "subscription not found",
		})
	}

	return c.JSON(http.StatusInternalServerError, map[string]string{
		"error": message,
	})
}
