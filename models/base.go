package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB

func ShowErrLog(err error) {
	log.Println(err)
}

func InitDb(dsn string) (*sql.DB, error) {
	var err error

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return db, err
	}

	if err = db.Ping(); err != nil {
		return db, err
	}

	return db, err
}
