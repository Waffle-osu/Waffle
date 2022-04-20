package objects

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type DatabaseUser struct {
	UserID       uint64
	Username     string
	Password     string
	Country      uint16
	Banned       int8
	BannedReason string
	Privileges   int32
	JoinedAt     time.Time
}

type DatabaseUserStats struct {
	UserID      uint64
	Mode        uint8
	RankedScore uint64
	TotalScore  uint64
	Level       float64
	Accuracy    float32
	Playcount   uint64
	CountSSH    uint64
	CountSS     uint64
	CountSH     uint64
	CountS      uint64
	CountA      uint64
	CountB      uint64
	CountC      uint64
	CountD      uint64
	Hit300      uint64
	Hit100      uint64
	Hit50       uint64
	HitMiss     uint64
	HitGeki     uint64
	HitKatu     uint64
}

func UserFromDatabaseById(id uint64) (int8, DatabaseUser) {
	returnUser := DatabaseUser{}

	db, connErr := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/waffle")

	if connErr != nil {
		fmt.Printf("Failed to Fetch User from Database, as a connection could not be successfully established.\n")

		return -2, returnUser
	}

	queryResult, queryErr := db.Query("SELECT user_id, username, password, country, banned, banned_reason, privileges, joined_at FROM waffle.users WHERE user_id = ?", id)

	if queryErr != nil {
		fmt.Printf("Failed to Fetch User from Database, MySQL query failed.\n")

		return -2, returnUser
	}

	if queryResult.Next() {
		scanErr := queryResult.Scan(&returnUser.UserID, &returnUser.Username, &returnUser.Password, &returnUser.Country, &returnUser.Banned, &returnUser.BannedReason, &returnUser.Privileges, &returnUser.JoinedAt)

		if scanErr != nil {
			fmt.Printf("Failed to Scan database results onto User object.\n")

			return -2, returnUser
		}
	}

	//User not found
	return -1, returnUser
}

func UserFromDatabaseByUsername(username string) (int8, DatabaseUser) {
	returnUser := DatabaseUser{}

	db, connErr := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/waffle")

	if connErr != nil {
		fmt.Printf("Failed to Fetch User from Database, as a connection could not be successfully established.\n")

		return -2, returnUser
	}

	queryResult, queryErr := db.Query("SELECT user_id, username, password, country, banned, banned_reason, privileges, joined_at FROM waffle.users WHERE username = ?", username)

	if queryErr != nil {
		fmt.Printf("Failed to Fetch User from Database, MySQL query failed.\n")

		return -2, returnUser
	}

	if queryResult.Next() {
		scanErr := queryResult.Scan(&returnUser.UserID, &returnUser.Username, &returnUser.Password, &returnUser.Country, &returnUser.Banned, &returnUser.BannedReason, &returnUser.Privileges, &returnUser.JoinedAt)

		if scanErr != nil {
			fmt.Printf("Failed to Scan database results onto User object.\n")

			return -2, returnUser
		}
	}

	//User not found
	return -1, returnUser
}
