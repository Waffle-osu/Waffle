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
[Initialization] //                                                    //
[Initialization] // ssl_silence_warning: silences the warning about    //
[Initialization] //                      IRC SSL Certificates missing  //
[Initialization] //                                                    //
[Initialization] // ssl_key:   Location of the Private Key File.       //
[Initialization] // ssl_cert:  Location of the SSL Certificate.        //
[Initialization] //                                                    //
[Initialization] // bancho_ip:  Where the Bancho TCP Listener is hosted//
[Initialization] // irc_ip:     Where the IRC TCP Listener is hosted   //
[Initialization] // irc_ssl_ip: Where IRC/SSL is hosted                //
[Initialization] //                                                    //
[Initialization] // host_irc:      Enables/Disables IRC Server         //
[Initialization] // host_irc_ssl:  Enables/Disables IRC SSL Server     //
[Initialization] //                                                    //
[Initialization] // using_waffle_web: Whether the original waffle-web  //
[Initialization] //                   is used, as opposed to a custom  //
[Initialization] //                   made frontend to go around it.   //
[Initialization] //                   it used around beatmap forum     //
[Initialization] //                   posts and ranking queue.         //
[Initialization] //                                                    //
[Initialization] ////////////////////////////////////////////////////////
`

type ExpectedKey struct {
	Key      string
	Critical bool
}

type KeyValuePair struct {
	Key   string
	Value string
}

var ExpectedKeys map[ExpectedKey]func() = map[ExpectedKey]func(){
	{"mysql_location", true}:       MySqlSettingsIncompleteError,
	{"mysql_database", true}:       MySqlSettingsIncompleteError,
	{"mysql_username", true}:       MySqlSettingsIncompleteError,
	{"mysql_password", true}:       MySqlSettingsIncompleteError,
	{"token_format", false}:        TokenFormatWarning,
	{"ssl_silence_warning", false}: nil,
	{"ssl_key", false}:             IrcSslCertsWarning,
	{"ssl_cert", false}:            IrcSslCertsWarning,
	{"bancho_ip", true}:            BanchoIpWarning,
	{"host_irc", false}:            nil,
	{"host_irc_ssl", false}:        nil,
	{"irc_ip", false}:              IrcIpMissing,
	{"irc_ssl_ip", false}:          IrcSSLIpMissing,
	{"ffmpeg_executable", false}:   FfmpegMissing,
}

var DefaultSettings map[string]string = map[string]string{
	"token_format": "wa%sff%sle%dto%dke%sn",
	"bancho_ip":    "127.0.0.1:13381",
	"host_irc":     "true",
	"host_irc_ssl": "true",
	"irc_ip":       "127.0.0.1:6667",
	"irc_ssl_ip":   "127.0.0.1:6697",
}

func ReadConfiguration() {
	existingKeys := map[string]bool{}

	// Check if .env file exists, if not print first run warning
	if _, err := os.Stat(".env"); errors.Is(err, os.ErrNotExist) {
		succeeded := CreateDefaultConfiguration()

		if !succeeded {
			return
		}

		for _, line := range strings.Split(EnvFileWarningString, "\n") {
			helpers.Logger.Printf(line)
		}

		return
	}

	data, err := os.ReadFile(".env")

	if err != nil {
		helpers.Logger.Fatalf("[Initialization] Failed to read configuration file, cannot start server!\n")
	}

	asString := string(data)

	// the .env is essentially a ini without groups
	splitLines := strings.Split(asString, "\n")

	lineCounter := 1

	for _, iniEntry := range splitLines {
		if strings.HasPrefix(iniEntry, "#") {
			continue
		}

		if len(iniEntry) == 0 {
			continue
		}

		splitEntry := strings.Split(iniEntry, "=")

		if len(splitEntry) != 2 {
			helpers.Logger.Printf("[Initialization] .env file has an error on line %d, ignoring.\n", lineCounter)

			continue
		}

		key := splitEntry[0]
		value := splitEntry[1]

		if value == "" {
			continue
		}

		existingKeys[key] = true

		UnsafeSetKey(key, value)
	}

	warningsToRun := []func(){}
	defaultSetWarningsToRun := []KeyValuePair{}

	displayAllWarnings := func() {
		for _, warningFunction := range warningsToRun {
			warningFunction()
		}
	}

	displayAllSetDefaultWarnings := func() {
		for _, kv := range defaultSetWarningsToRun {
			helpers.Logger.Printf("[Initialization] Key %s has been set to the default value %s\n", kv.Key, kv.Value)
		}
	}

	// This gives out warnings for all the unset parameters
	// Also fully exits if a critial key is missing.
	for key, value := range ExpectedKeys {
		exists := existingKeys[key.Key]

		if !exists {
			defaultValue, defaultExists := DefaultSettings[key.Key]

			if value != nil && !defaultExists {
				warningsToRun = append(warningsToRun, value)
			}

			if key.Critical && !defaultExists {
				displayAllWarnings()

				helpers.Logger.Fatalf("[Initialization] Critical Configuration key is missing. Exiting.\n")
			}

			if defaultExists {
				defaultSetWarningsToRun = append(defaultSetWarningsToRun, KeyValuePair{
					Key:   key.Key,
					Value: defaultValue,
				})

				UnsafeSetKey(key.Key, defaultValue)
			}
		}
	}

	go func() {
		time.Sleep(2 * time.Second)

		displayAllWarnings()

		helpers.Logger.Printf("[Initialization] ----- Keep in mind set to default does not mean it will save in .env -----\n")

		displayAllSetDefaultWarnings()
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

token_format=wa%sff%sle%dto%dke%sn

# Arcade PIN Salt
# Used to add an unknown to the Arcade PIN hashes
# So in the event of a database leak the PIN can't simply be brute forced (which would be quite easy)
arcade_pin_salt=(678) 999-8212@osu!arcade

# Uncomment the Following line if you wish to silence the SSL Certificate missing warning.
# ssl_silence_warning=true

bancho_ip=127.0.0.1:13381

# Uncomment this line if you wish to not create a IRC Server for Waffle.
#
#host_irc=false
irc_ip=127.0.0.1:6667

# Uncomment this line if you wish to not create a SSL IRC Server for Waffle.
#
#host_irc_ssl=false
irc_ssl_ip=127.0.0.1:6697

`

	writeErr := os.WriteFile(".env", []byte(defaultConfiguration), 0644)

	if writeErr != nil {
		helpers.Logger.Printf("[Initialization] Failed to create default configuration!\n\n")
		return false
	}

	return true
}

func UnsafeSetKey(key string, value string) {
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
	case "bancho_ip":
		BanchoIp = value
	case "host_irc":
		HostIrc = value
	case "host_irc_ssl":
		HostIrcSsl = value
	case "irc_ip":
		IrcIp = value
	case "irc_ssl_ip":
		IrcSslIp = value
	case "using_waffle_web":
		UsingWaffleWeb = value
	case "ffmpeg_path":
		FFMPEGPath = value
	case "arcade_pin_salt":
		ArcadePinSalt = value
	}
}
