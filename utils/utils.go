package utils

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
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
		logger.Info("HTTP response message",
			zap.String("message", message),
			zap.Int("status", status),
		)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("Failed to encode JSON response", zap.Error(err))
	}
}

func ParseNumber(s string) (int, error) {
	return strconv.Atoi(s)
}
