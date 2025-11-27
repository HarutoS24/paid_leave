package service

import (
	"backend/repo"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"fmt"

	"github.com/gorilla/sessions"
)

func ValidateUserByPassword(db *sql.DB, employeeID string, password string) (bool, error) {
	hash, err := repo.GetHash(db, employeeID)
	if err != nil {
		return false, err
	}
	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return false, err
	}
	hashedInput := sha256.Sum256([]byte(password))

	if subtle.ConstantTimeCompare(hashedInput[:], hashBytes) == 1 {
		return true, nil
	}
	return false, nil
}

func IsLoggedIn(session *sessions.Session) (bool, error) {
	isLoggedInRaw, ok := session.Values["isLoggedIn"]
	if !ok {
		return false, nil
	}
	isLoggedIn, ok := isLoggedInRaw.(bool)
	if !ok {
		return false, fmt.Errorf("invalid value found in the cookie")
	}
	return isLoggedIn, nil
}
