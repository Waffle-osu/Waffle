package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var database *sql.DB

func Initialize() {
	db, connErr := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/waffle")

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if connErr != nil {
		fmt.Printf("MySQL Connection could not be established...\n")

		return
	}

	database = db
}

func Deinitialize() {
	database.Close()
}
