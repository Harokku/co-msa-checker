package database

import (
	"co-msa-checker/utils"
	"fmt"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
	"io"
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

type Xls struct {
}

func (x Xls) BuildUsers(inputStream io.Reader) ([]User, error) {
	var (
		res []User
	)

	f, err := excelize.OpenReader(inputStream)
	if err != nil {
		return nil, err
	}

	usersSheet, err := f.GetRows("utenti")
	if err != nil {
		return nil, err
	}

	for i, user := range usersSheet {
		// Skip 1st row, contain only column header for xlsx human readability
		if i > 0 {
			var u User
			var randomPw string
			var hashedPw string

			u.Username = user[0]

			// Generate new random password and hash it
			randomPw, _ = utils.NewPw(8) // TODO: Check for error
			hashedPw = utils.Hash256(randomPw)
			u.Password = hashedPw

			// Add plain password to xlsx row
			f.SetCellValue("utenti", fmt.Sprintf("C%d", i+1), randomPw)

			// Cast xlsx text to boolean
			if user[1] == "TRUE" {
				u.ManagerRole = true
			}

			res = append(res, u)
		}
	}

	// TODO: implement in memory response to avoid temporary disk access
	f.SaveAs("Utenti.xlsx")

	return res, nil
}

func (x Xls) BuildFleet(inputStream io.Reader) ([]Fleet, error) {
	var (
		res []Fleet
	)

	f, err := excelize.OpenReader(inputStream)
	if err != nil {
		return nil, err
	}

	fleetSheet, err := f.GetRows("msa")
	if err != nil {
		return nil, err
	}

	for i, msa := range fleetSheet {
		// Skip 1st row, contain only column header for xlsx human readability
		if i > 0 {
			var m Fleet

			m.Radiocode = msa[0]
			m.Plate = msa[1]
			// Since "note" field can be blank, check if slice item is populated
			if len(msa) > 2 {
				m.Note = msa[2]
			}

			res = append(res, m)
		}
	}

	return res, nil
}
