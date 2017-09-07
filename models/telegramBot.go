package models

import (
	"github.com/cheviz/pitchdayBackend/config"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func TelegramHandler(update tgbotapi.Update, botId string) {
	Logger.Println("Got request from bot", botId)

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
