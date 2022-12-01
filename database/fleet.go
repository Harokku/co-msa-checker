package database

import (
	"co-msa-checker/utils"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
)

type Fleet struct {
	Id        string `json:"id"`
	Radiocode string `json:"radiocode"`
	Plate     string `json:"plate"`
	Note      string `json:"note"`
}

// -------------------------
// Db Call
// -------------------------

// FleetGetAll retrieve all fleet data
func FleetGetAll() ([]Fleet, error) {
	var (
		err          error
		rows         *sql.Rows
		sqlStatement string
		res          []Fleet
	)

	sqlStatement = `select id, radiocode, plate, COALESCE(note, '') from msa`

	rows, err = DbConnection.Query(sqlStatement)
	if err != nil {
		return nil, fmt.Errorf("error retrieving fleet from db:\tFleetGetAll\t%v", err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			utils.LogDeferError("FleetGetAll", err)
		}
	}(rows)

	for rows.Next() {
		var c Fleet
		err = rows.Scan(&c.Id, &c.Radiocode, &c.Plate, &c.Note)
		if err != nil {
			return nil, fmt.Errorf("error scanning row:\tFleetGetAll\t%v", err)
		}
		res = append(res, c)
	}

	return res, nil
}

// FleetGetById retrieve specified record based on id
func FleetGetById(id string) (Fleet, error) {
	var (
		err          error
		row          *sql.Row
		sqlStatement string
		res          Fleet
	)

	sqlStatement = `select id, radiocode, plate, COALESCE(note, '') from msa where id=$1`

	row = DbConnection.QueryRow(sqlStatement, id)

	switch err = row.Scan(&res.Id, &res.Radiocode, &res.Plate, &res.Note); err {
	case nil:
		// no error return response
		return res, nil
	default:
		// error, pass it to caller
		utils.Err(err)
		return Fleet{}, err
	}
}

func BulkCreateFleet(newFleet []Fleet) error {
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
	sqlStatement, err = txn.Prepare(pq.CopyIn("msa", "radiocode", "plate", "note"))

	//Exec insert for every passed document
	for _, msa := range newFleet {
		_, err = sqlStatement.Exec(msa.Radiocode, msa.Plate, msa.Note)
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
