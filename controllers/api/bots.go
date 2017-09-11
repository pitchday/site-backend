package api

import (
	"encoding/json"
	"github.com/pitchday/site-backend/models"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gorilla/mux"
	"net/http"
)

func Bot_Hook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	botId := vars["botId"]

	decoder := json.NewDecoder(r.Body)
	var update tgbotapi.Update

	err := decoder.Decode(&update)
	if err != nil {
		Logger.Println(err)
		resp := models.Response{
			Success: false,
			Debug:   "Cannot parse request",
			Message: "Message not received",
		}
		resp.Send(w, 400)
		return

	}

	resp := models.Response{
		Success: true,
		Message: "Received successfully",
	}
	resp.Send(w, 200)

	go func() {
		models.TelegramHandler(update, botId)
	}()
	return
}

func Get_Group_Member_Count(w http.ResponseWriter, r *http.Request) {
	telegramCount, err := models.GetMemberCountInChannel()
	if err != nil {
		Logger.Println(err)
		resp := models.Response{
			Success: false,
			Debug:   "There was an error getting the telegram count",
			Message: "Unable to get count",
		}
		resp.Send(w, 400)
		return
	}

	counts := struct {
		TelegramCount int `json:"telegramCount"`
	}{
		telegramCount,
	}

	resp := models.Response{
		Success: true,
		Message: "Received successfully",
		Data:    counts,
	}
	resp.Send(w, 200)
}
