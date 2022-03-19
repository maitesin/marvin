package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func DefaultRouter() http.Handler {
	router := chi.NewRouter()

	router.Get("/hello", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write([]byte("Hello, pinger!"))
		if err != nil {
			fmt.Println(err)
		}
	})

	return router
}
