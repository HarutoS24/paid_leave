package model

import (
	"database/sql"
	"time"
)

type Employee struct {
	ID           string
	Name         string
	JoiningDate  time.Time
	IsAdmin      bool
	RegisteredAt time.Time
	DeletedAt    sql.NullTime
}
