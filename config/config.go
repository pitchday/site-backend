package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var Conf Config

type Config struct {
	ApiURL    string `json:"apiUrl"`
	DBString  string `json:"dbConnectionString"`
	ApiPrefix string `json:"apiPrefix"`
}

func init() {
	// Get the config file
	configFile, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
	}
	json.Unmarshal(configFile, &Conf)
}
