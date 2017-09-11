package models

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pitchday/site-backend/config"
)

func TelegramHandler(update tgbotapi.Update, botId string) {
	switch {
	case update.Message != nil:
		telegramHandleMessage(update)

	case update.ChannelPost != nil, update.EditedChannelPost != nil:
		telegramHandleChannel(update)
	}

	return
}

func telegramHandleChannel(update tgbotapi.Update) {
	var message *tgbotapi.Message
	var isEdit bool

	switch {
	case update.ChannelPost != nil:
		message = update.ChannelPost
	case update.EditedChannelPost != nil:
		message = update.EditedChannelPost
		isEdit = true
	default:
		return
	}

	Logger.Println("GOT MESSAGE:", message, isEdit)

	if message.Chat.ID != config.Conf.AnnouncementsTelegramChannelId {
		Logger.Println("Received message from unknown channel. The given Id is:", message.Chat.ID)
		return
	}

	announcement := Announcement{
		MessageId: message.MessageID,
		ChannelId: message.Chat.ID,
	}

	switch {
	default:
		return
	case message.Photo != nil:
		images := *message.Photo

		url, err := telegramBot.GetFileDirectURL(images[2].FileID)
		if err != nil {
			Logger.Printf("On announcement got %s while requesting image link", fmt.Sprint(err))
		}
		announcement.ResourceType = "img"
		announcement.Resource = url
		fallthrough
	case message.Text != "":
		announcement.Body = message.Text
	}

	if isEdit {
		err := announcement.Edit()
		if err != nil {
			Logger.Println("Got an error while updating announcement, the error is", err)
		}
	} else {
		err := announcement.Create()
		if err != nil {
			Logger.Println("Got an error while storing announcement, the error is", err)
		}
	}

	return
}

func telegramHandleMessage(update tgbotapi.Update) {
	switch {
	case update.Message.Chat.IsPrivate():
		telegramPrivateHandler(update)

	case update.Message.Chat.IsGroup():

	case update.Message.Chat.IsSuperGroup():
		telegramSuperGroupHandler(update)
	}
}

func telegramPrivateHandler(update tgbotapi.Update) {
	if update.Message.IsCommand() {
		handleTelegramCommand(update.Message.Command())
	}

}

func handleTelegramCommand(command string) (err error) {
	switch command {
	case "start":
		Logger.Println("User registered")
	}

	return
}

func telegramSuperGroupHandler(update tgbotapi.Update) {
	if update.Message.Chat.ID != config.Conf.CommunityTelegramGroupId {
		Logger.Println("Received message from unknown group. The given Id is:", update.Message.Chat.ID)
		return
	}

	if update.Message != nil {
		if update.Message.NewChatMembers != nil && update.Message.Chat.Title == config.Conf.CommunityTelegramGroupName {
			var users []tgbotapi.User
			users = *update.Message.NewChatMembers

			for _, user := range users {
				MakeContributorFromTelegram(user)
			}
		}

		if update.Message.LeftChatMember != nil && update.Message.Chat.Title == config.Conf.CommunityTelegramGroupName {
			err := DeleteContributorUsingTelegramId(update.Message.LeftChatMember.ID)
			if err != nil {
				Logger.Println("Got error:", err)
			}
		}
	}
}

func GetMemberCountInChannel() (count int, err error) {
	groupData := tgbotapi.ChatConfig{
		ChatID:             config.Conf.CommunityTelegramGroupId,
		SuperGroupUsername: "",
	}

	count, err = telegramBot.GetChatMembersCount(groupData)
	count = count - 1
	return
}
