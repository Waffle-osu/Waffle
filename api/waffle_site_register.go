package api

import (
	"Waffle/database"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiRegisterResponse struct {
	WaffleRegisterStatus string `json:"WaffleRegisterStatus"`
	WaffleUsername       string `json:"WaffleUsername"`
	WaffleToken          string `json:"WaffleToken"`
	WaffleUserId         int64  `json:"WaffleUserId"`
}

func ApiHandleWaffleRegister(ctx *gin.Context) {
	formUsername := ctx.PostForm("username")
	formPassword := ctx.PostForm("password")

	ctx.Header("Access-Control-Allow-Origin", "127.0.0.1")
	ctx.Header("Access-Control-Allow-Origin", "http://localhost:3000")

	createSuccess := database.CreateNewUser(formUsername, formPassword)

	registerResponse := ApiRegisterResponse{}

	if createSuccess {
		queryResult, user := database.UserFromDatabaseByUsername(formUsername)

		if queryResult == -2 || queryResult == -1 {
			registerResponse.WaffleRegisterStatus = "User creation failed!"

			data, marshalErr := json.Marshal(registerResponse)

			if marshalErr != nil {
				ctx.String(http.StatusInternalServerError, "")
				return
			}

			ctx.Data(http.StatusOK, "waffle/blob", data)
			return
		}

		registerResponse.WaffleRegisterStatus = "User creation succeeded!"
		registerResponse.WaffleUsername = formUsername
		registerResponse.WaffleUserId = int64(user.UserID)
		registerResponse.WaffleToken = database.TokensCreateNewToken(user)

		data, marshalErr := json.Marshal(registerResponse)

		if marshalErr != nil {
			ctx.String(http.StatusInternalServerError, "")
			return
		}

		ctx.Data(http.StatusOK, "waffle/blob", data)
	} else {
		registerResponse.WaffleRegisterStatus = "User creation failed!"

		data, marshalErr := json.Marshal(registerResponse)

		if marshalErr != nil {
			ctx.String(http.StatusInternalServerError, "")
			return
		}

		ctx.Data(http.StatusOK, "waffle/blob", data)
		return
	}
}
