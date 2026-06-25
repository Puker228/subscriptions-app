package main

import (
	"context"
	"os"

	"github.com/Puker228/subscriptions-app/internal/subscriptions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	if err := godotenv.Load(); err != nil {
		e.Logger.Warn("failed to load .env file", "error", err)
	}

	ctx := context.Background()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("GOOSE_DBSTRING")
	}
	if dbURL == "" {
		e.Logger.Error("DATABASE_URL or GOOSE_DBSTRING is required")
		os.Exit(1)
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		e.Logger.Error("failed to create database pool", "error", err)
		os.Exit(1)
	}
	if err := pool.Ping(ctx); err != nil {
		e.Logger.Error("failed to connect to database", "error", err)
		pool.Close()
		os.Exit(1)
	}
	defer pool.Close()

	r := subscriptions.NewRepository(pool)
	s := subscriptions.NewService(r)
	h := subscriptions.NewHandler(s)

	apiV1 := e.Group("/api/v1")
	apiV1.GET("/sub", h.List)
	apiV1.GET("/sub/sum", h.Sum)
	apiV1.POST("/sub", h.Create)
	apiV1.GET("/sub/:id", h.GetOneByID)
	apiV1.PUT("/sub", h.Update)
	apiV1.DELETE("/sub/:id", h.Delete)

	if err := e.Start(":8800"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
