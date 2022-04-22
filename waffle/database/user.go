package database

import (
	"crypto/md5"
	"encoding/hex"
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
	Rank           uint64
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

	queryResult, queryErr := database.Query("SELECT user_id, username, password, country, banned, banned_reason, privileges, joined_at FROM waffle.users WHERE user_id = ?", id)
	defer queryResult.Close()

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

	queryResult, queryErr := database.Query("SELECT user_id, username, password, country, banned, banned_reason, privileges, joined_at FROM waffle.users WHERE username = ?", username)
	defer queryResult.Close()

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

	queryResult, queryErr := database.Query("SELECT * FROM (SELECT user_id, mode, ROW_NUMBER() OVER (ORDER BY ranked_score DESC) AS 'rank', ranked_score, total_score, user_level, accuracy, playcount, count_ssh, count_ss, count_sh, count_s, count_a, count_b, count_c, count_d, hit300, hit100, hit50, hitMiss, hitGeki, hitKatu, replays_watched FROM waffle.stats WHERE mode = ?) t WHERE user_id = ?", mode, id)
	defer queryResult.Close()

	if queryErr != nil {
		fmt.Printf("Failed to Fetch User Stats from Database, MySQL query failed.\n")

		return -2, returnStats
	}

	if queryResult.Next() {
		scanErr := queryResult.Scan(&returnStats.UserID, &returnStats.Mode, &returnStats.Rank, &returnStats.RankedScore, &returnStats.TotalScore, &returnStats.Level, &returnStats.Accuracy, &returnStats.Playcount, &returnStats.CountSSH, &returnStats.CountSS, &returnStats.CountSH, &returnStats.CountS, &returnStats.CountA, &returnStats.CountB, &returnStats.CountC, &returnStats.CountD, &returnStats.Hit300, &returnStats.Hit100, &returnStats.Hit50, &returnStats.HitMiss, &returnStats.HitGeki, &returnStats.HitKatu, &returnStats.ReplaysWatched)

		if scanErr != nil {
			fmt.Printf("Failed to Scan database results onto UserStats object.\n")

			return -2, returnStats
		}

		return 0, returnStats
	}

	return -1, returnStats
}

func CreateNewUser(username string, rawPassword string) bool {
	duplicateUsernameQuery, duplicateUsernameQueryErr := database.Query("SELECT COUNT(*) FROM waffle.users WHERE username = ?", username)
	defer duplicateUsernameQuery.Close()

	if duplicateUsernameQueryErr != nil {
		fmt.Printf("Failed to create new user, MySQL query failed.\n")

		return false
	}

	if duplicateUsernameQuery.Next() {
		return false
	}

	passwordHashed := md5.Sum([]byte(rawPassword))
	passwordHashedString := hex.EncodeToString(passwordHashed[:])

	var newUserId uint64
	var newUsername string

	insertResult, queryErr := database.Query("INSERT INTO waffle.users (username, password) VALUES (?, ?)", username, passwordHashedString)
	queryResult, queryErr := database.Query("SELECT user_id, username FROM waffle.users WHERE username = ?", username)

	defer insertResult.Close()
	defer queryResult.Close()

	if queryErr != nil {
		fmt.Printf("Failed to create new user, MySQL query failed.\n")

		return false
	}

	if queryResult.Next() {
		scanErr := queryResult.Scan(&newUserId, &newUsername)

		if scanErr != nil {
			return false
		}

		_, statsInsertErr := database.Query("INSERT INTO waffle.stats (user_id, mode) VALUES (?, 0)", newUserId)
		_, statsInsertErr = database.Query("INSERT INTO waffle.stats (user_id, mode) VALUES (?, 1)", newUserId)
		_, statsInsertErr = database.Query("INSERT INTO waffle.stats (user_id, mode) VALUES (?, 2)", newUserId)
		_, statsInsertErr = database.Query("INSERT INTO waffle.stats (user_id, mode) VALUES (?, 3)", newUserId)

		if statsInsertErr != nil {
			fmt.Printf("Failed to create new user, user stats creation failed. MySQL query failed.\n")
			return false
		}
	} else {
		return false
	}

	return true
}
