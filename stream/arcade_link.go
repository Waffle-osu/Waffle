package stream

import (
	"Waffle/config"
	"Waffle/database"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/blake2b"
)

func HandleArcadeLink(ctx *gin.Context) {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, ctx.Request.Body)

	if err != nil {
		ctx.Status(400)
		return
	}

	asString := buf.String()

	dataSplit := strings.Split(asString, "\n")

	if len(dataSplit) != 3 {
		ctx.Status(400)
		return
	}

	hashedCardId := blake2b.Sum512([]byte(dataSplit[0]))
	hashAsString := hex.EncodeToString(hashedCardId[:])

	linkCode := dataSplit[1]
	desiredPin := dataSplit[2]

	hashPinInput := fmt.Sprintf("%s%s", config.ArcadePinSalt, desiredPin)
	hashedPin := blake2b.Sum512([]byte(hashPinInput))
	hashedPinAsStr := hex.EncodeToString(hashedPin[:])

	linkCodeRow := database.Database.QueryRow(`
		SELECT 
			arcade_link_codes.user_id, users.username
		FROM 
			waffle.arcade_link_codes 
		LEFT JOIN 
			users ON arcade_link_codes.user_id = users.user_id
		WHERE 
			card_id = ? AND 
			link_code = ? AND 
			created_at >= CURRENT_TIMESTAMP - INTERVAL 5 MINUTE
		`, hashAsString, linkCode)

	var username string
	var userId uint64

	scanErr := linkCodeRow.Scan(&userId, &username)

	if scanErr != nil {
		ctx.Status(400)

		return
	}

	_, insertErr := database.Database.Exec("INSERT INTO waffle.arcade_cards (card_id, card_pin, user_id) VALUES (?, ?, ?)", hashAsString, hashedPinAsStr, userId)

	if insertErr != nil {
		ctx.Status(400)

		return
	}

	ctx.String(200, "success\n%d|%s|0.0", userId, username)
}
