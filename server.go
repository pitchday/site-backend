package main

import (
	"github.com/cheviz/pitchdayBackend/config"
	"github.com/cheviz/pitchdayBackend/controllers"
	"github.com/cheviz/pitchdayBackend/models"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"os"
)

var Logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	// Setup the global variables, settings, database, and cache stores
	err := models.Setup()
	if err != nil {
		// Fatal
		panic(err)
	}

	// Start the web servers
	Logger.Printf("API Server started at http://%s\n", config.Conf.ApiURL)
	http.ListenAndServe(config.Conf.ApiURL, handlers.CombinedLoggingHandler(os.Stdout, controllers.CreateRouter()))
}