package main

import (
	// Core functions to manage http
	"net/http"

	// Gorilla handlers? maybe helpers?
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Route struct
type Route struct {
	Name          string
	Method        string
	Pattern       string
	GzipMandatory bool
	HandlerFunc   http.HandlerFunc
}

// NewRouter add routes to somewhere?????
func NewRouter(routes []Route) http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var h http.Handler = route.HandlerFunc
		h = handlers.CompressHandler(h)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(h)
	}

	handler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"accept", "content-type", "origin", "x-custom-header", "authorization"}),
	)(router)
	return handler
}
