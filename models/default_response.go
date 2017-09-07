package models

import (
	"encoding/json"
	"fmt"
	"github.com/pitchday/site-backend/config"
	"net/http"
)

// Response contains the attributes found in an API response
type Response struct {
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Debug   string      `json:"debug,omitempty"`
}

// sets the status code, and marshals the response, so it can be written to the given ResponseWriter.
func (r *Response) Send(w http.ResponseWriter, statusCode int) {
	resp, err := json.Marshal(r)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		Logger.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", config.Conf.AccessControlAllowOrigin)
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "%s", resp)
	return
}
