package database

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"time"
)

type Token struct {
	TokenHash    string
	CreationDate time.Time
}

func TokensCreateNewToken(user User) string {
	newToken := Token{
		CreationDate: time.Now(),
	}

	hashInput := fmt.Sprintf("wa%sff%sle%dto%dke%sn", user.Username, user.Password, newToken.CreationDate.UnixMilli(), user.UserID, user.JoinedAt)
	hashedInput := sha512.Sum512([]byte(hashInput))
	stringHashedInput := hex.EncodeToString(hashedInput[:])

	tokenInsertQuery, tokenInsertQueryErr := Database.Query("INSERT INTO site_tokens (token_hash) VALUES (?)", stringHashedInput)

	if tokenInsertQueryErr != nil {
		return ""
	}

	tokenInsertQuery.Close()

	return stringHashedInput
}
