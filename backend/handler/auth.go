package handler

import (
	"backend/config"
	"backend/dto"
	"backend/service"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value(config.DBContextKey).(*sql.DB)
	session := r.Context().Value(config.AuthContextKey).(*sessions.Session)
	var req = dto.LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.EmployeeID == nil || req.Password == nil {
		http.Error(w, "空欄があります", http.StatusBadRequest)
		return
	}

	isLoggedIn, err := service.IsLoggedIn(session)
	// sessionに異常があった場合session情報削除
	if err != nil {
		session.Values = make(map[interface{}]interface{})
		session.Options.MaxAge = -1
		if saveErr := session.Save(r, w); saveErr != nil {
			err = fmt.Errorf("%w reset error: %s", err, saveErr)
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if isLoggedIn {
		fmt.Fprintln(w, "already logged in")
		return
	}

	valid, err := service.ValidateUserByPassword(db, *req.EmployeeID, *req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !valid {
		fmt.Fprintf(w, "%+v", req)
		http.Error(w, "incorrect credentials", http.StatusBadRequest)
		return
	}

	session.Values["isLoggedIn"] = true
	session.Values["employeeID"] = req.EmployeeID
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   config.AuthSessionDuration,
		HttpOnly: true,
	}
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "log in successful")
}
