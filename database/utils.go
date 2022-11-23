package database

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// TruncateTable Truncate (clean) actual table
func TruncateTable(t string) error {
	var (
		err          error
		sqlStatement string
	)

	sqlStatement = `TRUNCATE TABLE $1`

	_, err = DbConnection.Exec(sqlStatement, t)
	if err != nil {
		return fmt.Errorf("error truncating %s table", t)
	}

	return nil
}

// -------------------------
// bcrypt standard functions
// -------------------------

// HashPassword hash given password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash check validity of a give password
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
