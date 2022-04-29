package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var database *sql.DB

// Initialize initializes the MySQL database things
func Initialize(username string, password string, location string, dbDatabase string) {
	db, connErr := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, location, dbDatabase))

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
