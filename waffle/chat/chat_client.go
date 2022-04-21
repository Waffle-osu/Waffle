package chat

import (
	"container/list"
	"sync"
)

type ChatClient interface {
	IsOfAdminPrivileges() bool
	GetPacketQueue() *list.List
	GetPacketQueueMutex() *sync.Mutex
}
