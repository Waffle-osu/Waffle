package misc

import (
	"sync"
	"time"
)

//Fun Statistics
var StatsBytesRecieved uint64
var StatsBytesSent uint64
var StatsBanchoLaunch time.Time

var StatsSendLock sync.Mutex
var StatsRecvLock sync.Mutex

func InitializeStatistics() {
	StatsSendLock = sync.Mutex{}
	StatsRecvLock = sync.Mutex{}
	StatsBanchoLaunch = time.Now()
	StatsBytesRecieved = 0
	StatsBytesSent = 0
}

func ResetStatistics() {
	StatsSendLock.Lock()
	StatsRecvLock.Lock()

	StatsBytesRecieved = 0
	StatsBytesSent = 0
}