package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter creates a mux router using the routes variable
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler = route.HandlerFunc
		handler = LoggerHTTPHandler(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	return router
}
