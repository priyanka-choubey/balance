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
		log.Error()
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
		var loginDetails *tools.LoginDetails
		loginDetails, err = (*database).CreateUserLoginDetails(params.Username, params.Token)
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
