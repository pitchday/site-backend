package api

import (
	"github.com/pitchday/site-backend/models"
	"log"
	"net/http"
	"os"
)

var Logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

//Check server health
func API(w http.ResponseWriter, r *http.Request) {
	hello := models.Response{Success: true, Debug: "Server is healthy", Message: "Server is healthy"}
	hello.Send(w, 200)
	return
}
