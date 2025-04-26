package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"test-project/internal/domain/user"
	"test-project/internal/middleware/auth"
	"test-project/internal/redis"
	"test-project/internal/usecase"
	"test-project/utils"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Handler struct {
	svc      usecase.AuthUsecase
	jwtSvc   *usecase.JwtUsecase
	logger   *zap.Logger
	redisSvc *redis.Client
}

func NewHandler(svc usecase.AuthUsecase, jwtSvc *usecase.JwtUsecase, logger *zap.Logger, redisSvc *redis.Client) *Handler {
	return &Handler{svc, jwtSvc, logger, redisSvc}
}

func RegisterCargoRoute(r *mux.Router, db *pgxpool.Pool, logger *zap.Logger, jwtService *usecase.JwtUsecase, redisService *redis.Client) {
	userRepo := user.NewPostgresUserRepo(db)

	svc := usecase.NewService(userRepo, jwtService, redisService)
	h := NewHandler(svc, jwtService, logger, redisService)

	r.HandleFunc("/register", h.register).Methods("POST")
	r.HandleFunc("/login", h.login).Methods("POST")
	r.HandleFunc("/logout", h.logout).Methods("POST")
	r.HandleFunc("/refresh", h.refresh).Methods("POST")

	r.Handle("/profile", auth.JwtMiddleware(h.svc, jwtService, logger, h.profile)).Methods("GET")
	r.Handle("/online", auth.JwtMiddleware(h.svc, jwtService, logger, h.onlineList)).Methods("GET")
}

// refresh handles token refresh
// @Summary Refresh access token
// @Description Refreshes access token using refresh token stored in cookie
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} auth.RefreshResponse "Token successfully refreshed"
// @Failure 401 {object} auth.ErrorResponse "Invalid or missing refresh token"
// @Failure 500 {object} auth.ErrorResponse "Internal server error"
// @Router /refresh [post]
func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		utils.JSON(w, http.StatusUnauthorized, "Refresh токен отсутствует", nil, h.logger)
		return
	}

	userID, err := h.jwtSvc.ValidateRefresh(cookie.Value)

	if err != nil {
		utils.JSON(w, http.StatusUnauthorized, "Невалидный refresh токен", nil, h.logger)
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/",   // такой же path как при установке
			HttpOnly: true,  // тоже такой же
			Secure:   false, // такой же как при установке
			SameSite: http.SameSiteLaxMode,
			MaxAge:   -1,
		})
		return
	}

	// Генерируем новые access и refresh токены
	newAccessToken, err := h.jwtSvc.GenerateAccess(userID)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, "Ошибка генерации access токена", nil, h.logger)
		return
	}
	newRefreshToken, err := h.jwtSvc.GenerateRefresh(userID)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, "Ошибка генерации refresh токена", nil, h.logger)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/",   // доступен во всём приложении
		HttpOnly: true,  // недоступен из JS
		Secure:   false, // true на HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((30 * 24 * time.Hour).Seconds()), // 30 дней
	})

	utils.JSON(w, http.StatusOK, "Токен успешно обновлён", map[string]string{
		"access_token": newAccessToken,
	}, h.logger)
}

// register handles user registration
// @Summary Register a new user
// @Description Registers a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body auth.RegisterRequest true "User credentials"
// @Success 201 {object} auth.RegisterResponse "User successfully registered"
// @Failure 400 {object} auth.ErrorResponse "Invalid input"
// @Router /register [post]
func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var req struct{ Email, Password string }
	json.NewDecoder(r.Body).Decode(&req)
	u, err := h.svc.Register(req.Email, req.Password)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, err.Error(), nil, h.logger)
		return
	}

	utils.JSON(w, http.StatusCreated, "Пользователь успешно зарегистрирован", u, h.logger)
}

// login handles user authentication
// @Summary User login
// @Description Authenticates a user and returns access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body auth.LoginRequest true "User credentials"
// @Success 200 {object} auth.LoginResponse "User successfully authenticated"
// @Failure 400 {object} auth.ErrorResponse "Invalid input"
// @Failure 401 {object} auth.ErrorResponse "Unauthorized"
// @Router /login [post]
func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSON(w, http.StatusBadRequest, "Некорректные данные", nil, h.logger)
		return
	}

	accessToken, refreshToken, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		utils.JSON(w, http.StatusUnauthorized, err.Error(), nil, h.logger)
		return
	}

	// w.Header().Set("Authorization", "Bearer "+token)

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true если HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((30 * 24 * time.Hour).Seconds()), // 30 дней
	})

	utils.JSON(w, http.StatusOK, "Пользователь успешно авторизован", map[string]string{
		"access_token": accessToken,
	}, h.logger)
}

// onlineList retrieves a list of online user IDs
// @Summary List online users
// @Description Retrieves a list of user IDs who are currently online
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} auth.OnlineListResponse "List of online user IDs"
// @Failure 401 {object} auth.ErrorResponse "Unauthorized"
// @Failure 500 {object} auth.ErrorResponse "Internal server error"
// @Router /online [get]
func (h *Handler) onlineList(w http.ResponseWriter, r *http.Request) {
	ids, err := h.svc.OnlineUsers(5 * time.Minute)

	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, fmt.Sprintf("error getting online users: %v", err), nil, h.logger)
		return
	}

	// убираем префикс "online:"
	for i, k := range ids {
		ids[i] = strings.TrimPrefix(k, "online:")
	}
	utils.JSON(w, http.StatusOK, "Список ID-шников онлайн пользователей", ids, h.logger)
}

// logout handles user logout
// @Summary User logout
// @Description Logs out a user and clears the refresh token cookie
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} auth.LogoutResponse "User successfully logged out"
// @Router /logout [post]
func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true если HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	utils.JSON(w, http.StatusOK, "Успешный выход из системы", nil, h.logger)
}

// profile handles user profile
// @Summary User profile
// @Description Retrieves user profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} auth.ProfileResponse "User profile information"
// @Failure 401 {object} auth.ErrorResponse "Unauthorized"
// @Failure 404 {object} auth.ErrorResponse "User not found"
// @Failure 500 {object} auth.ErrorResponse "Internal server error"
// @Router /profile [get]
func (h *Handler) profile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(string)

	if !ok {
		utils.JSON(w, http.StatusUnauthorized, "unauthorized", nil, h.logger)
		return
	}

	id, err := utils.ParseNumber(userID)

	if err != nil {
		utils.JSON(w, http.StatusBadRequest, "Некорректный id", nil, h.logger)
		return
	}

	u, err := h.svc.GetUser(id)

	if err != nil {
		utils.JSON(w, http.StatusNotFound, err.Error(), nil, h.logger)
		return
	}

	utils.JSON(w, http.StatusOK, "Профиль", u, h.logger)
}
