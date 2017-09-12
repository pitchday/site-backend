package api

import (
	"github.com/pitchday/site-backend/models"
	"net/http"
	"strconv"
)

func Get_Announcements(w http.ResponseWriter, r *http.Request) {
	announcements := models.Announcements{}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	err := announcements.Get(limit, offset)
	if err != nil {
		resp := models.Response{
			Success: false,
			Debug:   "If this error persists please submit a bug report",
			Message: "Unable to retrieve announcements",
		}
		resp.Send(w, 500)
		return
	}

	resp := models.Response{
		Success: true,
		Message: "Retrieved successfully",
		Data:    announcements,
	}
	resp.Send(w, 200)
	return
}
