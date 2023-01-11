package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("sqlite3", "file:posts.sqlite?cache=shared")
	if err != nil {
		panic(err)
	}
}

func Close() {
	_ = DB.Close()
}
