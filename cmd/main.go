package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"test-project/config"
	router "test-project/internal/delivery"

	"github.com/jackc/pgx/v5/pgxpool"
)

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
