package api

import (
	"github.com/cheviz/pitchdayBackend/models"
	"net/http"
)

//Check server health
func API(w http.ResponseWriter, r *http.Request) {
	hello := models.Response{Success: true, Debug: "Server is healthy", Message: "Server is healthy"}
	hello.Send(w, 200)
	return
}


