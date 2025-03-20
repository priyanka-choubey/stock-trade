package handlers

import (
	"fmt"

	dbmanager "github.com/priyanka-choubey/stock-trade/db_manager"
)

type status struct {
	Code    int
	Message string
}

func CreateUser(username string, token string) status {

	var resp status

	db, err := dbmanager.SetupDatabase()
	if err != nil {
		return resp.getStatus(500, "Internal error")
	}

	err = db.CreateUserLoginDetails(username, token)
	if err != nil {
		return resp.getStatus(403, fmt.Sprintf("%v", err))
	}

	resp = resp.getStatus(200, "OK")
	defer db.Db.Close()
	return resp
}

func AuthenticateUser(username string, token string) status {

	var resp status

	db, err := dbmanager.SetupDatabase()
	if err != nil {
		return resp.getStatus(500, "Internal error")
	}

	err = db.CheckUserLoginDetails(username, token)
	if err != nil {
		fmt.Println("%v", resp.getStatus(401, fmt.Sprintf("%v", err)))
		return resp.getStatus(401, fmt.Sprintf("%v", err))
	}

	resp = resp.getStatus(200, "OK")
	defer db.Db.Close()
	return resp
}

func (s status) getStatus(code int, message string) status {
	s.Code = code
	s.Message = message
	return s
}
