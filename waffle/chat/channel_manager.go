package chat

import (
	"sync"
)

var channels map[string]*Channel
var channelList []*Channel

func InitializeChannels() {
	channels = map[string]*Channel{
		"#osu":      {"#osu", "The main channel of osu!", false, false, true, []ChatClient{}, sync.Mutex{}},
		"#announce": {"#announce", "The main channel of osu!", false, true, true, []ChatClient{}, sync.Mutex{}},
		"#lobby":    {"#lobby", "Find people to multi with here!", false, false, false, []ChatClient{}, sync.Mutex{}},
	}
}

func GetChannelByName(name string) (channel *Channel, exists bool) {
	foundChannel, found := channels[name]

	return foundChannel, found
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
