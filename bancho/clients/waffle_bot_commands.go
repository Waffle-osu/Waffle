package clients

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/lobby"
	"Waffle/bancho/misc"
	"Waffle/bancho/packets"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var helpStrings = []string{
	"!help  :: You're reading this right now",
	"!roll  :: Rolls a random number between 0 and 100",
	"!stats :: Shows Waffle Statistics",
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

	if !toAll {
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

// WaffleBotCommandRoll !roll <~max>
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

func WaffleBotCommandBanchoStatistics(sender client_manager.OsuClient, args []string) []string {
	//Calculating Uptime
	var uptimeString string

	duration := time.Since(misc.StatsBanchoLaunch)

	hours := duration.Hours()

	//Crazy math go get normal 60 seconds, 60 minutes and stuff
	hours, minuteFraction := math.Modf(hours)
	minutes := minuteFraction * 60

	minutes, secondFraction := math.Modf(minutes)
	seconds := secondFraction * 60

	if int(duration.Hours()) == 0 {
		//Less than an hour has passed
		if int(duration.Minutes()) == 0 {
			//Less than a minute has passed
			pluralSeconds := "seconds"

			if int(seconds) == 1 {
				pluralSeconds = "second"
			}

			uptimeString = fmt.Sprintf("%.0f %s", seconds, pluralSeconds)
		} else {
			pluralSeconds := "seconds"

			if int(seconds) == 1 {
				pluralSeconds = "second"
			}

			pluralMinutes := "minutes"

			if int(minutes) == 1 {
				pluralMinutes = "minute"
			}

			//More than a minute but less than an hour has passed
			uptimeString = fmt.Sprintf("%.0f %s and %.0f %s", minutes, pluralMinutes, seconds, pluralSeconds)
		}
	} else {
		pluralSeconds := "seconds"

		if int(seconds) == 1 {
			pluralSeconds = "second"
		}

		pluralMinutes := "minutes"

		if int(minutes) == 1 {
			pluralMinutes = "minute"
		}

		pluralHours := "hours"

		if int(hours) == 1 {
			pluralHours = "hour"
		}
		//At least an hour has passed
		uptimeString = fmt.Sprintf("%.0f %s %.0f %s and %.0f %s", hours, pluralHours, minutes, pluralMinutes, seconds, pluralSeconds)
	}

	//Calculating data sent/recieved
	var dataSentString string
	var dataRecvString string

	//Recieved
	if misc.StatsBytesRecieved < 1024*1024*1024 {
		//Less than a Gigabyte has been recieved
		if misc.StatsBytesRecieved < 1024*1024 {
			//Less than a Megabyte has been recieved
			if misc.StatsBytesRecieved < 1024 {
				//Less than a Kilobyte has been recieved
				dataRecvString = fmt.Sprintf("%d bytes", misc.StatsBytesRecieved)
			} else {
				//More than a Kilobyte has been recieved
				dataRecvString = fmt.Sprintf("%.2fkb", float64(misc.StatsBytesRecieved)/1024.0)
			}
		} else {
			//More than a Megabyte has been recieved
			dataRecvString = fmt.Sprintf("%.2fmb", (float64(misc.StatsBytesRecieved)/1024.0)/1024.0)
		}
	} else {
		//More than a Gigabyte has been recieved
		dataRecvString = fmt.Sprintf("%.2fgb", ((float64(misc.StatsBytesRecieved)/1024.0)/1024.0)/1024.0)
	}

	//Sent
	if misc.StatsBytesSent < 1024*1024*1024 {
		//Less than a Gigabyte has been sent
		if misc.StatsBytesSent < 1024*1024 {
			//Less than a Megabyte has been sent
			if misc.StatsBytesSent < 1024 {
				//Less than a Kilobyte has been sent
				dataSentString = fmt.Sprintf("%d bytes", misc.StatsBytesSent)
			} else {
				//More than a Kilobyte has been sent
				dataSentString = fmt.Sprintf("%.2fkb", float64(misc.StatsBytesSent)/1024.0)
			}
		} else {
			//More than a Megabyte has been sent
			dataSentString = fmt.Sprintf("%.2fmb", (float64(misc.StatsBytesSent)/1024.0)/1024.0)
		}
	} else {
		//More than a Gigabyte has been sent
		dataSentString = fmt.Sprintf("%.2fgb", ((float64(misc.StatsBytesSent)/1024.0)/1024.0)/1024.0)
	}

	pluralUsers := "users"
	pluralMatches := "matches"

	if client_manager.GetClientCount() == 1 {
		pluralUsers = "user"
	}

	if lobby.GetMatchCount() == 1 {
		pluralMatches = "match"
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	mbAlloc := (float64(m.Alloc) / 1024.0) / 1024.0

	return []string{
		fmt.Sprintf("[WAFFLE-STATS] Waffle has been up for %s", uptimeString),
		fmt.Sprintf("[WAFFLE-STATS] Serving %d %s, playing %d %s", client_manager.GetClientCount(), pluralUsers, lobby.GetMatchCount(), pluralMatches),
		fmt.Sprintf("[WAFFLE-STATS] %s have been sent", dataSentString),
		fmt.Sprintf("[WAFFLE-STATS] %s have been recieved", dataRecvString),
		fmt.Sprintf("[WAFFLE-STATS] %d Goroutines are currently running", runtime.NumGoroutine()),
		fmt.Sprintf("[WAFFLE-STATS] Currently using approximately %.2fmb RAM", mbAlloc),
	}
}
