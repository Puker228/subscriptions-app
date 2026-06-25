package subscriptions

import (
	"errors"
	"net/http"
	"strconv"

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

// List godoc
// @Summary Получить список подписок
// @Description Возвращает список подписок с фильтрацией, сортировкой и пагинацией.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param service_name query string false "Фильтр по названию сервиса"
// @Param q query string false "Алиас для service_name"
// @Param user_id query string false "UUID пользователя"
// @Param sort query string false "Поле сортировки" Enums(service_name, price, start_date, end_date)
// @Param order query string false "Направление сортировки" Enums(asc, desc)
// @Param page query int false "Номер страницы" minimum(1) default(1)
// @Success 200 {object} ListResult
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sub [get]
func (h *Handler) List(c *echo.Context) error {
	ctx := c.Request().Context()
	p := ListParams{
		ServiceName: c.QueryParam("service_name"),
		Sort:        c.QueryParam("sort"),
		Order:       c.QueryParam("order"),
		PageSize:    10,
	}
	if p.ServiceName == "" {
		p.ServiceName = c.QueryParam("q")
	}
	if page, err := strconv.Atoi(c.QueryParam("page")); err == nil && page > 0 {
		p.Page = page
	}
	if userID := c.QueryParam("user_id"); userID != "" {
		id, err := uuid.Parse(userID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid user_id",
			})
		}
		p.UserID = id
	}

	result, err := h.s.List(ctx, p)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to list subscriptions",
		})
	}

	return c.JSON(http.StatusOK, result)
}

// Sum godoc
// @Summary Посчитать сумму подписок
// @Description Возвращает суммарную стоимость подписок за период с optional-фильтрами по пользователю и сервису.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param start_date query string true "Начало периода в формате MM-YYYY"
// @Param end_date query string true "Конец периода в формате MM-YYYY"
// @Param service_name query string false "Фильтр по названию сервиса"
// @Param q query string false "Алиас для service_name"
// @Param user_id query string false "UUID пользователя"
// @Success 200 {object} SumResult
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sub/sum [get]
func (h *Handler) Sum(c *echo.Context) error {
	ctx := c.Request().Context()
	p := SumParams{
		StartDate:   c.QueryParam("start_date"),
		EndDate:     c.QueryParam("end_date"),
		ServiceName: c.QueryParam("service_name"),
	}
	if p.ServiceName == "" {
		p.ServiceName = c.QueryParam("q")
	}
	if userID := c.QueryParam("user_id"); userID != "" {
		id, err := uuid.Parse(userID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid user_id",
			})
		}
		p.UserID = id
	}

	result, err := h.s.Sum(ctx, p)
	if err != nil {
		return h.handleError(c, err, "failed to sum subscriptions")
	}

	return c.JSON(http.StatusOK, result)
}

// Create godoc
// @Summary Создать подписку
// @Description Создаёт новую подписку. Если id не передан, он генерируется автоматически.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body Subscription true "Данные подписки"
// @Success 201 {object} Subscription
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sub [post]
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

// GetOneByID godoc
// @Summary Получить подписку по ID
// @Description Возвращает одну подписку по UUID.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "UUID подписки"
// @Success 200 {object} Subscription
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sub/{id} [get]
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

// Update godoc
// @Summary Обновить подписку
// @Description Обновляет подписку по id из тела запроса.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body Subscription true "Данные подписки"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sub [put]
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

// Delete godoc
// @Summary Удалить подписку
// @Description Удаляет подписку по UUID.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "UUID подписки"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sub/{id} [delete]
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
	if errors.Is(err, ErrEmptyServiceName) || errors.Is(err, ErrInvalidPrice) || errors.Is(err, ErrInvalidUserID) || errors.Is(err, ErrEmptyStartDate) || errors.Is(err, ErrEmptyEndDate) || errors.Is(err, ErrInvalidDate) || errors.Is(err, ErrInvalidPeriod) {
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
