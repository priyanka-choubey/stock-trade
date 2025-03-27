package dbmanager

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

type MySqlDatabase struct {
	Db *sql.DB
}

type DatabaseInterface interface {
	CheckUserLoginDetails(username string) (LoginDetails, error)
	CreateUserLoginDetails(username string, token string) (LoginDetails, error)
	UpdateUserLoginDetails(username string, token string) error
	DeleteUserLoginDetails(usernamr string) error
	SetupDatabase() error
}

// Database collections
type LoginDetails struct {
	AuthToken string
	Username  string
}

func SetupDatabase() (*MySqlDatabase, error) {
	var database MySqlDatabase

	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   os.Getenv("DBURL"),
		DBName: "trade",
	}
	// Get a database handle.
	var err error
	database.Db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
		return &database, err
	}

	pingErr := database.Db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
		return &database, fmt.Errorf("cannot connect to the database")
	}
	log.Debug("Connected to the databse!")
	return &database, nil
}

func (d *MySqlDatabase) CheckUserLoginDetails(username string, token string) error {

	var exists bool

	row := d.Db.QueryRow("SELECT username,token FROM user WHERE username = ? AND token = ?", username, token)
	if err := row.Scan(&exists); err != nil || err == sql.ErrNoRows {
		return fmt.Errorf("user %s: invalid credentials: %v", username, err)
	}
	return nil
}

func (d *MySqlDatabase) CreateUserLoginDetails(username string, token string) error {

	var exists bool

	d.Db.QueryRow("INSERT INTO user (username,token) VALUES (?,?)", username, token)
	row := d.Db.QueryRow("SELECT username,token FROM user WHERE username = ? AND token = ?", username, token)
	if err := row.Scan(&exists); err != nil || err == sql.ErrNoRows {
		return fmt.Errorf("user %s: cannot create user: %v", username, err)
	}
	return nil
}

func (d *MySqlDatabase) UpdateUserLoginDetails(username string, token string) error {

	var clientData LoginDetails

	d.Db.QueryRow("UPDATE user SET token = ? WHERE username = ?", token, username)
	row := d.Db.QueryRow("SELECT username,token FROM user WHERE username = ?", username)
	if err := row.Scan(&clientData.Username, &clientData.AuthToken); err != nil {
		return fmt.Errorf("unexpected error: %v", err)
	}

	if clientData.AuthToken != token {
		return fmt.Errorf("user %s: cannot update user token", clientData.Username)
	}
	return nil
}

func (d *MySqlDatabase) DeleteUserLoginDetails(username string) error {

	var clientData LoginDetails

	d.Db.QueryRow("DELETE FROM user WHERE username = ?", username)
	row := d.Db.QueryRow("SELECT (username,token) FROM user WHERE username = ?", username)
	if err := row.Scan(&clientData.Username, &clientData.AuthToken); err != nil {
		return nil
	}
	return fmt.Errorf("user %s: cannot delete user", clientData.Username)
}
