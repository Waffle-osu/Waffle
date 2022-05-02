package clients

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/packets"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

var helpStrings = []string{
	"!help :: You're reading this right now",
	"!roll :: Rolls a random number between 0 and 100",
}

var adminHelpStrings = []string{
	"---------------------------------",
	"!announce target <client username> : <message> :: Sends a Notification to a client",
	"^^^ That : seperator is important there!!",
	"!announce all <message> :: Sends a Notification to everyone on the server",
}

func WaffleBotCommandTemplate(sender client_manager.OsuClient, args []string) []string {
	return []string{}
}

// WaffleBotCommandHelp !help
func WaffleBotCommandHelp(sender client_manager.OsuClient, args []string) []string {
	returnStrings := helpStrings

	if (sender.GetUserData().Privileges & (chat.PrivilegesBAT | chat.PrivilegesAdmin)) > 0 {
		returnStrings = append(returnStrings, adminHelpStrings...)
	}

	return returnStrings
}

// WaffleBotCommandAnnounce !announce (both variants)
func WaffleBotCommandAnnounce(sender client_manager.OsuClient, args []string) []string {
	//Check privileges
	if (chat.PrivilegesAdmin & sender.GetUserData().Privileges) <= 0 {
		return []string{
			fmt.Sprintf("%s - you don't have the required privileges to execute !announce", sender.GetUserData().Username),
		}
	}

	toAll := false

	if args[0] == "all" {
		toAll = true
	} else if args[0] == "target" {
		toAll = false
	}

	if toAll == false {
		target := ""
		index := 0

		//This is handling for Usernames with Spaces
		for {
			index++

			currentElement := args[index]

			//this is the delimiter, if we hit this, it's the end of the username part
			if currentElement != ":" {
				target += currentElement + " "
			} else {
				index++
				target = strings.TrimSpace(target)
				break
			}
		}

		targetClient := client_manager.GetClientByName(target)

		if targetClient != nil {
			totalString := strings.Join(args[index:], " ")

			packets.BanchoSendAnnounce(targetClient.GetPacketQueue(), totalString)
			packets.BanchoSendGetAttention(targetClient.GetPacketQueue())
		}
	} else {
		client_manager.BroadcastPacket(func(packetQueue chan packets.BanchoPacket) {
			totalString := strings.Join(args[1:], " ")

			packets.BanchoSendAnnounce(packetQueue, totalString)
		})
	}

	return []string{
		"Announcement has been issued.",
	}
}

func WaffleBotCommandRoll(sender client_manager.OsuClient, args []string) []string {
	max := 100.0

	if len(args) != 0 {
		float, err := strconv.ParseFloat(args[0], 64)

		if err == nil {
			max = float
		}
	}

	rolled := math.Round(rand.Float64() * max)

	return []string{
		fmt.Sprintf("%s rolled %d!", sender.GetUserData().Username, int(rolled)),
	}
}
