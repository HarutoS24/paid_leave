package handler

import (
	"backend/config"
	"backend/dto"
	"backend/model"
	"backend/service"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func AddPaidLeaveHandler(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value(config.DBContextKey).(*sql.DB)
	employee := r.Context().Value(config.LoginUserContextKey).(model.Employee)

	var req dto.AddPaidLeaveRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Duration == nil || req.StartAtHour == nil || req.VacationDate == nil {
		http.Error(w, "空欄があります", http.StatusBadRequest)
		return
	}
	if *req.StartAtHour < 9 || *req.StartAtHour > 18 {
		http.Error(w, "有給開始時間の指定が不正です", http.StatusBadRequest)
		return
	}
	if *req.StartAtHour+*req.Duration > 18 || *req.Duration <= 0 || *req.Duration%2 != 0 {
		http.Error(w, "有給時間の指定が不正です", http.StatusBadRequest)
		return
	}
	today := time.Now()
	if !req.VacationDate.After(today) {
		http.Error(w, "今日以降の日付について申請してください", http.StatusBadRequest)
		return
	}

	var params = service.AddPaidLeaveParams{
		EmployeeID:   employee.ID,
		VacationDate: *req.VacationDate,
		StartAtHour:  *req.StartAtHour,
		Duration:     *req.Duration,
	}
	err := service.AddPaidLeave(db, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "success")
}

func GetInfoHandler(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value(config.DBContextKey).(*sql.DB)
	employee := r.Context().Value(config.LoginUserContextKey).(model.Employee)
	info, err := service.CalculateSumPaidLeaveInfo(db, employee.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var res dto.GetInfoResponse
	res.Total = info.TotalCount
	res.Used = info.Used
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
