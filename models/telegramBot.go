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

	case update.CallbackQuery != nil:
		handleTelegramQuery(update)
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
	switch {
	case update.Message.IsCommand():
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
		markup := tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonURL(config.Conf.BotMessages["joinCommunityLabel"], config.Conf.CommunityTelegramGroupLink)})
		sendMessage(update.Message.Chat.ID, config.Conf.BotMessages["start"], markup)


	case "help":

	case "refer":
		contributor := Contributor{}
		err = contributor.GetByServiceId(fmt.Sprint(update.Message.From.ID))
		if err != nil {
			sendMessage(update.Message.Chat.ID,  "It seems there was a problem retrieving your referral code. If you have seen this message before please report on the community channel", nil)
			//TODO notify admins
			return
		}
		sendMessage(update.Message.Chat.ID, fmt.Sprintf("Your referral code is: %s You can also share this link with your friend.", contributor.ReferralCode), nil)
		sendMessage(update.Message.Chat.ID, fmt.Sprintf("https://t.me/PitcherBot?start=%s", contributor.ReferralCode), nil)
	}

	return
}

func handleTelegramQuery(update tgbotapi.Update) {
	markup := tgbotapi.InlineKeyboardMarkup{}
	msg := ""

	switch update.CallbackQuery.Data {
	default:
		//MAIN MENU
		msg = config.Conf.BotMessages["welcomeToCommunity"]
		buttons := []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(config.Conf.BotMessages["earnTokensLabel"], "earnTokens"),
		}
		markup = tgbotapi.NewInlineKeyboardMarkup(buttons)

	case "earnTokens":
		msg = config.Conf.BotMessages["earnTokensMessage"]

		markup.InlineKeyboard = [][]tgbotapi.InlineKeyboardButton{
			[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(config.Conf.BotMessages["bringNewMembersLabel"], "members")},
			[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(config.Conf.BotMessages["joinBountyLabel"], "bounty")},
			[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(config.Conf.BotMessages["joinDesignLabel"], "design")},
			[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(config.Conf.BotMessages["joinEngineeringLabel"], "engineering")},
			[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(config.Conf.BotMessages["beAdvisorLabel"], "advise")},
			[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("Back", "")},
		}

	case "members":
		msg = config.Conf.BotMessages["bringNewMembersMessage"]
		buttons := []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("Back", "earnTokens"),
		}
		markup = tgbotapi.NewInlineKeyboardMarkup(buttons)


	case "bounty":
		msg = config.Conf.BotMessages["bringNewMembersMessage"]
		buttons := []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("Back", "earnTokens"),
		}
		markup = tgbotapi.NewInlineKeyboardMarkup(buttons)
	}
	Logger.Println("GOT MESSAGE FROM:", update.CallbackQuery.From.ID, "IN CHAT", update.CallbackQuery.Message.Chat.ID, "MessageId is", update.CallbackQuery.Message.MessageID, "DATA IS", update.CallbackQuery.Data)

	updateMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, msg)
	updateKeyboard(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, markup)
	return
}

func telegramSuperGroupHandler(update tgbotapi.Update) {
	if update.Message.Chat.ID != config.Conf.CommunityTelegramGroupId {
		Logger.Println("Received message from unknown group. The given Id is:", update.Message.Chat.ID)
		return
	}

	if update.Message == nil{
		return
	}

	switch {
	case update.Message.NewChatMembers != nil:
		var users []tgbotapi.User
		users = *update.Message.NewChatMembers
		for _, user := range users {
			contributor := Contributor{}
			err := contributor.GetByServiceId(fmt.Sprint(user.ID))
			if err != nil {
				if isNotFoundDBError(err) {
					// Register user as community member
					MakeContributorFromTelegram(user, true, "", 0)
				} else {
					//	TODO notify admins
				}
			}

			if contributor.DeletedAt.Valid {
				err = contributor.MakeMember()
				if err != nil {
					//	TODO notify admins
				}

				if contributor.PrivateChat != 0 {
				//	TODO send message on private chat
					markup := tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(config.Conf.BotMessages["earnTokensLabel"], "earnTokens")})
					sendMessage(contributor.PrivateChat, config.Conf.BotMessages["welcomeToCommunity"], markup)
				}
				return
			}
		}
	case update.Message.LeftChatMember != nil:
		err := DeleteContributorUsingTelegramId(update.Message.LeftChatMember.ID)
		if err != nil {
			Logger.Println("Got error:", err)
			//	TODO notify admins
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

func sendMessage(chatId int64, message string, markup interface{}) (){
	msg := tgbotapi.NewMessage(chatId, message)
	msg.DisableWebPagePreview = true

	if markup != nil {
		msg.ReplyMarkup = markup
	}

	_, err := telegramBot.Send(msg)
	if err != nil {
		Logger.Println("Got error while sending message", err)
	}
	return
}

func updateMessage(chatId int64, messageId int, message string) (){
	msg := tgbotapi.NewEditMessageText(chatId, messageId, message)

	_, err := telegramBot.Send(msg)
	if err != nil {
		Logger.Println("Got error while sending message", err)
	}
	return
}

func updateKeyboard(chatId int64, messageId int, markup tgbotapi.InlineKeyboardMarkup) (){
	msg := tgbotapi.NewEditMessageReplyMarkup(chatId, messageId, markup)

	_, err := telegramBot.Send(msg)
	if err != nil {
		Logger.Println("Got error while sending message", err)
	}
	return
}