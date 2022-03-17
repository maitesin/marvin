package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func DefaultRouter(newRelicApp *newrelic.Application) http.Handler {
	router := chi.NewRouter()

	//http.HandleFunc(newrelic.WrapHandleFunc(newRelicApp, "/users", usersHandler))

	return router
}
