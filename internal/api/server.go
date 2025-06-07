package api

import (
	"log"
	"net/http"
)

func StartServer() error {
	router := NewRouter()
	log.Println("API server running on :8080")
	return http.ListenAndServe(":8080", router)
}
