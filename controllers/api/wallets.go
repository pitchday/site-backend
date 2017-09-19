package api

import (
	"net/http"
	"github.com/pitchday/site-backend/models"
)

func Get_Wallet_Balance(w http.ResponseWriter, r *http.Request) {
	walletAddress := r.URL.Query().Get("wallet")

	walletBalance, err := models.GetWalletBalance(walletAddress)
	if err != nil {
		resp := models.Response{
			Success: false,
			Debug:   "There was an error getting the wallet balance",
			Message: "Unable to get balance",
		}
		resp.Send(w, 400)
		return
	}

	resp := models.Response{
		Success: true,
		Message: "Received successfully",
		Data:    walletBalance,
	}
	resp.Send(w, 200)
}
