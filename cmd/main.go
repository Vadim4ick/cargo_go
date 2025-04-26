package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"test-project/config"
	router "test-project/internal/delivery"

	_ "test-project/docs"

	"github.com/jackc/pgx/v5/pgxpool"
)

// @title Cargo Project API
// @version 1.0
// @description This is the API for the Cargo Project
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token for authenticated requests
func main() {

	dsn := config.Envs.POSTGRES_URI

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}
	defer pool.Close()

	// Проверка подключения к базе данных
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("Не удалось проверить подключение к базе данных:", err)
	}
	fmt.Println("Успешное подключение к базе данных!")

	mux := router.Setup(pool)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
