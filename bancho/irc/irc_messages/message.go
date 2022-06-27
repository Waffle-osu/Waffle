package irc_messages

import (
	"fmt"
	"strings"
)

type Message struct {
	Source   string
	Command  string
	Params   []string
	Trailing string
}

func ParseMessage(line string) Message {
	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSuffix(line, "\r")

	returnMessage := Message{}

	getToken := func() string {
		split := strings.SplitN(line, " ", 2)

		if len(split) > 1 {
			line = split[1]
		} else {
			line = ""
		}

		return split[0]
	}

	if line[0] == ':' {
		sourceSplit := strings.SplitN(line[1:], " ", 2)

		returnMessage.Source = sourceSplit[0]

		if len(sourceSplit) > 1 {
			line = sourceSplit[1]
		}
	}

	returnMessage.Command = getToken()

	for len(line) != 0 {
		if line[0] == ':' {
			returnMessage.Trailing = strings.TrimSpace(line[1:])
			break
		} else {
			param := getToken()

			returnMessage.Params = append(returnMessage.Params, param)
		}
	}

	return returnMessage
}

func (message *Message) FormatMessage() (formatted string, formatErr string) {
	if len(message.Params) == 0 || message.Trailing == "" {
		return "", "Either Parameters or Trailing has to be set!"
	}

	returnString := fmt.Sprintf("%s ", message.Command)

	if len(message.Params) != 0 {
		returnString = fmt.Sprintf("%s %s", returnString, strings.Join(message.Params, " "))
	}

	if len(message.Trailing) != 0 {
		returnString = fmt.Sprintf("%s %s", returnString, message.Trailing)
	}

	return returnString, ""
}
