package handler

import (
	"backend/config"
	"backend/dto"
	"backend/model"
	"backend/repo"
	"backend/service"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func AddEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value(config.DBContextKey).(*sql.DB)

	var req dto.AddEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.EmployeeID == nil || req.EmployeeName == nil || req.JoiningDate == nil || req.Password == nil || req.IsAdmin == nil {
		http.Error(w, "空欄があります", http.StatusBadRequest)
		return
	}

	employee, err := repo.GetEmployeeByID(db, *req.EmployeeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if employee != (model.Employee{}) {
		http.Error(w, "すでに同じIDでの登録があります", http.StatusConflict)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			tx.Rollback()
		} else {
			tx.Commit()
			fmt.Fprintln(w, "successfully added the user")
		}
	}()

	err = service.AddEmployee(tx, *req.EmployeeID, *req.IsAdmin, *req.EmployeeName, req.JoiningDate.Time)
	if err != nil {
		return
	}
	if err = service.AddAuth(tx, *req.EmployeeID, *req.Password); err != nil {
		return
	}
}

func GetLoggedInEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	employee := r.Context().Value(config.LoginUserContextKey).(model.Employee)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(employee); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
