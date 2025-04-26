package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Укажите название новой папки (например: user)")
		return
	}

	name := strings.ToLower(os.Args[1])
	basePath := filepath.Join("internal", "delivery", "http", name)

	// Создание директории
	err := os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		fmt.Printf("Ошибка при создании папки: %v\n", err)
		return
	}

	// Путь к handler.go
	handlerFilePath := filepath.Join(basePath, "handler.go")

	// Проверка: если уже существует
	if _, err := os.Stat(handlerFilePath); err == nil {
		fmt.Println("handler.go уже существует.")
		return
	}

	// Содержимое handler.go
	handlerContent := fmt.Sprintf(`package %s

import (
	"net/http"
)

type Handler struct {}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/%s", h.ExampleHandler)
}

func (h *Handler) ExampleHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from %s handler!"))
}
`, name, name, name)

	// Запись в файл
	err = os.WriteFile(handlerFilePath, []byte(handlerContent), 0644)
	if err != nil {
		fmt.Printf("Ошибка при создании handler.go: %v\n", err)
		return
	}

	fmt.Printf("Успешно создано: %s\n", handlerFilePath)
}
