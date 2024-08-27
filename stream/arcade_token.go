package stream

import (
	"Waffle/config"
	"Waffle/database"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/blake2b"
)

func HandleArcadeToken(ctx *gin.Context) {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, ctx.Request.Body)

	if err != nil {
		ctx.Status(400)
		return
	}

	asString := buf.String()

	dataSplit := strings.Split(asString, "\n")

	if len(dataSplit) != 2 {
		ctx.Status(400)
		return
	}

	hashedCardId := blake2b.Sum512([]byte(dataSplit[0]))
	hashAsString := hex.EncodeToString(hashedCardId[:])

	enteredPin := dataSplit[1]

	hashPinInput := fmt.Sprintf("%s%s", config.ArcadePinSalt, enteredPin)
	hashedPin := blake2b.Sum512([]byte(hashPinInput))
	hashedPinAsStr := hex.EncodeToString(hashedPin[:])

	cardRow := database.Database.QueryRow("SELECT COUNT(*) FROM arcade_cards WHERE card_id = ? AND card_pin = ?", hashAsString, hashedPinAsStr)

	var count int64

	scanErr := cardRow.Scan(&count)

	if scanErr != nil {
		ctx.Status(502)
		return
	}

	if count == 0 {
		ctx.String(200, "pin")
		return
	}

	hashInput := fmt.Sprintf(config.TokenFormatString, hashedCardId, hashedPin, time.Now().UnixMilli(), 0, time.Now().UnixMilli())
	hashedToken := blake2b.Sum512([]byte(hashInput))
	hashedTokenAsStr := hex.EncodeToString(hashedToken[:])

	ctx.String(200, hashedTokenAsStr)
}
