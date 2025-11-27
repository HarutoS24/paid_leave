package model

import (
	"database/sql"
	"time"
)

type LeaveEmployee struct {
	Id           int
	EmployeeId   string
	VacationDate time.Time
	StartAtHour  int
	Duration     int
	GivenAt      time.Time
	RegisteredAt time.Time
	DeletedAt    sql.NullTime
}
