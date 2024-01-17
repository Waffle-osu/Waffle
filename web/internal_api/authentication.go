package internal_api

import (
	"Waffle/config"
	"Waffle/database"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type DoAuthRequest struct {
	Username string
	Password string
}

func InternalDoAuth(ctx *gin.Context) {
	data := DoAuthRequest{}
	err := ctx.BindJSON(&data)

	if err != nil {
		ctx.String(500, "")
	}

	userId, authQueryErr := database.AuthenticateUser(data.Username, data.Password)

	if authQueryErr {
		ctx.String(500, "")
	}

	queryErr, userData := database.UserFromDatabaseById(uint64(userId))

	if queryErr != 0 {
		ctx.String(500, "")
	}

	tokenString := fmt.Sprintf(config.TokenFormatString, userData.Password, userData.JoinedAt, userId, time.Now().Unix(), data.Password)
	tokenHashed := sha256.Sum256([]byte(tokenString))
	tokenBytes := tokenHashed[:]
	returnToken := hex.EncodeToString(tokenBytes)

	ctx.String(200, returnToken)
}
