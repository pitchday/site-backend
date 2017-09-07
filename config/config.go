package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var Conf Config

type Config struct {
	ApiURL                     string `json:"apiUrl"`
	AccessControlAllowOrigin   string `json:"AccessControlAllowOrigin"`
	DBString                   string `json:"dbConnectionString"`
	ApiPrefix                  string `json:"apiPrefix"`
	AWSAccessKey               string `json:"awsAccessKey"`
	AWSSecretKey               string `json:"awsSecretKey"`
	TelegramBotWebHook         string `json:"telegramBotWebHook"`
	AWSBucketName              string `json:"awsBucketName"`
	TelegramBotToken           string `json:"telegramBotToken"`
	CommunityTelegramGroupName string `json:"communityTelegramGroupName"`
	CommunityTelegramGroupId   int64  `json:"communityTelegramGroupId"`
}

func init() {
	// Get the config file
	configFile, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
	}
	json.Unmarshal(configFile, &Conf)
}
