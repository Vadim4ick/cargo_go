package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func JSON(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := Response{
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(resp)
}

func ParseNumber(s string) (int, error) {
	return strconv.Atoi(s)
}
