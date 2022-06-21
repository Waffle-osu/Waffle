package config

import (
	"Waffle/helpers"
	"errors"
	"os"
	"strings"
	"time"
)

var EnvFileWarningString string = `
[Initialization] ////////////////////////////////////////////////////////
[Initialization] //////////////////  First Run Advice  //////////////////
[Initialization] ////////////////////////////////////////////////////////
[Initialization] // No .env file detected in the Waffle.exe directory! //
[Initialization] //    This file stores important configuration for    //
[Initialization] //      The server, such as Database Credentials,     //
[Initialization] //                                                    //
[Initialization] // A .env file with default settings has been created //
[Initialization] //      Please change the settings as necessary       //
[Initialization] //                                                    //
[Initialization] //    Explanation to all the keys in the .env file:   //
[Initialization] //                                                    //
[Initialization] // mysql_location: Location of the MySQL Server       //
[Initialization] // mysql_database: Name of the Database to use        //
[Initialization] // mysql_username: Under which user to log in         //
[Initialization] // mysql_password: Password for said user             //
[Initialization] // token_format:   Format string specifier for tokens //
[Initialization] //                 make sure to change it, otherwise  //
[Initialization] //                 you run the risk (even if small)   //
[Initialization] //                 that your tokens will be forged.   //
[Initialization] //                 Your format string needs to        //
[Initialization] //                 contain 2 %%s's 2 %%d's and a %%s     //
[Initialization] //                 in this exact order.               //
[Initialization] //                                                    //
[Initialization] // example:        wa%%sff%%sle%%dto%%dke%%sn              //
[Initialization] ////////////////////////////////////////////////////////
`

type ExpectedKey struct {
	Key      string
	Critical bool
}

var ExpectedKeys map[ExpectedKey]func()string = map[ExpectedKey]func()string{
	{"mysql_location", true}:       MySqlSettingsIncompleteError,
	{"mysql_database", true}:       MySqlSettingsIncompleteError,
	{"mysql_username", true}:       MySqlSettingsIncompleteError,
	{"mysql_password", true}:       MySqlSettingsIncompleteError,
	{"token_format", false}:        TokenFormatWarning,
	{"ssl_silence_warning", false}: nil,
	{"ssl_key", false}:             IrcSslCertsWarning,
	{"ssl_cert", false}:            IrcSslCertsWarning,
}

var DefaultSettings map[string]string = map[string]string{
	"token_format": "wa%sff%sle%dto%dke%sn",
}

func ReadConfiguration() {
	existingKeys := map[string]bool{}

	if _, err := os.Stat(".env"); errors.Is(err, os.ErrNotExist) {
		go func() {
			time.Sleep(2 * time.Second)

			succeeded := CreateDefaultConfiguration()

			if !succeeded {
				return
			}

			for _, line := range strings.Split(EnvFileWarningString, "\n") {
				helpers.Logger.Printf(line)
			}
		}()

		return
	}

	data, err := os.ReadFile(".env")

	if err != nil {
		helpers.Logger.Fatalf("[Initialization] Failed to read configuration file, cannot start server!")
	}

	asString := string(data)

	splitLines := strings.Split(asString, "\n")

	lineCounter := 1

	for _, iniEntry := range splitLines {
		if strings.HasPrefix(iniEntry, "#") {
			continue
		}

		splitEntry := strings.Split(iniEntry, "=")

		if len(splitEntry) != 2 {
			helpers.Logger.Printf("[Initialization] .env file has an error on line %d, ignoring.", lineCounter)

			continue
		}

		key := splitEntry[0]
		value := splitEntry[1]

		existingKeys[key] = true

		switch key {
		case "mysql_username":
			MySqlUsername = value
		case "mysql_password":
			MySqlPassword = value
		case "mysql_database":
			MySqlDatabase = value
		case "mysql_location":
			MySqlLocation = value
		case "token_format":
			TokenFormatString = value
		case "ssl_silence_warning":
			SSLSilenceWarning = value
		case "ssl_key":
			SSLKeyLocation = value
		case "ssl_cert":
			SSLCertLocation = value
		}
	}

	warningsToRun := []func(){}
	warningsAlreadyRun := map[string]bool{}

	displayAllWarnings := func() {
		for _, warningFunction := range warningsToRun {
			warningFunction()
		}
	}

	for key, value := range ExpectedKeys {
		exists := existingKeys[key.Key]

		if !exists {
			if value != nil {
				warningsToRun = append(warningsToRun, value)
			}

			if key.Critical {
				displayAllWarnings()

				helpers.Logger.Fatalf("[Initialization] Critical Configuration key is missing. Exiting.")
			}
		}
	}

	go func() {
		time.Sleep(2 * time.Second)

		displayAllWarnings()
	}()
}

func CreateDefaultConfiguration() bool {
	defaultConfiguration :=
		`
mysql_username=root
mysql_password=root
mysql_location=127.0.0.1:3306
mysql_database=waffle

# Format string specifier for tokens.
# make sure to change it, otherwise you run the risk (even if small) that your tokens will be forged.
# Your format string needs to contain 2 %s's 2 %d's and a %s in this exact order.
# 
# example: wa%sff%sle%dto%dke%sn

tokenformat=wa%sff%sle%dto%dke%sn

# Uncomment the Following line if you wish to silence the SSL Certificate missing warning.
# ssl_silence_warning=true

`

	writeErr := os.WriteFile(".env", []byte(defaultConfiguration), 0644)

	if writeErr != nil {
		helpers.Logger.Printf("[Initialization] Failed to create default configuration!\n\n")
		return false
	}

	return true
}
