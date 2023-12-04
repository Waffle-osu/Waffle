package helpers

import (
	"fmt"
	"math"
	"time"
)

func FormatTime(milliseconds uint64) string {
	var formatted string

	duration := time.Duration(milliseconds) * time.Millisecond

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

			formatted = fmt.Sprintf("%.0f %s", seconds, pluralSeconds)
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
			formatted = fmt.Sprintf("%.0f %s and %.0f %s", minutes, pluralMinutes, seconds, pluralSeconds)
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
		formatted = fmt.Sprintf("%.0f %s %.0f %s and %.0f %s", hours, pluralHours, minutes, pluralMinutes, seconds, pluralSeconds)
	}

	return formatted
}
