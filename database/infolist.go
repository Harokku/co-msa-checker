package database

import (
	"co-msa-checker/utils"
	"database/sql"
	"fmt"
	"time"
)

type Info struct {
	Id        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Note      string    `json:"note"`
	Status    bool      `json:"status"`
	Operator  string    `json:"operator"`
	Priority  string    `json:"priority"`
	MsaId     string    `json:"msa_id"`
}

type Update struct {
	Id         string    `json:"id"`
	Timestamp  time.Time `json:"timestamp"`
	Note       string    `json:"note"`
	Operator   string    `json:"operator"`
	Deprecated bool      `json:"deprecated"`
	InfoId     string    `json:"info_id"`
}

type FullInfo struct {
	Info    Info     `json:"info"`
	Updates []Update `json:"updates"`
}

// -------------------------
// Db Call
// -------------------------

// -------------------------
// Info
// -------------------------

// InfoGetAllByMsaId retrieve all base info based on MSA id without updates
func InfoGetAllByMsaId(id string) ([]Info, error) {
	var (
		err          error
		rows         *sql.Rows
		sqlStatement string
		res          []Info
	)

	sqlStatement = `select id,timestamp,note,status,operator,priority from infolist where msa_id=$1 and status=false`

	rows, err = DbConnection.Query(sqlStatement, id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving info list from db:\tInfoGetAllMsaId\t%v", err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			utils.LogDeferError("InfoGetAllByMsaId", err)
		}
	}(rows)

	for rows.Next() {
		var c Info
		err = rows.Scan(&c.Id, &c.Timestamp, &c.Note, &c.Status, &c.Operator, &c.Priority)
		if err != nil {
			return []Info{}, fmt.Errorf("error scanning row:\tInfoGetAllByMsaId\t%w", err)
		}
		res = append(res, c)
	}

	return res, nil
}

// -------------------------
// Updates
// -------------------------

func InfoGetAllUpdatesByInfoId(id string) ([]Update, error) {
	var (
		err          error
		rows         *sql.Rows
		sqlStatement string
		res          []Update
	)

	sqlStatement = `select id, timestamp, note,operator,deprecated from updates where info_id=$1`

	rows, err = DbConnection.Query(sqlStatement, id)
	if err != nil {
		return []Update{}, fmt.Errorf("error retrieving updates from db:\tInfoGetAllUpdatesByInfoId\t%w", err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			utils.LogDeferError("InfoGetAllUpdatesByInfoId", err)
		}
	}(rows)

	for rows.Next() {
		var c Update
		err = rows.Scan(&c.Id, &c.Timestamp, &c.Note, &c.Operator, &c.Deprecated)
		if err != nil {
			return []Update{}, fmt.Errorf("error scanning row:\tInfoGetAllUpdatesByInfoId\t%w", err)
		}
		res = append(res, c)
	}

	return res, nil
}

func NewUpdate(data Update) (Update, error) {
	var (
		err          error
		sqlStatement string
		res          Update
	)

	res = data

	sqlStatement = `
		INSERT INTO updates (info_id,note,operator)
		VALUES ($1,$2,$3)
		RETURNING id,timestamp
`

	err = DbConnection.QueryRow(sqlStatement, data.InfoId, data.Note, data.Operator).Scan(&res.Id, &res.Timestamp)
	if err != nil {
		return Update{}, fmt.Errorf("error inserting record:\taddUpdate\t%w", err)
	}

	return res, nil
}
