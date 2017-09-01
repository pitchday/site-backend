package controllers

import (
	"github.com/cheviz/pitchdayBackend/controllers/api"
	"github.com/cheviz/pitchdayBackend/middleware"
	"github.com/gorilla/mux"
	"github.com/justinas/nosurf"
	"log"
	"net/http"
	"os"
)

var Logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

func CreateRouter() http.Handler {
	router := mux.NewRouter()

	router = router.StrictSlash(true)
	router.HandleFunc("/", Use(api.API)).Methods("GET")
	router.HandleFunc("/newsletter", Use(api.Newsletter_Subscribe)).Methods("POST")

	// Setup CSRF Protection
	csrfHandler := nosurf.New(router)

	// Exempt API routes and Static files
	csrfHandler.ExemptGlob("/*")

	return Use(csrfHandler.ServeHTTP, middleware.GetContext)
}

func Use(handler http.HandlerFunc, mid ...func(http.Handler) http.HandlerFunc) http.HandlerFunc {
	for _, m := range mid {
		handler = m(handler)
	}
	return handler
}
