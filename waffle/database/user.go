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

func CreateNewUser(username string, rawPassword string) bool {
	duplicateUsernameQuery, duplicateUsernameQueryErr := database.Query("SELECT COUNT(*) FROM waffle.users WHERE username = ?", username)
	defer duplicateUsernameQuery.Close()

	if duplicateUsernameQueryErr != nil {
		fmt.Printf("Failed to create new user, MySQL query failed.\n")

		return false
	}

	if duplicateUsernameQuery.Next() {
		var count uint64

		duplicateUsernameQuery.Scan(&count)

		if count != 0 {
			return false
		}
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

		osuStatsInsert, statsInsertErr := database.Query("INSERT INTO waffle.stats (user_id, mode) VALUES (?, 0)", newUserId)
		taikoStatsInsert, statsInsertErr := database.Query("INSERT INTO waffle.stats (user_id, mode) VALUES (?, 1)", newUserId)
		catchStatsInsert, statsInsertErr := database.Query("INSERT INTO waffle.stats (user_id, mode) VALUES (?, 2)", newUserId)
		maniaStatsInsert, statsInsertErr := database.Query("INSERT INTO waffle.stats (user_id, mode) VALUES (?, 3)", newUserId)

		osuStatsInsert.Close()
		taikoStatsInsert.Close()
		catchStatsInsert.Close()
		maniaStatsInsert.Close()

		if statsInsertErr != nil {
			fmt.Printf("Failed to create new user, user stats creation failed. MySQL query failed.\n")
			return false
		}
	} else {
		return false
	}

	return true
}
