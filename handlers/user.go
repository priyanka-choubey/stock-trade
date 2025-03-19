package handlers

import (
	"github.com/priyanka-choubey/stock-trade/db_manager"
)

type status struct {
	code    int
	message string
}

func CreateUser(username string, token string) status {

	var resp status

	db, err := db_manager.SetupDatabase()
	if err != nil {
		resp.getStatus(500, "Internal error")
	}

	err = db.CreateUserLoginDetails(username, token)
	if err != nil {
		resp.getStatus(403, err)
	}

	resp.getStatus(200, "OK")
	defer db.Close()
	return resp
}

func AuthenticateUser(username string, token string) status {

	var resp status

	db, err := db_manager.SetupDatabase()
	if err != nil {
		resp.getStatus(500, "Internal error")
	}

	err = db.CheckUserLoginDetails(username, token)
	if err != nil {
		resp.getStatus(401, err)
	}
	resp.getStatus(200, "OK")
	defer db.Close()
	return resp
}

func (s status) getStatus(code int, message string) {
	s.code = code
	s.message = message
}
