package dto

import (
	"fmt"
	"time"
)

type DateOnly struct {
	time.Time
}

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := string(b)
	if len(s) < 2 {
		return fmt.Errorf("invalid date string")
	}
	s = s[1 : len(s)-1]
	t, err := time.Parse("2006/01/02", s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

type AddEmployeeRequest struct {
	EmployeeID   *string   `json:"employeeId"`
	EmployeeName *string   `json:"employeeName"`
	Password     *string   `json:"password"`
	IsAdmin      *bool     `json:"isAdmin"`
	JoiningDate  *DateOnly `json:"joiningDate"`
}

type AddEmployeeResponse = BaseReponse
