package clients

import (
	"time"
)

func (client *WaffleBot) WaffleBotMaintainClient() {
	for client.continueRunning {
		//Maybe i'll add some fancy stuff here like funny statuses but as it stands this will be empty
		time.Sleep(time.Second)
	}
}
