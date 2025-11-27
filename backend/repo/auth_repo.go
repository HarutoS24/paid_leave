package repo

import (
	"database/sql"
	"fmt"
	"os"
)

type AddAuthParams struct {
	EmployeeID string
	Hash       string
}

var dbAuthName string

func init() {
	dbAuthName = os.Getenv("MYSQL_AUTH_DATABASE")
	if dbAuthName == "" {
		panic("認証用db名が設定されていません")
	}
}

func AddAuth(tx *sql.Tx, auth AddAuthParams) error {
	query := fmt.Sprintf("INSERT INTO %s.auth_tbl (employee_id, hash) VALUES (?,?);", dbAuthName)
	_, err := tx.Exec(query, auth.EmployeeID, auth.Hash)
	return err
}

func GetHash(db *sql.DB, employeeID string) (string, error) {
	var passwordHash string
	query := fmt.Sprintf("SELECT hash FROM %s.auth_tbl WHERE employee_id=?", dbAuthName)
	row := db.QueryRow(query, employeeID)
	err := row.Scan(&passwordHash)
	// レコードが見つからない場合はエラーとして扱わない
	if err == sql.ErrNoRows {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return passwordHash, nil
}
