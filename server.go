package main

import (
	"github.com/gorilla/handlers"
	"github.com/pitchday/site-backend/config"
	"github.com/pitchday/site-backend/controllers"
	"github.com/pitchday/site-backend/models"
	"log"
	"net/http"
	"os"
)

var Logger = log.New(os.Stdout, "Info: ", log.Ldate|log.Ltime|log.Lshortfile)

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
