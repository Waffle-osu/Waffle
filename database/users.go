package database

import (
	"Waffle/helpers"
	"crypto/md5"
	"encoding/hex"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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

// UserFromDatabaseById retrieves a users information given a user id
func UserFromDatabaseById(id uint64) (int8, User) {
	returnUser := User{}

	queryResult, queryErr := Database.Query("SELECT user_id, username, password, country, banned, banned_reason, privileges, joined_at FROM waffle.users WHERE user_id = ?", id)

	if queryErr != nil {
		helpers.Logger.Printf("[Database] Failed to Fetch User from Database, MySQL query failed.\n")

		if queryResult != nil {
			queryResult.Close()
		}

		return -2, returnUser
	}

	if queryResult.Next() {
		scanErr := queryResult.Scan(&returnUser.UserID, &returnUser.Username, &returnUser.Password, &returnUser.Country, &returnUser.Banned, &returnUser.BannedReason, &returnUser.Privileges, &returnUser.JoinedAt)

		queryResult.Close()

		if scanErr != nil {
			helpers.Logger.Printf("[Database] Failed to Scan Database results onto User object.\n")

			return -2, returnUser
		}

		return 0, returnUser
	}

	queryResult.Close()
	//User not found
	return -1, returnUser
}

// UserFromDatabaseByUsername retrieves a users information given a username
func UserFromDatabaseByUsername(username string) (int8, User) {
	returnUser := User{}

	queryResult, queryErr := Database.Query("SELECT user_id, username, password, country, banned, banned_reason, privileges, joined_at FROM waffle.users WHERE username = ?", username)
	defer queryResult.Close()

	if queryErr != nil {
		helpers.Logger.Printf("[Database] Failed to Fetch User from Database, MySQL query failed.\n")

		if queryResult != nil {
			queryResult.Close()
		}

		return -2, returnUser
	}

	//If there is a result
	if queryResult.Next() {
		scanErr := queryResult.Scan(&returnUser.UserID, &returnUser.Username, &returnUser.Password, &returnUser.Country, &returnUser.Banned, &returnUser.BannedReason, &returnUser.Privileges, &returnUser.JoinedAt)

		queryResult.Close()

		if scanErr != nil {
			helpers.Logger.Printf("[Database] Failed to Scan Database results onto User object.\n")

			return -2, returnUser
		}

		return 0, returnUser
	}

	queryResult.Close()
	//User not found
	return -1, returnUser
}

// CreateNewUser creates a new user given a username and a password
func CreateNewUser(username string, rawPassword string) bool {
	duplicateUsernameQuery, duplicateUsernameQueryErr := Database.Query("SELECT COUNT(*) FROM waffle.users WHERE username = ?", username)

	if duplicateUsernameQueryErr != nil {
		if duplicateUsernameQuery != nil {
			duplicateUsernameQuery.Close()
		}

		helpers.Logger.Printf("[Database] Failed to create new user, MySQL query failed.\n")

		return false
	}

	if duplicateUsernameQuery.Next() {
		var count uint64

		scanErr := duplicateUsernameQuery.Scan(&count)

		duplicateUsernameQuery.Close()

		if count != 0 || scanErr != nil {
			return false
		}
	}

	passwordHashed := md5.Sum([]byte(rawPassword))
	passwordHashedString := hex.EncodeToString(passwordHashed[:])
	bcryptPassword, bcryptErr := bcrypt.GenerateFromPassword([]byte(passwordHashedString), bcrypt.DefaultCost)

	if bcryptErr != nil {
		return false
	}

	var newUserId uint64
	var newUsername string

	insertResult, queryErrInsert := Database.Query("INSERT INTO waffle.users (username, password) VALUES (?, ?)", username, bcryptPassword)
	queryResult, queryErrGet := Database.Query("SELECT user_id, username FROM waffle.users WHERE username = ?", username)

	insertResult.Close()

	if queryErrInsert != nil || queryErrGet != nil {
		helpers.Logger.Printf("[Database] Failed to create new user, MySQL query failed.\n")

		return false
	}

	if queryResult.Next() {
		scanErr := queryResult.Scan(&newUserId, &newUsername)

		queryResult.Close()

		if scanErr != nil {
			return false
		}

		osuStatsInsert, statsInsertErrOsu := Database.Query("INSERT INTO waffle.stats (user_id, mode) VALUES (?, 0)", newUserId)
		taikoStatsInsert, statsInsertErrTaiko := Database.Query("INSERT INTO waffle.stats (user_id, mode) VALUES (?, 1)", newUserId)
		catchStatsInsert, statsInsertErrCatch := Database.Query("INSERT INTO waffle.stats (user_id, mode) VALUES (?, 2)", newUserId)
		maniaStatsInsert, statsInsertErrMania := Database.Query("INSERT INTO waffle.stats (user_id, mode) VALUES (?, 3)", newUserId)

		osuStatsInsert.Close()
		taikoStatsInsert.Close()
		catchStatsInsert.Close()
		maniaStatsInsert.Close()

		if statsInsertErrOsu != nil || statsInsertErrTaiko != nil || statsInsertErrCatch != nil || statsInsertErrMania != nil {
			helpers.Logger.Printf("[Database] Failed to create new user, user stats creation failed. MySQL query failed.\n")
			return false
		}
	} else {
		return false
	}

	return true
}

func AuthenticateUser(username string, password string) (userId int32, authSuccess bool) {
	query, queryErr := Database.Query("SELECT user_id, username, password FROM waffle.users WHERE username = ?", username)

	var scanUsername, scanPassword string
	var scanUserId int32

	if queryErr != nil {
		if query != nil {
			query.Close()
		}

		return -2, false
	}

	if query.Next() {
		scanErr := query.Scan(&scanUserId, &scanUsername, &scanPassword)

		query.Close()

		if scanErr != nil {
			return -2, false
		}

		if bcrypt.CompareHashAndPassword([]byte(scanPassword), []byte(password)) == nil {
			return scanUserId, true
		} else {
			return scanUserId, false
		}
	} else {
		return -1, false
	}
}
