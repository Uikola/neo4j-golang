package http

import (
	"github.com/Uikola/neo4j-golang/internal/server/http/host"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewServer(
	cloudEventHandler *host.Handler,
) http.Handler {
	router := chi.NewRouter()

	addRoutes(router, cloudEventHandler)

	var handler http.Handler = router

	return handler
}

func addRoutes(
	router *chi.Mux,
	cloudEventHandler *host.Handler,
) {
	router.Post("/save", cloudEventHandler.Save)
}
