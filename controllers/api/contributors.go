package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pitchday/site-backend/models"
	"net/http"
	"strconv"
)

func Get_Contributors(w http.ResponseWriter, r *http.Request) {
	contributors := models.ContributorList{}

	err := contributors.Get()
	if err != nil {
		resp := models.Response{
			Success: false,
			Debug:   "If this error persists please submit a bug report",
			Message: "Unable to retrieve contributors",
		}
		resp.Send(w, 500)
		return
	}

	resp := models.Response{
		Success: true,
		Message: "Retrieved successfully",
		Data:    contributors,
	}
	resp.Send(w, 200)
	return
}

//Add given email to mailing list
func Add_Contributor(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var c models.Contributor

	err := decoder.Decode(&c)
	if err != nil {
		resp := models.Response{
			Success: false,
			Debug:   "Cannot parse request",
			Message: "Unable to add contributor",
		}
		resp.Send(w, 400)
		return
	}

	err, isDuplicated := c.Create()
	if err != nil {
		if isDuplicated {
			resp := models.Response{
				Success: false,
				Debug:   "Email already registered",
				Message: "Unable to add contributor",
			}
			resp.Send(w, 409)
			return
		}
		resp := models.Response{
			Success: false,
			Debug:   "If this error persists please submit a bug report",
			Message: "Unable to add contributor",
		}
		resp.Send(w, 500)
		return
	}

	resp := models.Response{
		Success: true,
		Message: "Added successfully",
	}
	resp.Send(w, 200)
	return
}

func Remove_Contributor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["contributorId"]
	Id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		failParams := models.Response{
			Success: false,
			Debug:   fmt.Sprintf("Invalid id, %s", idStr),
			Message: "Unable to remove contributor.",
		}
		failParams.Send(w, 400)
		return
	}

	c := models.Contributor{
		Id: Id,
	}

	err, isNotFound := c.Get()
	if err != nil {
		if isNotFound {
			failParams := models.Response{
				Success: false,
				Debug:   fmt.Sprintf("Unable to find contributor %s.", idStr),
				Message: "Unable to remove contributor.",
			}
			failParams.Send(w, 404)
			return
		}
		failParams := models.Response{
			Success: false,
			Debug:   "If this error persists please submit a bug report",
			Message: "Unable to remove contributor.",
		}
		failParams.Send(w, 500)
		return
	}

	err = c.Delete()
	if err != nil {
		failParams := models.Response{
			Success: false,
			Debug:   "If this error persists please submit a bug report",
			Message: "Unable to remove contributor.",
		}
		failParams.Send(w, 500)
		return
	}

	resp := models.Response{
		Success: true,
		Message: "Removed successfully",
	}
	resp.Send(w, 200)
}
