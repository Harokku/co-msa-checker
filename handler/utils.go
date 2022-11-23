package handler

import "co-msa-checker/database"

// checkUser Check if user exist in db and retrieve username
func checkUser(id string) (database.User, error) {
	var (
		user database.User
		err  error
	)

	user, err = database.GetUserDataById(id)
	if err != nil {
		return database.User{}, err
	}

	return user, nil
}
