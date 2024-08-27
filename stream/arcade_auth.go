package stream

import (
	"Waffle/database"
	"bytes"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/blake2b"
)

func HandleArcadeAuth(ctx *gin.Context) {
	buf := new(bytes.Buffer)

	_, err := buf.ReadFrom(ctx.Request.Body)

	if err != nil {
		ctx.String(400, "")
		return
	}

	hashedCardId := blake2b.Sum512(buf.Bytes())
	hashAsString := hex.EncodeToString(hashedCardId[:])

	row := database.Database.QueryRow("SELECT COUNT(*) FROM arcade_cards WHERE card_id = ?", hashAsString)

	var count int64

	scanErr := row.Scan(&count)

	if scanErr != nil {
		ctx.String(502, "")
		return
	}

	if count == 1 {
		userRow := database.Database.QueryRow("SELECT users.username, arcade_cards.user_id FROM arcade_cards LEFT JOIN users ON arcade_cards.user_id = users.user_id WHERE card_id = ?", hashAsString)

		var username string
		var userId uint64

		userScanErr := userRow.Scan(&username, &userId)

		if userScanErr != nil {
			ctx.String(502, "")
		}

		ctx.String(200, "exists\n%d|%s|69.69", userId, username)
	} else {
		ctx.String(200, "new")
	}
}
