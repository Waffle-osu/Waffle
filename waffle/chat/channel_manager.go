package chat

import (
	"sync"
)

var channels map[string]*Channel
var channelList []*Channel

func InitializeChannels() {
	channels = map[string]*Channel{
		"#osu":      {"#osu", "The main channel of osu!", false, []AdminPrivilegable{}, sync.Mutex{}},
		"#announce": {"#announce", "The main channel of osu!", false, []AdminPrivilegable{}, sync.Mutex{}},
	}
}

func GetAvailableChannels() []*Channel {
	if channelList == nil {
		for _, value := range channels {
			channelList = append(channelList, value)
		}
	}

	return channelList
}

func TryJoinChannel(client AdminPrivilegable, channelName string) bool {
	channel, exists := channels[channelName]

	if exists == false {
		return false
	}

	return channel.Join(client)
}

func LeaveChannel(client AdminPrivilegable, channelName string) {
	channel, exists := channels[channelName]

	if exists == false {
		return
	}

	channel.Leave(client)
}
