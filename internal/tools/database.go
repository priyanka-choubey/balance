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
	Coins    int64
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
