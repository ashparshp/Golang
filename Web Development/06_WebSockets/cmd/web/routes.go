package main

import (
	"net/http"

	"github.com/ashparshp/ws/internal/handlers"
	"github.com/bmizerany/pat"
)

func routes() http.Handler {
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(handlers.Home))
	
	return mux
}
