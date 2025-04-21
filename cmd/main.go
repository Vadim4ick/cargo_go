package main

import (
	"fmt"
	"log"
	"net/http"

	router "test-project/internal/delivery"
)

func main() {
	mux := router.Setup()

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
