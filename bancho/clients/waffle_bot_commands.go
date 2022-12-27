package clients

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/lobby"
	"Waffle/bancho/misc"
	"Waffle/database"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func WaffleBotCommandTemplate(sender client_manager.WaffleClient, args []string) []string {
	return []string{}
}

// WaffleBotCommandHelp !help
func WaffleBotCommandHelp(sender client_manager.WaffleClient, args []string) []string {
	returnStrings := helpStrings

	if (sender.GetUserData().Privileges & (chat.PrivilegesBAT | chat.PrivilegesAdmin)) > 0 {
		returnStrings = append(returnStrings, adminHelpStrings...)
	}

	return returnStrings
}

// WaffleBotCommandAnnounce !announce (both variants)
func WaffleBotCommandAnnounce(sender client_manager.WaffleClient, args []string) []string {
	defer func() {
		recover()
	}()

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

			targetClient.BanchoAnnounce(totalString)
			targetClient.BanchoGetAttention()
		}
	} else {
		client_manager.BroadcastPacketOsu(func(client client_manager.WaffleClient) {
			totalString := strings.Join(args[1:], " ")

			client.BanchoAnnounce(totalString)
		})
	}

	return []string{
		"Announcement has been issued.",
	}
}

// WaffleBotCommandRoll !roll <~max>
func WaffleBotCommandRoll(sender client_manager.WaffleClient, args []string) []string {
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

func WaffleBotCommandBanchoStatistics(sender client_manager.WaffleClient, args []string) []string {
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

//!rank <?mode> <username>
func WaffleBotCommandRank(sender client_manager.WaffleClient, args []string) []string {
	username := sender.GetUserData().Username
	mode := int8(0)
	writtenMode := "osu!"
	beginUsernameIndex := 1

	if len(args) != 0 {
		switch args[0] {
		case "osu!":
			mode = 0
			writtenMode = "osu!"
		case "osu!taiko":
			mode = 1
			writtenMode = "osu!taiko"
		case "osu!catch":
			mode = 2
			writtenMode = "osu!catch"
		default:
			beginUsernameIndex = 0
		}

		if len(args) > 1 {
			constructedUsername := ""

			if len(args) < 2 && beginUsernameIndex == 1 {
				return []string{
					"Invalid Command Format.",
				}
			}

			for i := beginUsernameIndex; i != len(args); i++ {
				constructedUsername += args[i] + " "
			}

			username = strings.TrimSpace(constructedUsername)
		}
	}

	userQueryResult, user := database.UserFromDatabaseByUsername(username)

	if userQueryResult == -2 {
		return []string{
			"Server Error occured. Could not retrieve user stats.",
		}
	} else if userQueryResult == -1 {
		return []string{
			"User not found: " + username,
		}
	}

	userStatsQueryResult, userStats := database.UserStatsFromDatabase(user.UserID, mode)

	if userStatsQueryResult == -2 {
		return []string{
			"Server Error occured. Could not retrieve user stats.",
		}
	}

	return []string{
		fmt.Sprintf("---------- User Statistics of %s for %s", username, writtenMode),
		fmt.Sprintf("Rank: %d", userStats.Rank),
		fmt.Sprintf("Ranked Score: %d", userStats.RankedScore),
		fmt.Sprintf("Total Score: %d", userStats.TotalScore),
		fmt.Sprintf("Level: %.2f", userStats.Level),
		fmt.Sprintf("Accuracy: %.2f%%", userStats.Accuracy*100.0),
		fmt.Sprintf("Playcount: %d", userStats.Playcount),
	}
}

func WaffleBotCommandLeaderboards(sender client_manager.WaffleClient, args []string) []string {
	offset := 0
	mode := int8(0)
	writtenMode := "osu!"

	if len(args) != 0 {
		if len(args) == 2 {
			parsedOffset, parseErr := strconv.ParseInt(args[0], 10, 64)

			if parseErr != nil {
				return []string{
					"Failed to load leaderboards, invalid offset.",
				}
			}

			offset = int(parsedOffset)

			switch args[1] {
			case "osu!":
				mode = 0
				writtenMode = "osu!"
			case "osu!taiko":
				mode = 1
				writtenMode = "osu!taiko"
			case "osu!catch":
				mode = 2
				writtenMode = "osu!catch"
			}
		} else {
			switch args[0] {
			case "osu!":
				mode = 0
				writtenMode = "osu!"
			case "osu!taiko":
				mode = 1
				writtenMode = "osu!taiko"
			case "osu!catch":
				mode = 2
				writtenMode = "osu!catch"
			default:
				parsedOffset, parseErr := strconv.ParseInt(args[0], 10, 64)

				if parseErr != nil {
					return []string{
						"Failed to load leaderboards, invalid offset.",
					}
				}

				offset = int(parsedOffset)
			}
		}
	}

	leaderboardQuery, leaderboardQueryErr := database.Database.Query("SELECT users.username, stats.* FROM (SELECT user_id, mode, ROW_NUMBER() OVER (ORDER BY ranked_score DESC) AS 'rank', ranked_score, total_score, user_level FROM waffle.stats WHERE mode = ? AND user_id != 1) stats LEFT JOIN waffle.users ON stats.user_id = users.user_id LIMIT 10 OFFSET ?", mode, offset)

	if leaderboardQueryErr != nil {
		if leaderboardQuery != nil {
			leaderboardQuery.Close()
		}

		return []string{
			"Failed to load leaderboards, query failed.",
		}
	}

	returnResults := []string{
		fmt.Sprintf("Showing leaderboards for %s #%d - #%d", writtenMode, offset, offset+10),
	}

	for leaderboardQuery.Next() {
		var username string
		partialStats := database.UserStats{}

		scanErr := leaderboardQuery.Scan(&username, &partialStats.UserID, &partialStats.Mode, &partialStats.Rank, &partialStats.RankedScore, &partialStats.TotalScore, &partialStats.Level)

		if scanErr != nil {
			leaderboardQuery.Close()

			return []string{
				"Failed to load leaderboards, query failed.",
			}
		}

		returnResults = append(returnResults, fmt.Sprintf("#%d %s - Score: %d (%d) Level %.0f", partialStats.Rank, username, partialStats.RankedScore, partialStats.TotalScore, partialStats.Level))
	}

	leaderboardQuery.Close()

	return returnResults
}

//!silence <duration in minutes> <username>
func WaffleBotCommandSilence(sender client_manager.WaffleClient, args []string) []string {
	//Check privileges
	if (chat.PrivilegesAdmin & sender.GetUserData().Privileges) <= 0 {
		if (chat.PrivilegesBAT & sender.GetUserData().Privileges) <= 0 {
			return []string{
				fmt.Sprintf("%s - you don't have the required privileges to execute !announce", sender.GetUserData().Username),
			}
		}
	}

	if len(args) < 2 {
		return []string{
			"Invalid Command Format.",
		}
	}

	duration, err := strconv.ParseInt(args[0], 10, 64)

	if err != nil {
		return []string{
			"Invalid Command Format.",
		}
	}

	constructedUsername := ""

	for i := 1; i != len(args); i++ {
		constructedUsername += args[i] + " "
	}

	constructedUsername = strings.TrimSpace(constructedUsername)

	client := client_manager.GetClientByName(constructedUsername)

	if client == nil {
		return []string{
			"Client not found!.",
		}
	}

	client.SetSilencedUntilUnix(time.Now().Unix() + (duration * 60))

	return []string{
		fmt.Sprintf("%s is now silenced for %d minutes", constructedUsername, duration),
	}
}
