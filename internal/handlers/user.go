package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/priyanka-choubey/balance/api"
	"github.com/priyanka-choubey/balance/internal/tools"
	log "github.com/sirupsen/logrus"
)

var usernameInUseError = errors.New("Given username is already in use")
var improperCredentialsError = errors.New("Username or token cannot be empty")
var invalidCredentialsError = errors.New("Username or token is invalid")

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var params = api.UserParams{}
	var decoder *schema.Decoder = schema.NewDecoder()
	var err error

	err = decoder.Decode(&params, r.URL.Query())

	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	if params.Username == "" || params.Token == "" {
		log.Error(improperCredentialsError)
		api.RequestErrorHandler(w, improperCredentialsError)
		return
	}

	var database *tools.MySqlDatabase
	database, err = tools.NewDatabase()
	if err != nil {
		api.InternalErrorHandler(w)
		return
	}

	_, err = (*database).GetUserLoginDetails(params.Username)
	if err == nil {
		log.Error(usernameInUseError)
		api.RequestErrorHandler(w, usernameInUseError)
		return
	} else {

		loginDetails, err := (*database).CreateUserLoginDetails(params.Username, params.Token)
		if err != nil {
			api.RequestErrorHandler(w, err)
			return
		}

		coinDetails, err := (*database).CreateAccountBalanceDetails(params.Username)
		if err != nil {
			api.RequestErrorHandler(w, err)
			return
		}

		if loginDetails == nil {
			api.InternalErrorHandler(w)
			return
		}

		var response = api.UserResponse{
			Username: loginDetails.Username,
			Token:    loginDetails.AuthToken,
			Balance:  int(coinDetails.Coins),
			Code:     http.StatusOK,
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Error(err)
			api.InternalErrorHandler(w)
			return
		}
	}

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params, err := getParams(w, r)
	if err != nil {
		return
	}

	database, err := getDatabase(w, r)
	if err != nil {
		return
	}

	err = (*database).DeleteUserLoginDetails(params.Username)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	err = (*database).DeleteAccountBalanceDetails(params.Username)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	var response = api.Response{
		Message: "User " + params.Username + " has been deleted",
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

func UpdateUserToken(w http.ResponseWriter, r *http.Request) {
	params, err := getParams(w, r)
	if err != nil {
		return
	}

	database, err := getDatabase(w, r)
	if err != nil {
		return
	}

	err = (*database).UpdateUserLoginDetails(params.Username, params.Token)
	if err != nil {
		api.RequestErrorHandler(w, err)
		return
	}

	var response = api.Response{
		Message: "Token for user " + params.Username + " has been updated",
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

func getParams(w http.ResponseWriter, r *http.Request) (api.UserParams, error) {
	var params = api.UserParams{}
	var decoder *schema.Decoder = schema.NewDecoder()
	var err error

	err = decoder.Decode(&params, r.URL.Query())

	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return params, err
	}

	return params, err

}

func getDatabase(w http.ResponseWriter, r *http.Request) (*tools.MySqlDatabase, error) {
	var database *tools.MySqlDatabase
	var err error

	database, err = tools.NewDatabase()
	if err != nil {
		api.InternalErrorHandler(w)
		return database, err
	}

	return database, err
}
