package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"test-project/config"
	router "test-project/internal/delivery"
	"test-project/internal/redis"
	"test-project/internal/usecase"
	logger "test-project/pkg"
	"test-project/utils"

	_ "test-project/docs"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
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

	if err := pool.Ping(context.Background()); err != nil {
		logger.Fatal("Не удалось проверить подключение к базе данных", zap.Error(err))
	}
	logger.Info("Успешное подключение к базе данных")

	jwtService := usecase.NewJWTService(
		config.Envs.JWTSecretAccess,
		config.Envs.JWTSecretRefresh,
		config.Envs.JWTAccessTTL,
		config.Envs.JWTRefreshTTL,
	)

	redisService := redis.New(config.Envs.RedisAddr, config.Envs.RedisPassword)

	mux := router.Setup(pool, logger, jwtService, redisService)

	// Настройка CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{config.Envs.FRONT_URI}, // Укажите разрешённые домены
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true, // Разрешить отправку куки и заголовков авторизации
		MaxAge:           300,  // Кэширование CORS-запросов (в секундах)
	})

	handler := corsHandler.Handler(mux)

	utils.StartInvitationCleaner(pool, logger)
	logger.Info("Starting server on :8080")
	fmt.Printf("Swagger UI available at %s/api/v1/swagger/index.html", config.Envs.API_URI)
	if err := http.ListenAndServe(":8080", handler); err != nil {
		logger.Fatal("Ошибка запуска сервера", zap.Error(err))
	}
}
