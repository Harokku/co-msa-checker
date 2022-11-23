package database

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
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

	// new method with only hashing
	h := sha256.New()
	h.Write([]byte(pw))
	passwordHash = fmt.Sprintf("%x", h.Sum(nil))

	if err != nil {
		return User{}, fmt.Errorf("error calculating password hash")
	}

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
