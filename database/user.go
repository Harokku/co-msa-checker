package database

import (
	"co-msa-checker/utils"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
)

type User struct {
	Id          string `json:"id,omitempty"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	ManagerRole bool   `json:"manager_role,omitempty"`
	Token       string `json:"token,omitempty"`
}

// -------------------------
// Db Call
// -------------------------

func GetUserDataById(id string) (User, error) {
	var (
		err          error
		row          *sql.Row
		sqlStatement string
		res          User
	)

	sqlStatement = `
			SELECT username, manager_role
			FROM users
			WHERE id = $1
`
	row = DbConnection.QueryRow(sqlStatement, id)
	switch err = row.Scan(&res.Username, &res.ManagerRole); err {
	case sql.ErrNoRows:
		return User{}, fmt.Errorf("error retrieving user data from db:\tGetUserData\t%v", err)
	case nil:
		return res, nil
	default:
		return User{}, fmt.Errorf("error retrieving user data from db:\tGetUserData\t%v", err)
	}
}

func GetUserDataByPassword(pw string) (User, error) {
	var (
		err          error
		row          *sql.Row
		sqlStatement string
		passwordHash string
		res          User
	)

	sqlStatement = `
			SELECT id, username, manager_role
			FROM users
			WHERE password = $1
`
	// old method with salting
	//passwordHash, err = hashPassword(pw)
	//if err != nil {
	//	return User{}, fmt.Errorf("error calculating password hash")
	//}

	// new method with only hashing
	passwordHash = utils.Hash256(pw)

	row = DbConnection.QueryRow(sqlStatement, passwordHash)
	switch err = row.Scan(&res.Id, &res.Username, &res.ManagerRole); err {
	case sql.ErrNoRows:
		return User{}, fmt.Errorf("error retrieving user data from db:\tGetUserDataByPassword\t%v", err)
	case nil:
		return res, nil
	default:
		return User{}, fmt.Errorf("error retrieving user data from db:\tGetUserDataByPassword\t%v", err)
	}
}

func BulkCreate(newUsers []User) error {
	var (
		err          error
		sqlStatement *sql.Stmt //Prepared sql statement
		txn          *sql.Tx   //DB transaction
	)

	//Begin new transaction
	txn, err = DbConnection.Begin()
	if err != nil {
		return err
	}

	//Prepare insert statement
	sqlStatement, err = txn.Prepare(pq.CopyIn("users", "username", "password", "manager_role"))

	//Exec insert for every passed document
	for _, user := range newUsers {
		_, err = sqlStatement.Exec(user.Username, user.Password, user.ManagerRole)
		if err != nil {
			return err
		}
	}

	//Flush actual data
	_, err = sqlStatement.Exec()
	if err != nil {
		return err
	}

	err = sqlStatement.Close()
	if err != nil {
		return err
	}

	//Execute transaction
	err = txn.Commit()
	if err != nil {
		return err
	}

	return nil
}
