package dto

import "time"

type AddPaidLeaveRequest struct {
	VacationDate *time.Time `json:"vacationDate"`
	StartAtHour  *int       `json:"startAtHour"`
	Duration     *int       `json:"duration"`
}

type AddPaidLeaveResponse = BaseReponse

type GetInfoResponse struct {
	Total int     `json:"total"`
	Used  float64 `json:"used"`
}
