package http

import (
	"net/http"

	"github.com/go-chi/chi"
)

func DefaultRouter() http.Handler {
	router := chi.NewRouter()

	return router
}
