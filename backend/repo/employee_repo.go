package repo

import (
	"backend/model"
	"database/sql"
	"fmt"
	"os"
	"time"
)

var dbName string

func init() {
	dbName = os.Getenv("MYSQL_DATABASE")
	if dbAuthName == "" {
		panic("従業員管理用db名が設定されていません")
	}
}

type AddEmployeeParams struct {
	EmployeeID   string
	EmployeeName string
	JoiningDate  time.Time
	IsAdmin      bool
	RegisteredAt time.Time
	DeletedAt    sql.NullTime
}
func GetEmployeeByID(db *sql.DB, employeeID string) (model.Employee, error) {
	query := fmt.Sprintf("SELECT id, name, is_admin, joining_date, registered_at, deleted_at FROM %s.employees_tbl WHERE id=?;", dbName)
	row := db.QueryRow(query, employeeID)

	var employee model.Employee
	err := row.Scan(&employee.ID, &employee.Name, &employee.IsAdmin, &employee.JoiningDate, &employee.RegisteredAt, &employee.DeletedAt)
	// レコードが見つからない場合はエラーとしない
	if err == sql.ErrNoRows {
		return model.Employee{}, nil
	} else if err != nil {
		return model.Employee{}, err
	}
	return employee, nil
}
