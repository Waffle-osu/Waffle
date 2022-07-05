package api

import (
	"Waffle/database"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiLoginResponse struct {
	WaffleUsername string `json:"WaffleUsername"`
	WaffleToken    string `json:"WaffleToken"`
	WaffleUserId   int64  `json:"WaffleUserId"`
}

func ApiHandleWaffleLogin(ctx *gin.Context) {
	formUsername := ctx.PostForm("username")
	formPassword := ctx.PostForm("password")

	ctx.Header("Access-Control-Allow-Origin", "127.0.0.1")
	ctx.Header("Access-Control-Allow-Origin", "http://localhost:3000")

	passwordHashed := md5.Sum([]byte(formPassword))
	passwordHashedString := hex.EncodeToString(passwordHashed[:])

	userId, authSuccess := database.AuthenticateUser(formUsername, passwordHashedString)

	loginResponse := ApiLoginResponse{}

	if userId == -2 {
		ctx.String(http.StatusInternalServerError, "")
		return
	}

	if userId == -1 {
		loginResponse.WaffleToken = ""
		loginResponse.WaffleUsername = ""
		loginResponse.WaffleUserId = -1

		data, marshalErr := json.Marshal(loginResponse)

		if marshalErr != nil {
			ctx.String(http.StatusInternalServerError, "")
			return
		}

		ctx.Data(http.StatusOK, "waffle/blob", data)
		return
	}

	userQueryResult, user := database.UserFromDatabaseById(uint64(userId))

	if userQueryResult == -2 {
		ctx.String(http.StatusInternalServerError, "")
		return
	}

	if userQueryResult == -1 {
		loginResponse.WaffleToken = ""
		loginResponse.WaffleUsername = ""
		loginResponse.WaffleUserId = -1

		data, marshalErr := json.Marshal(loginResponse)

		if marshalErr != nil {
			ctx.String(http.StatusInternalServerError, "")
			return
		}

		ctx.Data(http.StatusOK, "waffle/blob", data)
		return
	}

	if authSuccess {
		loginResponse.WaffleToken = database.TokensCreateNewToken(user)
		loginResponse.WaffleUserId = int64(userId)
		loginResponse.WaffleUsername = user.Username

		data, marshalErr := json.Marshal(loginResponse)

		if marshalErr != nil {
			ctx.String(http.StatusInternalServerError, "")
			return
		}

		ctx.Data(http.StatusOK, "waffle/blob", data)
	} else {
		loginResponse.WaffleToken = ""
		loginResponse.WaffleUsername = ""
		loginResponse.WaffleUserId = -1

		data, marshalErr := json.Marshal(loginResponse)

		if marshalErr != nil {
			ctx.String(http.StatusInternalServerError, "")
			return
		}

		ctx.Data(http.StatusOK, "waffle/blob", data)
		return
	}
}
