package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

// DB активное соединение с БД
var DB *sql.DB

// Setup
func Setup(dbFileName string) {
	if dbFileName == "" {
		dbFileName = "posts.sqlite"
	}

	wordDir, err := os.Getwd()
	if err != nil {
		panic("не удалось получить рабочую директорию")
	}

	IsNewDb := false

	_, err = os.Stat(fmt.Sprintf("%s/%s", wordDir, dbFileName))
	if os.IsNotExist(err) {
		err = os.WriteFile(dbFileName, []byte{}, 0764)
		IsNewDb = true
		if err != nil {
			panic("не удалось создать файл бд")
		}
	} else if err != nil {
		panic("не удалось получить информацию о бд")
	}

	DB, err = sql.Open("sqlite3", fmt.Sprintf("file:%s?", dbFileName))
	if err != nil {
		panic(err)
	}

	if IsNewDb {
		_, err = DB.Exec(sqlCreateTable())
		if err != nil {
			panic("не удалось создать таблицу для записей в бд")
		}
	}
}

// Close
func Close() {
	_ = DB.Close()
}

// sqlCreateTable
func sqlCreateTable() string {
	return `create table posts
(
    id         INTEGER
        constraint posts_pk
            primary key autoincrement,
    created_at TEXT not null,
    text       TEXT
);`
}
