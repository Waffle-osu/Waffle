package spectator

import (
	"Waffle/utils"
)

var ClientManager utils.ClientManager[SpectatorClient] = utils.ClientManager[SpectatorClient]{}

func InitializeClientManager() {
	ClientManager.Initialize()
}
