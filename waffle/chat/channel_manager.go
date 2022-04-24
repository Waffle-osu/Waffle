package chat

import (
	"sync"
)

var channels map[string]*Channel
var channelList []*Channel

const (
	PrivilegesNormal    int32 = 1
	PrivilegesBAT       int32 = 2
	PrivilegesSupporter int32 = 4
	PrivilegesFriend    int32 = 8
	PrivilegesAdmin     int32 = 16
)

func InitializeChannels() {
	channels = map[string]*Channel{
		"#osu":      {"#osu", "The main channel of osu!", PrivilegesNormal, PrivilegesNormal, true, []ChatClient{}, sync.Mutex{}},
		"#announce": {"#announce", "The main channel of osu!", PrivilegesNormal, PrivilegesBAT | PrivilegesAdmin, true, []ChatClient{}, sync.Mutex{}},
		"#lobby":    {"#lobby", "Find people to multi with here!", PrivilegesNormal, PrivilegesNormal, false, []ChatClient{}, sync.Mutex{}},
		"#bat":      {"#bat", "Staff channel for BAT's", PrivilegesBAT | PrivilegesAdmin, PrivilegesBAT | PrivilegesAdmin, false, []ChatClient{}, sync.Mutex{}},
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
