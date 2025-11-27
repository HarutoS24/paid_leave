package service

import (
	"backend/model"
	"backend/repo"
	"database/sql"
	"fmt"
	"time"

	"github.com/gorilla/sessions"
)

func GetLoggedInEmployee(db *sql.DB, session *sessions.Session) (model.Employee, error) {
	isLoggedIn, err := IsLoggedIn(session)
	if err != nil {
		return model.Employee{}, err
	}
	if !isLoggedIn {
		return model.Employee{}, nil
	}

	employeeIDRaw, ok := session.Values["employeeID"]
	if !ok {
		return model.Employee{}, nil
	}
	employeeID, ok := employeeIDRaw.(string)
	if !ok {
		return model.Employee{}, fmt.Errorf("invalid value found in session")
	}
	employee, err := repo.GetEmployeeByID(db, employeeID)
	if err != nil {
		return model.Employee{}, err
	}
	return employee, nil
}
