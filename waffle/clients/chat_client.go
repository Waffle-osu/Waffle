package clients

import (
	"container/list"
	"sync"
)

func (client *Client) IsOfAdminPrivileges() bool {
	return client.UserData.Privileges&16 > 0
}

func (client *Client) GetPacketQueue() *list.List {
	return client.PacketQueue
}

func (client *Client) GetPacketQueueMutex() *sync.Mutex {
	return &client.PacketQueueMutex
}
