package database

import (
	"github.com/jackc/pgx"
	"io/ioutil"
)

func InitDB(db *pgx.ConnPool) error {
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
