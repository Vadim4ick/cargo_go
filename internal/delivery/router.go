package router

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"test-project/config"
	"test-project/internal/delivery/http/auth"
	"test-project/internal/delivery/http/cargo"
	"test-project/internal/delivery/http/invitation"
	"test-project/internal/delivery/http/truck"
	"test-project/internal/delivery/http/user"
	authDomain "test-project/internal/domain/auth"
	userDomain "test-project/internal/domain/user"
	"test-project/internal/middleware"
	"test-project/internal/redis"
	"test-project/internal/usecase"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	httpSwagger "github.com/swaggo/http-swagger"
)

func Setup(pool *pgxpool.Pool, logger *zap.Logger, jwtService *usecase.JwtUsecase, redisService *redis.Client) *mux.Router {
	r := mux.NewRouter()

	// 2) Статика: всё из ./uploads по /api/v1/uploads/*
	_, err := os.ReadDir("./uploads")
	if err != nil {
		log.Fatalf("не удалось прочитать папку uploads: %v", err)
	}

	subrouter := r.PathPrefix("/api/v1").Subrouter()

	// Настройка статического файлового сервера для папки uploads
	// Маршрут /api/v1/uploads/ будет обслуживать файлы из ./uploads
	fileServer := http.FileServer(http.Dir("./uploads"))
	subrouter.PathPrefix("/uploads/").Handler(http.StripPrefix("/api/v1/uploads/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем заголовок Content-Disposition для скачивания
		filename := filepath.Base(r.URL.Path)
		w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)

		// Устанавливаем CORS-заголовки
		w.Header().Set("Access-Control-Allow-Origin", config.Envs.FRONT_URI) // Укажите ваш фронтенд URL, например, http://localhost:3000
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Обслуживаем файл
		fileServer.ServeHTTP(w, r)
	})))

	swaggerUsername := config.Envs.SWAGGER_LOGIN
	swaggerPassword := config.Envs.SWAGGER_PASS
	swaggerHandler := middleware.AuthSwagger(httpSwagger.WrapHandler, swaggerUsername, swaggerPassword)

	subrouter.PathPrefix("/swagger/").Handler(swaggerHandler)

	userRepo := userDomain.NewPostgresUserRepo(pool)
	authSvc := usecase.NewService(userRepo, jwtService, redisService)

	deps := &authDomain.Deps{
		Logger:      logger,
		JwtService:  jwtService,
		AuthService: authSvc,
		Redis:       redisService,
		DB:          pool,
	}

	user.RegisterUserRoutes(subrouter, deps)
	truck.RegisterUserRoutes(subrouter, deps)
	cargo.RegisterCargoRoute(subrouter, deps)
	auth.RegisterCargoRoute(subrouter, deps)
	invitation.RegisterInvitationRoutes(subrouter, deps)

	return subrouter
}
