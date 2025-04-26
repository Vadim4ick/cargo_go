package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"test-project/config"
	router "test-project/internal/delivery"
	logger "test-project/pkg"

	_ "test-project/docs"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
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
	logger, err := logger.NewLogger("logs/app.log")

	if err != nil {
		log.Fatal("Ошибка инициализации логгера:", err)
	}

	defer logger.Sync()

	dsn := config.Envs.POSTGRES_URI

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		logger.Fatal("Ошибка подключения к базе данных", zap.Error(err))
	}
	defer pool.Close()

	// Проверка подключения к базе данных
	if err := pool.Ping(context.Background()); err != nil {
		logger.Fatal("Не удалось проверить подключение к базе данных", zap.Error(err))
	}
	logger.Info("Успешное подключение к базе данных")

	mux := router.Setup(pool, logger)

	logger.Info("Starting server on :8080")
	fmt.Println("Swagger UI available at http://localhost:8080/swagger/index.html")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		logger.Fatal("Ошибка запуска сервера", zap.Error(err))
	}
}
