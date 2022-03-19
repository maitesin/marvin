package http

import (
	"net/http"

	"github.com/go-chi/chi"
)

func DefaultRouter() http.Handler {
	router := chi.NewRouter()

	router.Get("/hello", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hello, pinger!"))
	})

	return router
}
