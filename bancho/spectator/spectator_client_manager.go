package spectator

import (
	"Waffle/common"
)

var ClientManager common.ClientManager[SpectatorClient] = common.ClientManager[SpectatorClient]{}

func InitializeClientManager() {
	ClientManager.Initialize()
}
