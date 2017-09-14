package models

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pitchday/site-backend/config"
	"log"
	"os"
)

var db *gorm.DB
var telegramBot *tgbotapi.BotAPI
var Logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

// Setup initializes the Conn object
func Setup() error {
	var err error
	db, err = gorm.Open("mysql", config.Conf.DBString)
	if err != nil {
		return err
	}
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.LogMode(true)

	db.AutoMigrate(&NewsletterSubscription{})
	db.AutoMigrate(&Contributor{})
	db.AutoMigrate(&Announcement{})

	telegramBot, err = tgbotapi.NewBotAPI(config.Conf.TelegramBotToken)
	if err != nil {
		Logger.Println(err)
		return err
	}
	telegramBot.Debug = false
	log.Printf("Authorized on account %s", telegramBot.Self.UserName)

	_, err = telegramBot.SetWebhook(tgbotapi.NewWebhook(config.Conf.TelegramBotWebHook + config.Conf.TelegramBotToken))
	if err != nil {
		Logger.Println(err)
	}

	return err
}

func isDuplicatedDBError(err error) (isDuplicated bool) {
	if fmt.Sprint(err)[:27] == "Error 1062: Duplicate entry" {
		isDuplicated = true
	}
	return
}

func isNotFoundDBError(err error) (isDuplicated bool) {
	if fmt.Sprint(err) == "record not found" {
		isDuplicated = true
	}
	return
}
