package server

import (
	http "net/http"

	db "api/internal/db"
	handlers "api/internal/handlers"

	mux "github.com/gorilla/mux"
)

// NewRouter constructs the Gorilla mux router, registers
// the application routes and applies middleware.
func NewRouter(s *db.Store) http.Handler {
    router := mux.NewRouter()

    router.HandleFunc("/api/go/actors", handlers.GetActors(s)).Methods("GET")
    router.HandleFunc("/api/go/actors", handlers.CreateActor(s)).Methods("POST")
    router.HandleFunc("/api/go/actors/{id}", handlers.GetActor(s)).Methods("GET")
    router.HandleFunc("/api/go/actors/{id}", handlers.UpdateActor(s)).Methods("PUT")
    router.HandleFunc("/api/go/actors/{id}", handlers.DeleteActor(s)).Methods("DELETE")

    // apply middleware chain
    return enableCORS(jsonContentTypeMiddleware(router))
}
