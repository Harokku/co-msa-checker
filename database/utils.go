package database

import (
	"fmt"
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
