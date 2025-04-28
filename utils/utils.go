package utils

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/form"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var formDecoder = form.NewDecoder()

// ParseFormData парсит form-data из запроса в переданную структуру
func ParseFormData(r *http.Request, dst interface{}) error {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return err
	}
	return formDecoder.Decode(dst, r.Form)
}

func JSON(w http.ResponseWriter, status int, message string, data interface{}, logger *zap.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := map[string]interface{}{
		"message": message,
		"data":    data,
	}

	// Логируем сообщение ответа
	if logger != nil {
		logger.WithOptions(zap.AddCallerSkip(1)).Info("HTTP response message",
			zap.String("message", message),
			zap.Int("status", status),
		)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil && logger != nil {
		logger.WithOptions(zap.AddCallerSkip(1)).Error("Failed to encode JSON response", zap.Error(err))
	}
}

func ParseNumber(s string) (int, error) {
	return strconv.Atoi(s)
}

const deleteExpiredInvitationsQuery = `
	DELETE FROM invitations
	WHERE "createdAt" + INTERVAL '5 minutes' < NOW()
`

func StartInvitationCleaner(db *pgxpool.Pool, logger *zap.Logger) {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		// Сразу очистить при старте
		cleanup(db, logger)

		for {
			select {
			case <-ticker.C:
				cleanup(db, logger)
			}
		}
	}()
}

func cleanup(db *pgxpool.Pool, logger *zap.Logger) {
	start := time.Now()
	_, err := db.Exec(context.Background(), deleteExpiredInvitationsQuery)
	duration := time.Since(start)

	if err != nil {
		logger.Error("Ошибка очистки устаревших приглашений", zap.Error(err))
	} else {
		logger.Info("Очистка устаревших приглашений выполнена", zap.Duration("duration", duration))
	}
}

func init() {
	// Для time.Time
	formDecoder.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
		s := strings.TrimSpace(vals[0])
		if s == "" {
			return time.Time{}, nil
		}
		return time.Parse(time.RFC3339, s)
	}, time.Time{})

	// Для *time.Time
	formDecoder.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
		s := strings.TrimSpace(vals[0])
		if s == "" {
			return (*time.Time)(nil), nil
		}
		t, err := time.Parse(time.RFC3339, s)
		return &t, err
	}, (*time.Time)(nil))
}
