package tools

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

// Database collections
type LoginDetails struct {
	AuthToken string
	Username  string
}
type CoinDetails struct {
	Coins    int
	Username string
}

type MySqlDatabase struct {
	Db *sql.DB
}

type DatabaseInterface interface {
	GetUserLoginDetails(username string) (LoginDetails, error)
	GetUserCoins(username string) (CoinDetails, error)
	SetupDatabase() error
}

func NewDatabase() (*MySqlDatabase, error) {
	var database MySqlDatabase

	err := database.SetupDatabase()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &database, nil
}

func (d *MySqlDatabase) CreateUserLoginDetails(username string, token string) (*LoginDetails, error) {

	var clientData LoginDetails

	d.Db.QueryRow("INSERT INTO user (username,token) VALUES (?,?)", username, token)
	row := d.Db.QueryRow("SELECT username,token FROM user WHERE username = ?", username)
	if err := row.Scan(&clientData.Username, &clientData.AuthToken); err != nil {
		return &clientData, fmt.Errorf("User %d: Cannot create user: %v", username, err)
	}
	return &clientData, nil
}

func (d *MySqlDatabase) CreateAccountBalanceDetails(username string) (*CoinDetails, error) {

	var clientData CoinDetails

	d.Db.QueryRow("INSERT INTO balance (username) VALUES (?)", username)
	row := d.Db.QueryRow("SELECT * FROM balance WHERE username = ?", username)
	if err := row.Scan(&clientData.Username, &clientData.Coins); err != nil {
		return &clientData, fmt.Errorf("User %d: Cannot create user: %v", username, err)
	}
	return &clientData, nil
}

func (d *MySqlDatabase) DeleteUserLoginDetails(username string) error {

	var clientData LoginDetails

	d.Db.QueryRow("DELETE FROM user WHERE username = ?", username)
	row := d.Db.QueryRow("SELECT (username,token) FROM user WHERE username = ?", username)
	if err := row.Scan(&clientData.Username, &clientData.AuthToken); err != nil {
		return nil
	}
	return fmt.Errorf("User %d: Cannot delete user", clientData.Username)
}

func (d *MySqlDatabase) DeleteAccountBalanceDetails(username string) error {

	var clientData CoinDetails

	d.Db.QueryRow("DELETE FROM balance WHERE username = ?", username)
	row := d.Db.QueryRow("SELECT * FROM balance WHERE username = ?", username)
	if err := row.Scan(&clientData.Username, &clientData.Coins); err != nil {
		return nil
	}
	return fmt.Errorf("User %d: Cannot delete user", clientData.Username)
}

func (d *MySqlDatabase) UpdateUserLoginDetails(username string, token string) error {

	var clientData LoginDetails

	d.Db.QueryRow("UPDATE user SET token = ? WHERE username = ?", token, username)
	row := d.Db.QueryRow("SELECT username,token FROM user WHERE username = ?", username)
	if err := row.Scan(&clientData.Username, &clientData.AuthToken); err != nil {
		return fmt.Errorf("Unexpected Error: %v", err)
	}

	if clientData.AuthToken != token {
		return fmt.Errorf("User %d: Cannot update user token", clientData.Username)
	}
	return nil
}

func (d *MySqlDatabase) GetUserLoginDetails(username string) (*LoginDetails, error) {

	var clientData LoginDetails

	row := d.Db.QueryRow("SELECT username,token FROM user WHERE username = ?", username)
	if err := row.Scan(&clientData.Username, &clientData.AuthToken); err != nil {
		if err == sql.ErrNoRows {
			return &clientData, fmt.Errorf("User %d: no such user", username)
		}
		return &clientData, fmt.Errorf("User %d: %v", username, err)
	}
	return &clientData, nil
}

func (d *MySqlDatabase) GetUserCoins(username string) (*CoinDetails, error) {

	var clientData = CoinDetails{}
	row := d.Db.QueryRow("SELECT * FROM balance WHERE username = ?", username)
	if err := row.Scan(&clientData.Username, &clientData.Coins); err != nil {
		if err == sql.ErrNoRows {
			return &clientData, fmt.Errorf("User %d: no such user", username)
		}
		return &clientData, fmt.Errorf("User %d: %v", username, err)
	}
	return &clientData, nil
}

func (d *MySqlDatabase) UpdateAccountBalance(username string, balance int) error {

	var clientData CoinDetails

	d.Db.QueryRow("UPDATE balance SET balance = ? WHERE username = ?", balance, username)
	row := d.Db.QueryRow("SELECT username,balance FROM balance WHERE username = ?", username)
	if err := row.Scan(&clientData.Username, &clientData.Coins); err != nil {
		return fmt.Errorf("Unexpected Error: %v", err)
	}

	if clientData.Coins != balance {
		return fmt.Errorf("User %d: Cannot update accountbalance", clientData.Username)
	}
	return nil
}

func (d *MySqlDatabase) SetupDatabase() error {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "trade",
	}
	// Get a database handle.
	var err error
	d.Db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := d.Db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	log.Debug("Connected to the databse!")
	return nil
}
