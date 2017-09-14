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
		user := *update.Message.From

		contributor := Contributor{}
		err = contributor.GetByServiceId(fmt.Sprint(user.ID))
		if err != nil {
			if isNotFoundDBError(err){
				MakeContributorFromTelegram(user, false, update.Message.CommandArguments(), update.Message.Chat.ID)
			}
		}

		if contributor.PrivateChat == 0 {
			contributor.PrivateChat = update.Message.Chat.ID
			if len(contributor.ReferralCode) < 1 {
				contributor.ReferralCode = RandStringBytesMaskImprSrc(config.Conf.ReferalTokenLength, true)
			}

			contributor.Update()
		}

		markup := tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonURL(config.Conf.BotMessages["joinCommunityLabel"], config.Conf.CommunityTelegramGroupLink)})
		sendMessage(update.Message.Chat.ID, fmt.Sprintf(config.Conf.BotMessages["start"], config.Conf.CommunityTelegramGroupLink), markup)


	case "menu":
		markup := tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(config.Conf.BotMessages["earnTokensLabel"], "earnTokens")})
		sendMessage(update.Message.Chat.ID, config.Conf.BotMessages["welcomeToCommunity"], markup)

	//	TODO make this useful
	//case "refer":
	//	contributor := Contributor{}
	//	err = contributor.GetByServiceId(fmt.Sprint(update.Message.From.ID))
	//	if err != nil {
	//		sendMessage(update.Message.Chat.ID,  "It seems there was a problem retrieving your referral code. If you have seen this message before please report on the community channel", nil)
	//		//TODO notify admins
	//		return
	//	}
	//	sendMessage(update.Message.Chat.ID, fmt.Sprintf("Your referral code is: %s You can also share this link with your friend.", contributor.ReferralCode), nil)
	//	sendMessage(update.Message.Chat.ID, fmt.Sprintf("https://t.me/PitcherBot?start=%s", contributor.ReferralCode), nil)
	}

	return
}

func handleTelegramQuery(update tgbotapi.Update) {
	user := Contributor{}

	err := user.GetByServiceId(fmt.Sprint(update.CallbackQuery.From.ID))
	if err != nil {
		if isNotFoundDBError(err) {

		}
		Logger.Println("Got an error while retrieving user.", err)
	}

	if (len(user.ReferralCode) < 1 && user.Id > 0) || (len(user.ReferralCode) > config.Conf.ReferalTokenLength){
		user.ReferralCode = RandStringBytesMaskImprSrc(config.Conf.ReferalTokenLength, true)
		err = user.Update()
		if err != nil {
			Logger.Println("Failed to update user without referral code:", err)
		}
	}

	markup := tgbotapi.InlineKeyboardMarkup{}
	msg := ""

	switch update.CallbackQuery.Data {
	case "buyTokens":
		msg = config.Conf.BotMessages["buyTokensMessage"]
		markup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("Contact us", config.Conf.PRUrl)),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Back", "default")),
		)

	case "earnTokens":
		msg = config.Conf.BotMessages["earnTokensMessage"]
		markup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(config.Conf.BotMessages["bringNewMembersLabel"], "members")),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL(config.Conf.BotMessages["joinBountyLabel"], config.Conf.BotMessages["bountyPage"])),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL(config.Conf.BotMessages["joinDesignLabel"], config.Conf.BotMessages["designPage"])),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL(config.Conf.BotMessages["joinEngineeringLabel"], config.Conf.BotMessages["engineeringPage"])),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(config.Conf.BotMessages["beAdvisorLabel"], "advise")),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Back", "default")),
		)

	case "members":
		msg = fmt.Sprintf(config.Conf.BotMessages["bringNewMembersMessage"], user.ReferralCode)
		buttons := []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonURL("More info", config.Conf.BotMessages["bringNewMembersUrl"]),
			tgbotapi.NewInlineKeyboardButtonData("Back", "earnTokens"),
		}
		markup = tgbotapi.NewInlineKeyboardMarkup(buttons)

	case "advise":
		msg = fmt.Sprintf(config.Conf.BotMessages["beAdvisorMessage"], config.Conf.CommunityTelegramGroupLink)
		buttons := []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonURL("Go to community", config.Conf.CommunityTelegramGroupLink),
			tgbotapi.NewInlineKeyboardButtonData("Back", "earnTokens"),
		}
		markup = tgbotapi.NewInlineKeyboardMarkup(buttons)

	default:
		//MAIN MENU
		msg = config.Conf.BotMessages["welcomeToCommunity"]
		buttons := []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(config.Conf.BotMessages["earnTokensLabel"], "earnTokens"),
			tgbotapi.NewInlineKeyboardButtonData(config.Conf.BotMessages["buyTokensLabel"], "buyTokens"),
		}
		markup = tgbotapi.NewInlineKeyboardMarkup(buttons)
	}
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
	msg.DisableWebPagePreview = true

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