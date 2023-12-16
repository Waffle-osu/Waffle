package database

import (
	"Waffle/config"
	"Waffle/helpers"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Database *sql.DB

// Initialize initializes the MySQL Database things
func Initialize() {
	if Database != nil {
		return
	}

	for {
		db, _ := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", config.MySqlUsername, config.MySqlPassword, config.MySqlLocation, config.MySqlDatabase))

		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)

		pingErr := db.Ping()

		if pingErr != nil {
			helpers.Logger.Printf("[Database] MySQL Connection could not be established... Retrying...\n")

			time.Sleep(time.Second * 5)
			continue
		} else {
			Database = db

			break
		}
	}
}

func Deinitialize() {
	Database.Close()
}
