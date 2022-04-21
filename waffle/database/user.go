package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	UserID       uint64
	Username     string
	Password     string
	Country      uint16
	Banned       int8
	BannedReason string
	Privileges   int32
	JoinedAt     string
}

type UserStats struct {
	UserID         uint64
	Mode           uint8
	RankedScore    uint64
	TotalScore     uint64
	Level          float64
	Accuracy       float32
	Playcount      uint64
	CountSSH       uint64
	CountSS        uint64
	CountSH        uint64
	CountS         uint64
	CountA         uint64
	CountB         uint64
	CountC         uint64
	CountD         uint64
	Hit300         uint64
	Hit100         uint64
	Hit50          uint64
	HitMiss        uint64
	HitGeki        uint64
	HitKatu        uint64
	ReplaysWatched uint64
}

func UserFromDatabaseById(id uint64) (int8, User) {
	returnUser := User{}

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

		return 0, returnUser
	}

	//User not found
	return -1, returnUser
}

func UserFromDatabaseByUsername(username string) (int8, User) {
	returnUser := User{}

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

	//If there is a result
	if queryResult.Next() {
		scanErr := queryResult.Scan(&returnUser.UserID, &returnUser.Username, &returnUser.Password, &returnUser.Country, &returnUser.Banned, &returnUser.BannedReason, &returnUser.Privileges, &returnUser.JoinedAt)

		if scanErr != nil {
			fmt.Printf("Failed to Scan database results onto User object.\n")

			return -2, returnUser
		}

		return 0, returnUser
	}

	//User not found
	return -1, returnUser
}

func UserStatsFromDatabase(id uint64, mode int8) (int8, UserStats) {
	returnStats := UserStats{}

	db, connErr := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/waffle")

	if connErr != nil {
		fmt.Printf("Failed to Fetch User Stats from Database, as a connection could not be successfully established.\n")

		return -2, returnStats
	}

	queryResult, queryErr := db.Query("SELECT user_id, mode, ranked_score, total_score, user_level, accuracy, playcount, count_ssh, count_ss, count_sh, count_s, count_a, count_b, count_c, count_d, hit300, hit100, hit50, hitMiss, hitGeki, hitKatu, replays_watched FROM waffle.stats WHERE user_id = ? AND mode = ?", id, mode)

	if queryErr != nil {
		fmt.Printf("Failed to Fetch User Stats from Database, MySQL query failed.\n")

		return -2, returnStats
	}

	if queryResult.Next() {
		scanErr := queryResult.Scan(&returnStats.UserID, &returnStats.Mode, &returnStats.RankedScore, &returnStats.TotalScore, &returnStats.Level, &returnStats.Accuracy, &returnStats.Playcount, &returnStats.CountSSH, &returnStats.CountSS, &returnStats.CountSH, &returnStats.CountS, &returnStats.CountA, &returnStats.CountB, &returnStats.CountC, &returnStats.CountD, &returnStats.Hit300, &returnStats.Hit100, &returnStats.Hit50, &returnStats.HitMiss, &returnStats.HitGeki, &returnStats.HitKatu, &returnStats.ReplaysWatched)

		if scanErr != nil {
			fmt.Printf("Failed to Scan database results onto UserStats object.\n")

			return -2, returnStats
		}

		return 0, returnStats
	}

	return -1, returnStats
}
