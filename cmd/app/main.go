package main

import (
	"context"
	"log"
	"os"

	"github.com/Puker228/subscriptions-app/internal/subscriptions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	echoSwagger "github.com/swaggo/echo-swagger/v2"

	_ "github.com/Puker228/subscriptions-app/docs"
)

// @title Subscriptions API
// @version 1.0
// @description API для управления подписками.
// @host localhost:8800
// @BasePath /api/v1
// @schemes http
func main() {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	if err := godotenv.Load(); err != nil {
		log.Println("env not loaded:", err)
	}

	ctx := context.Background()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("GOOSE_DBSTRING")
	}
	if dbURL == "" {
		log.Println("DATABASE_URL or GOOSE_DBSTRING is empty")
		os.Exit(1)
	}

	log.Println("db config ok")

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Println("db pool error:", err)
		os.Exit(1)
	}
	if err := pool.Ping(ctx); err != nil {
		log.Println("db ping error:", err)
		pool.Close()
		os.Exit(1)
	}
	defer pool.Close()
	log.Println("db connected")

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

	addr := ":8800"
	log.Println("starting server on", addr)
	if err := e.Start(addr); err != nil {
		log.Println("server error:", err)
	}
}
