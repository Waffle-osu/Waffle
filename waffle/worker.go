package waffle

import (
	"math"
	"time"
)

type BanchoWorker struct {
	Bancho                  *Bancho
	Id                      int
	LastProcessdIndex       int32
	LastClientHandleRequest time.Time
}

func CreateNewWorker(id int, bancho *Bancho, decommision chan struct{}) {
	continueWork := true

	go WorkerWorkFunction(&continueWork, BanchoWorker{bancho, id, 0, time.Now()})

	<-decommision //If something's received on the decommission channel, we decommission

	continueWork = false
}

func WorkerWorkFunction(continueWork *bool, workerInformation BanchoWorker) {
	for *continueWork == true {
		clientCount := len(workerInformation.Bancho.Clients)

		if clientCount == 0 {
			continue
		}

		//The amount of casts in these disgusts me
		indexRange := math.Ceil(float64(clientCount) / float64(len(workerInformation.Bancho.WorkerChannels)))
		indexStart := indexRange * float64(workerInformation.Id)

		index := math.Min(float64(clientCount-1), indexStart+float64(workerInformation.LastProcessdIndex))

		if index < indexStart {
			continue
		}

		workerInformation.Bancho.ClientMutex.Lock()
		client := &workerInformation.Bancho.Clients[int32(index)]
		workerInformation.Bancho.ClientMutex.Unlock()

		workerInformation.LastProcessdIndex = (workerInformation.LastProcessdIndex + 1) % int32(indexRange)

		(*client).HandleIncoming()
		(*client).SendOutgoing()
	}
}
