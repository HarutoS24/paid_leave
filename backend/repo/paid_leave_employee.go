package repo

import (
	"backend/model"
	"database/sql"
	"fmt"
	"time"
)

type AddLeaveEmployeeParams struct {
	EmployeeID   string
	VacationDate time.Time
	StartAtHour  int
	Duration     int
	GivenAt      time.Time
	RegisteredAt time.Time
}

func GetRegisteredLeaveListByGivenAt(db *sql.DB, employeeID string, givenAt time.Time) ([]model.LeaveEmployee, error) {
	query := fmt.Sprintf("SELECT vacation_date, start_at_hour, duration FROM %s.paid_leave_employee_tbl WHERE employee_id=? AND given_at=? AND deleted_at is NULL;", dbName)
	fmt.Println(employeeID)
	fmt.Println(givenAt.Format("2006-01-02"))
	rows, err := db.Query(query, employeeID, givenAt.Format("2006-01-02"))
	var leaveList []model.LeaveEmployee
	if err == sql.ErrNoRows {
		return []model.LeaveEmployee{}, nil
	} else if err != nil {
		return []model.LeaveEmployee{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var leaveEmployee model.LeaveEmployee
		if err := rows.Scan(&leaveEmployee.VacationDate, &leaveEmployee.StartAtHour, &leaveEmployee.Duration); err != nil {
			return []model.LeaveEmployee{}, err
		}
		leaveList = append(leaveList, leaveEmployee)
	}

	return leaveList, nil
}

func AddLeaveEmployee(tx *sql.Tx, leaveEmployee []AddLeaveEmployeeParams) error {
	query := fmt.Sprintf("INSERT INTO %s.paid_leave_employee_tbl (employee_id, vacation_date, start_at_hour, duration, given_at, registered_at) VALUES (?, ?, ?, ?, ?, ?)", dbName)
	today := time.Now()

	for _, row := range leaveEmployee {
		_, err := tx.Exec(query, row.EmployeeID, row.VacationDate, row.StartAtHour, row.Duration, row.GivenAt, today)
		if err != nil {
			return err
		}
	}

	return nil
}
