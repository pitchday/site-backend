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
		return

	case update.Message.Chat.IsSuperGroup():
		telegramSuperGroupHandler(update)
	}
}

func telegramPrivateHandler(update tgbotapi.Update) {
	if update.Message.IsCommand() {
		handleTelegramCommand(update)
	}

}

func handleTelegramCommand(update tgbotapi.Update) (err error) {
	switch update.Message.Command() {
	case "start":
		go func(){
			user := *update.Message.From
			MakeContributorFromTelegram(user, false, update.Message.CommandArguments(), update.Message.Chat.ID)
		}()
		sendMessate(update.Message.Chat.ID, "WELCOME please join https://t.me/joinchat/CIYV7EMVhacy2y4KKbRveA\n You can always ask for /help")

	case "help":

	case "refer":
		contributor := Contributor{}
		err, isNotFound := contributor.GetByServiceId(fmt.Sprint(update.Message.From.ID))
		if err != nil {
			if isNotFound {

			}
			sendMessate(update.Message.Chat.ID,  "It seems there was a problem retrieving your referral code. If you have seen this message before please report on the community channel")
			return
		}

		sendMessate(update.Message.Chat.ID, fmt.Sprintf("Your referral code is: %s You can also share this link with your friend.", contributor.ReferralCode))
		sendMessate(update.Message.Chat.ID, fmt.Sprintf("https://t.me/PitcherBot?start=%s", contributor.ReferralCode))
	}

	return
}

func telegramSuperGroupHandler(update tgbotapi.Update) {
	if update.Message.Chat.ID != config.Conf.CommunityTelegramGroupId {
		Logger.Println("Received message from unknown group. The given Id is:", update.Message.Chat.ID)
		return
	}

	Logger.Println("Received message on group.")

	if update.Message != nil {
		if update.Message.NewChatMembers != nil && update.Message.Chat.ID == config.Conf.CommunityTelegramGroupId {
			Logger.Println("New user on group.")

			var users []tgbotapi.User
			users = *update.Message.NewChatMembers

			for _, user := range users {
				MakeContributorFromTelegram(user, true, "", 0)
			}
		}

		if update.Message.LeftChatMember != nil && update.Message.Chat.ID == config.Conf.CommunityTelegramGroupId {
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

func sendMessate(chatId int64, message string) (){
	msg := tgbotapi.NewMessage(chatId, message)

	_, err := telegramBot.Send(msg)
	if err != nil {
		Logger.Println("Got error while sending message")
	}
	return
}