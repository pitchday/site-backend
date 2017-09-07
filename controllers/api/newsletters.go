package api

import (
	"encoding/json"
	"github.com/pitchday/site-backend/models"
	"net/http"
)

//Add given email to mailing list
func Newsletter_Subscribe(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var ns models.NewsletterSubscription

	err := decoder.Decode(&ns)
	if err != nil {
		resp := models.Response{
			Success: false,
			Debug:   "Cannot parse request",
			Message: "Unable to subscribe",
		}
		resp.Send(w, 400)
		return
	}

	err = ns.Validate()
	if err != nil {
		resp := models.Response{
			Success: false,
			Debug:   "Invalid email",
			Message: "Unable to subscribe",
		}
		resp.Send(w, 400)
		return
	}

	err, isDuplicated := ns.Create()
	if err != nil {
		if isDuplicated {
			resp := models.Response{
				Success: false,
				Debug:   "Email already registered",
				Message: "Unable to subscribe",
			}
			resp.Send(w, 409)
			return
		}
		resp := models.Response{
			Success: false,
			Debug:   "Internal server error",
			Message: "Unable to subscribe",
		}
		resp.Send(w, 500)
		return
	}

	resp := models.Response{
		Success: true,
		Message: "Subscribed successfully",
	}
	resp.Send(w, 200)
	return
}
