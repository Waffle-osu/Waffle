package chat

import (
	"sync"
)

var channels map[string]*Channel
var channelList []*Channel

func InitializeChannels() {
	channels = map[string]*Channel{
		"#osu":      {"#osu", "The main channel of osu!", false, false, []ChatClient{}, sync.Mutex{}},
		"#announce": {"#announce", "The main channel of osu!", false, true, []ChatClient{}, sync.Mutex{}},
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

func TryJoinChannel(client ChatClient, channelName string) (joinSuccess bool, joinedChannel *Channel) {
	channel, exists := channels[channelName]

	if exists == false {
		return false, nil
	}

	return channel.Join(client), channel
}
