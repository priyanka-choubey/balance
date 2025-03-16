package tools

import (
	"time"
)

type mockDB struct{}

var mockLoginDetails = map[string]LoginDetails{
	"ruhi": {
		AuthToken: "12345",
		Username:  "ruhi",
	},
	"kritika": {
		AuthToken: "67890",
		Username:  "kritika",
	},
	"priyanka": {
		AuthToken: "13579",
		Username:  "priyanka",
	},
}

var mockCoinDetails = map[string]CoinDetails{
	"ruhi": {
		Coins:    100,
		Username: "ruhi",
	},
	"kritika": {
		Coins:    200,
		Username: "kritika",
	},
	"priyanka": {
		Coins:    300,
		Username: "priyanka",
	},
}

func (d *mockDB) GetUserLoginDetails(username string) *LoginDetails {
	// Simulate DB call
	time.Sleep(time.Second * 1)

	var clientData = LoginDetails{}
	clientData, ok := mockLoginDetails[username]
	if !ok {
		return nil
	}

	return &clientData
}

func (d *mockDB) GetUserCoins(username string) *CoinDetails {
	// Simulate DB call
	time.Sleep(time.Second * 1)

	var clientData = CoinDetails{}
	clientData, ok := mockCoinDetails[username]
	if !ok {
		return nil
	}

	return &clientData
}

func (d *mockDB) SetupDatabase() error {
	return nil
}
