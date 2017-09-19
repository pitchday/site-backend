package models

import (
	"fmt"
	"net/http"
	"net/url"
	"encoding/json"
	"time"
	"strconv"
	"errors"
	"github.com/pitchday/site-backend/config"
)

type Wallet struct {
	Balance int `json:"balance"`
}

func GetWalletBalance(walletAddress string) (balance int, err error){
	vc := ValueCache{
		Key: "walletBalance",
	}

	switch len(walletAddress){
	case 0:
		vc.GetByKey()
		if vc.ID == 0 || vc.ExpiresAt.Before(time.Now()) {
			Logger.Println("Wallet balance has expired... Updating now")

			newBalance := 0
			newBalance, err = GetUpdatedWalletBalance(config.Conf.WalletAddress)
			if err != nil {
				Logger.Println("Failed getting updated wallet balance for pd wallet,", err)
				return
			}

			vc.ExpiresAt = time.Now().Add(15 * time.Minute)
			vc.Value = fmt.Sprint(newBalance)
			err = vc.Update()
			if err != nil {
				return
			}
			err = nil
		}

		balance, _ = strconv.Atoi(vc.Value)
	case 40, 42:
		balance, err = GetUpdatedWalletBalance(walletAddress)
		if err != nil {
			Logger.Printf("Failed getting updated wallet balance for wallet '%s', %s\n", walletAddress, err)
		}

	default:
		err = errors.New("Invalid wallet")
		Logger.Printf("Got request with invalid wallet '%s'\n", walletAddress)
	}


	return
}

func GetUpdatedWalletBalance(walletAddress string) (balance int, err error){
	urlTemplate := "https://api.blockcypher.com/v1/eth/main/addrs/%s/balance"
	address := fmt.Sprintf(urlTemplate, walletAddress)

	requestUrl, err := url.Parse(address)
	if err != nil {
		Logger.Println("There was an error parsing the adress")
		return
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", requestUrl.String(), nil)

	//Do request
	resp, err := client.Do(req)
	if err != nil {
		Logger.Printf("There was an error retrieving the wallet balance: %s\n", err)
		return
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var c Wallet

	err = decoder.Decode(&c)
	if err != nil {
		Logger.Println("There was an error parsing the response while updating the wallet balance:", err)
		return
	}

	balance = c.Balance

	return
}
