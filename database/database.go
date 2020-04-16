package database

import (
	"database/sql"
	"io/ioutil"
)

func InitDB(db *sql.DB) error {
	bFile, err := ioutil.ReadFile("database/init.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(bFile))
	if err != nil {
		return err
	}

	return nil
}