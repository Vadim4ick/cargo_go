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
)

type Handler struct {
	svc    usecase.AuthUsecase
	jwtSvc *usecase.JwtUsecase
}

func NewHandler(svc usecase.AuthUsecase, jwtSvc *usecase.JwtUsecase) *Handler {
	return &Handler{svc, jwtSvc}

}

func RegisterCargoRoute(r *mux.Router, db *pgxpool.Pool) {

	userRepo := user.NewPostgresUserRepo(db)
	jwtService := usecase.NewJWTService("secretAccess", "secretRefresh", time.Minute*15, time.Hour*24*30)
	redisService := redis.New("localhost:6379", "")

	svc := usecase.NewService(userRepo, jwtService, redisService)
	h := NewHandler(svc, jwtService)

	r.HandleFunc("/register", h.register).Methods("POST")
	r.HandleFunc("/login", h.login).Methods("POST")
	r.HandleFunc("/refresh", h.refresh).Methods("POST")

	// защищённый роут: список онлайн пользователей
	r.Handle("/online", auth.JwtMiddleware(h.svc, jwtService, time.Minute*0, h.onlineList)).Methods("GET")
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		utils.JSON(w, http.StatusUnauthorized, "Refresh токен отсутствует", nil)
		return
	}

	userID, err := h.jwtSvc.ValidateRefresh(cookie.Value)

	if err != nil {
		utils.JSON(w, http.StatusUnauthorized, "Невалидный refresh токен", nil)
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
		utils.JSON(w, http.StatusInternalServerError, "Ошибка генерации access токена", nil)
		return
	}
	newRefreshToken, err := h.jwtSvc.GenerateRefresh(userID)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, "Ошибка генерации refresh токена", nil)
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
	})
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var req struct{ Email, Password string }
	json.NewDecoder(r.Body).Decode(&req)
	u, err := h.svc.Register(req.Email, req.Password)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.JSON(w, http.StatusCreated, "Пользователь успешно зарегистрирован", u)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSON(w, http.StatusBadRequest, "Некорректные данные", nil)
		return
	}

	accessToken, refreshToken, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		utils.JSON(w, http.StatusUnauthorized, err.Error(), nil)
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
	})
}

func (h *Handler) onlineList(w http.ResponseWriter, r *http.Request) {
	ids, err := h.svc.OnlineUsers(5 * time.Minute)

	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, fmt.Sprintf("error getting online users: %v", err), nil)
		return
	}

	// убираем префикс "online:"
	for i, k := range ids {
		ids[i] = strings.TrimPrefix(k, "online:")
	}
	utils.JSON(w, http.StatusOK, "Список ID-шников онлайн пользователей", ids)
}
