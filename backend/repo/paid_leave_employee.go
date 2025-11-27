package repo

import (
	"backend/model"
	"database/sql"
	"fmt"
	"time"
)

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
