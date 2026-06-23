package main

import (
	"github.com/Puker228/subscriptions-app/internal/subscriptions"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	h := subscriptions.Handler{}

	apiV1 := e.Group("/api/v1")
	apiV1.POST("/sub", h.Create)
	apiV1.GET("/sub/:id", h.GetOneByID)
	apiV1.PUT("/sub", h.Update)
	apiV1.DELETE("/sub/:id", h.Delete)

	if err := e.Start(":8800"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
