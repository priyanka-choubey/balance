package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/priyanka-choubey/balance/api"
	"github.com/priyanka-choubey/balance/internal/tools"
	log "github.com/sirupsen/logrus"
)

func GetBalance(w http.ResponseWriter, r *http.Request) {
	var params = api.BalanceParams{}
	var decoder *schema.Decoder = schema.NewDecoder()
	var err error

	err = decoder.Decode(&params, r.URL.Query())

	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	var database *tools.MySqlDatabase
	database, err = tools.NewDatabase()
	if err != nil {
		api.InternalErrorHandler(w)
		return
	}

	var tokenDetails *tools.CoinDetails
	tokenDetails, err = (*database).GetUserCoins(params.Username)
	if err != nil {
		api.InternalErrorHandler(w)
		return
	}
	if tokenDetails == nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	var response = api.BalanceResponse{
		Balance: int((*tokenDetails).Coins),
		Code:    http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}
