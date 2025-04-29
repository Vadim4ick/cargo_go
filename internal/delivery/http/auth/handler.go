package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"test-project/internal/domain/auth"
	"test-project/internal/middleware"
	"test-project/utils"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	deps auth.Deps
}

func NewHandler(deps *auth.Deps) *Handler {
	return &Handler{
		deps: *deps,
	}
}

func RegisterCargoRoute(r *mux.Router, deps *auth.Deps) {
	h := NewHandler(deps)

	r.HandleFunc("/auth/register", h.register).Methods(http.MethodPost)
	r.HandleFunc("/auth/login", h.login).Methods(http.MethodPost)
	r.HandleFunc("/auth/logout", h.logout).Methods(http.MethodPost)
	r.HandleFunc("/auth/refresh", h.refresh).Methods(http.MethodPost)

	r.Handle("/validate-token", middleware.JwtMiddleware(deps, h.validateToken)).Methods(http.MethodPost)
	r.Handle("/profile", middleware.JwtMiddleware(deps, h.profile)).Methods(http.MethodGet)
	r.Handle("/auth/online", middleware.JwtMiddleware(deps, h.onlineList)).Methods(http.MethodGet)
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
// @Router /auth/refresh [post]
func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		utils.JSON(w, http.StatusUnauthorized, "Refresh токен отсутствует", nil, h.deps.Logger)
		return
	}

	userID, err := h.deps.JwtService.ValidateRefresh(cookie.Value)

	if err != nil {
		utils.JSON(w, http.StatusUnauthorized, "Невалидный refresh токен", nil, h.deps.Logger)
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
	newAccessToken, err := h.deps.JwtService.GenerateAccess(userID)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, "Ошибка генерации access токена", nil, h.deps.Logger)
		return
	}
	newRefreshToken, err := h.deps.JwtService.GenerateRefresh(userID)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, "Ошибка генерации refresh токена", nil, h.deps.Logger)
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
	}, h.deps.Logger)
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
// @Router /auth/register [post]
func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var req struct{ Email, InviteToken, Password string }
	json.NewDecoder(r.Body).Decode(&req)

	res, err := h.deps.JwtService.ValidateInvite(req.InviteToken)

	if err != nil {
		utils.JSON(w, http.StatusBadRequest, err.Error(), nil, h.deps.Logger)
		return
	}

	if res != req.Email {
		utils.JSON(w, http.StatusBadRequest, "Указана неверная почта которая была указана при приглашении", nil, h.deps.Logger)
		return
	}

	u, err := h.deps.AuthService.Register(req.Email, req.Password)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, err.Error(), nil, h.deps.Logger)
		return
	}

	utils.JSON(w, http.StatusCreated, "Пользователь успешно зарегистрирован", u, h.deps.Logger)
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
// @Router /auth/login [post]
func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSON(w, http.StatusBadRequest, "Некорректные данные", nil, h.deps.Logger)
		return
	}

	accessToken, refreshToken, err := h.deps.AuthService.Login(req.Email, req.Password)
	if err != nil {
		utils.JSON(w, http.StatusUnauthorized, err.Error(), nil, h.deps.Logger)
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
	}, h.deps.Logger)
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
// @Router /auth/online [get]
func (h *Handler) onlineList(w http.ResponseWriter, r *http.Request) {
	ids, err := h.deps.AuthService.OnlineUsers(5 * time.Minute)

	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, fmt.Sprintf("error getting online users: %v", err), nil, h.deps.Logger)
		return
	}

	// убираем префикс "online:"
	for i, k := range ids {
		ids[i] = strings.TrimPrefix(k, "online:")
	}
	utils.JSON(w, http.StatusOK, "Список ID-шников онлайн пользователей", ids, h.deps.Logger)
}

// logout handles user logout
// @Summary User logout
// @Description Logs out a user and clears the refresh token cookie
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} auth.LogoutResponse "User successfully logged out"
// @Router /auth/logout [post]
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

	utils.JSON(w, http.StatusOK, "Успешный выход из системы", nil, h.deps.Logger)
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
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)

	if !ok {
		utils.JSON(w, http.StatusUnauthorized, "unauthorized", nil, h.deps.Logger)
		return
	}

	u, err := h.deps.AuthService.GetUser(userID)

	if err != nil {
		utils.JSON(w, http.StatusNotFound, err.Error(), nil, h.deps.Logger)
		return
	}

	utils.JSON(w, http.StatusOK, "Профиль", u, h.deps.Logger)
}

// validateToken handles token validation
// @Summary Validate token
// @Description Validates an access token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} auth.ErrorResponse "Token is valid"
// @Failure 401 {object} auth.ErrorResponse "Unauthorized"
// @Failure 500 {object} auth.ErrorResponse "Internal server error"
// @Router /validate-token [post]
func (h *Handler) validateToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		utils.JSON(w, http.StatusUnauthorized, "missing token", map[string]bool{"isValid": false}, h.deps.Logger)
		return
	}

	_, err := h.deps.JwtService.ValidateAccess(parts[1])

	if err != nil {
		utils.JSON(w, http.StatusUnauthorized, err.Error(), map[string]bool{"isValid": false}, h.deps.Logger)
		return
	}

	utils.JSON(w, http.StatusOK, "token is valid", map[string]bool{"isValid": true}, h.deps.Logger)
}
