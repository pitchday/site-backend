package models

import (
	"github.com/cheviz/pitchdayBackend/config"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

func TelegramHandler(update tgbotapi.Update, botId string) {
	Logger.Println("Got request from bot", botId)

	switch update.Message.Chat.Type {
	case "private":
		telegramPrivateHandler(update)

	case "group":

	case "supergroup":
		telegramSuperGroupHandler(update)
	}
}

func telegramPrivateHandler(update tgbotapi.Update) {
	messageString := update.Message.Text

	Logger.Printf("%s says: %s", update.Message.From.UserName, messageString)

	if update.Message.IsCommand() {
		handleTelegramCommand(update.Message.Command())
	}

}

func handleTelegramCommand(msg string) (err error) {
	message := strings.SplitN(msg, " ", 1)

	switch message[0] {
	case "start":
		Logger.Println("User registered")
	}

	return
}

func telegramSuperGroupHandler(update tgbotapi.Update) {
	if update.Message != nil {
		if update.Message.NewChatMembers != nil && update.Message.Chat.Title == config.Conf.CommunityTelegramGroupName {
			Logger.Print("Registering new members")
			var users []tgbotapi.User
			users = *update.Message.NewChatMembers

			for _, user := range users {

				Logger.Print("Registering members", user.ID, user.FirstName, user.LastName)
				MakeContributorFromTelegram(user)
			}

		}

		if update.Message.LeftChatMember != nil && update.Message.Chat.Title == config.Conf.CommunityTelegramGroupName {
			err := DeleteContributorUsingTelegramId(update.Message.LeftChatMember.ID)
			if err != nil {
				Logger.Println("Got error:", err)
			}
		}
		//Logger.Println("The message is:", update.Message.Text)
		Logger.Println("Sent by:", update.Message.From.UserName, "with id:", update.Message.From.ID)
		//Logger.Println("The id is:", update.Message.Chat.ID, "The name is:", update.Message.Chat.Title, "The type is:", update.Message.Chat.Type)

	}
}

func GetMemberCountInChannel(channelName string) (count int64, err error) {
	groupData := tgbotapi.ChatConfig{
		ChatID: 0,
		SuperGroupUsername: "",
	}

	telegramBot.GetChatMembersCount(groupData)
}