package config

import "Waffle/helpers"

// All the warnings that've been displayed already
// So we don't display them twice.
var displayedWarnings map[string]bool = map[string]bool{}

/*
	Contains all the Warnings for not having certain keys set
	in the .env file, which may or may not be important.
*/

func MySqlSettingsIncompleteError() {
	_, runAlready := displayedWarnings["mysql"]

	if runAlready {
		return
	}

	helpers.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")
	helpers.Logger.Printf("[Initialization] //////////////////  Attention!!!!!!!  //////////////////\n")
	helpers.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")
	helpers.Logger.Printf("[Initialization] // Server cannot run unless all of those keys are set //\n")
	helpers.Logger.Printf("[Initialization] // All of those relate to the MySQL Database          //\n")
	helpers.Logger.Printf("[Initialization] // connection, which is required to run Waffle.       //\n")
	helpers.Logger.Printf("[Initialization] // It seems that one of the following keys is missing //\n")
	helpers.Logger.Printf("[Initialization] // Please make sure to fill all of those out and      //\n")
	helpers.Logger.Printf("[Initialization] // restart Waffle.                                    //\n")
	helpers.Logger.Printf("[Initialization] //                                                    //\n")
	helpers.Logger.Printf("[Initialization] // mysql_location: Location of the MySQL Server       //\n")
	helpers.Logger.Printf("[Initialization] // mysql_database: Name of the Database to use        //\n")
	helpers.Logger.Printf("[Initialization] // mysql_username: Under which user to log in         //\n")
	helpers.Logger.Printf("[Initialization] // mysql_password: Password for said user             //\n")
	helpers.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")

	displayedWarnings["mysql"] = true
}

func TokenFormatWarning() {
	_, runAlready := displayedWarnings["token_format"]

	if runAlready {
		return
	}

	helpers.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")
	helpers.Logger.Printf("[Initialization] //////////////////  Attention!!!!!!!  //////////////////\n")
	helpers.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")
	helpers.Logger.Printf("[Initialization] // There is a missing key in your configuration file. //\n")
	helpers.Logger.Printf("[Initialization] // The missing key is the token_format.               //\n")
	helpers.Logger.Printf("[Initialization] // The example Token Format has been set as it is     //\n")
	helpers.Logger.Printf("[Initialization] // the default option.                                //\n")
	helpers.Logger.Printf("[Initialization] //                                                    //\n")
	helpers.Logger.Printf("[Initialization] // token_format:   Format string specifier for tokens //\n")
	helpers.Logger.Printf("[Initialization] //                 make sure to change it, otherwise  //\n")
	helpers.Logger.Printf("[Initialization] //                 you run the risk (even if small)   //\n")
	helpers.Logger.Printf("[Initialization] //                 that your tokens will be forged.   //\n")
	helpers.Logger.Printf("[Initialization] //                 Your format string needs to        //\n")
	helpers.Logger.Printf("[Initialization] //                 contain 2 %%s's 2 %%d's and a %%s     //\n")
	helpers.Logger.Printf("[Initialization] //                 in this exact order.               //\n")
	helpers.Logger.Printf("[Initialization] //                                                    //\n")
	helpers.Logger.Printf("[Initialization] // example:        wa%%sff%%sle%%dto%%dke%%sn              //\n")
	helpers.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")

	displayedWarnings["token_format"] = true
}

func BanchoIpWarning() {
	_, runAlready := displayedWarnings["bancho_ip"]

	if runAlready {
		return
	}

	helpers.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")
	helpers.Logger.Printf("[Initialization] //////////////////  Attention!!!!!!!  //////////////////\n")
	helpers.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")
	helpers.Logger.Printf("[Initialization] // The Bancho IP in the .env file has been left empty //\n")
	helpers.Logger.Printf("[Initialization] // The Bancho IP is essential in running Waffle.      //\n")
	helpers.Logger.Printf("[Initialization] // You should definetly fill it out if you wish to    //\n")
	helpers.Logger.Printf("[Initialization] // run Waffle.                                        //\n")
	helpers.Logger.Printf("[Initialization] //                                                    //\n")
	helpers.Logger.Printf("[Initialization] // bancho_ip: Where the Bancho TCP Listener is hosted //\n")
	helpers.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")

	displayedWarnings["bancho_ip"] = true
}

func IrcSslCertsWarning() {
	if SSLSilenceWarning == "true" || HostIrcSsl == "false" {
		return
	}

	_, runAlready := displayedWarnings["ssl"]

	if runAlready {
		return
	}

	helpers.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")
	helpers.Logger.Printf("[Initialization] //////////////////  Attention!!!!!!!  //////////////////\n")
	helpers.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")
	helpers.Logger.Printf("[Initialization] // One of the two keys listed below is missing.       //\n")
	helpers.Logger.Printf("[Initialization] // These configuration keys are critical for the      //\n")
	helpers.Logger.Printf("[Initialization] // function of IRC over SSL.                          //\n")
	helpers.Logger.Printf("[Initialization] // If you wish to only offer unencrypted SSL, you can //\n")
	helpers.Logger.Printf("[Initialization] // ignore this warning.                               //\n")
	helpers.Logger.Printf("[Initialization] //                                                    //\n")
	helpers.Logger.Printf("[Initialization] // ssl_key:   Location of the Private Key File.       //\n")
	helpers.Logger.Printf("[Initialization] // ssl_cert:  Location of the SSL Certificate.        //\n")
	helpers.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")

	displayedWarnings["ssl"] = true
}

func IrcIpMissing() {
	if HostIrc == "false" {
		return
	}

	_, runAlready := displayedWarnings["irc_ip"]

	if runAlready {
		return
	}

	helpers.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")
	helpers.Logger.Printf("[Initialization] //////////////////  Attention!!!!!!!  //////////////////\n")
	helpers.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")
	helpers.Logger.Printf("[Initialization] // The IRC IP in the .env file has been left empty    //\n")
	helpers.Logger.Printf("[Initialization] // If you wish to not run a IRC server, you can       //\n")
	helpers.Logger.Printf("[Initialization] // set host_irc to false, which will disable IRC and  //\n")
	helpers.Logger.Printf("[Initialization] // silence this warning.                              //\n")
	helpers.Logger.Printf("[Initialization] //                                                    //\n")
	helpers.Logger.Printf("[Initialization] // irc_ip: Where the IRC TCP Listener is hosted       //\n")
	helpers.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")

	displayedWarnings["irc_ip"] = true
}

func IrcSSLIpMissing() {
	if HostIrcSsl == "false" {
		return
	}

	_, runAlready := displayedWarnings["irc_ssl_ip"]

	if runAlready {
		return
	}

	helpers.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")
	helpers.Logger.Printf("[Initialization] //////////////////  Attention!!!!!!!  //////////////////\n")
	helpers.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")
	helpers.Logger.Printf("[Initialization] // The IRC/SSL IP in the .env file has been left      //\n")
	helpers.Logger.Printf("[Initialization] // empty. if you wish to not run a IRC server, you    //\n")
	helpers.Logger.Printf("[Initialization] // can set host_irc_ssl to false, which will disable  //\n")
	helpers.Logger.Printf("[Initialization] // IRC/SSL and silence this warning.                  //\n")
	helpers.Logger.Printf("[Initialization] // Not providing certificate locations will also      //\n")
	helpers.Logger.Printf("[Initialization] // automatically disable IRC/SSL, although it won't   //\n")
	helpers.Logger.Printf("[Initialization] // silence this warning.                              //\n")
	helpers.Logger.Printf("[Initialization] //                                                    //\n")
	helpers.Logger.Printf("[Initialization] // irc_ssl_ip: Where IRC/SSL is hosted                //\n")
	helpers.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")

	displayedWarnings["irc_ssl_ip"] = true
}

func WaffleWebConfigMissing() {
	if UsingWaffleWeb == "false" {
		return
	}

	_, runAlready := displayedWarnings["using_waffle_web"]

	if runAlready {
		return
	}

	helpers.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")
	helpers.Logger.Printf("[Initialization] //////////////////  Attention!!!!!!!  //////////////////\n")
	helpers.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")
	helpers.Logger.Printf("[Initialization] // You have not specified whether you're running      //\n")
	helpers.Logger.Printf("[Initialization] // waffle-web or not, this is important to the        //\n")
	helpers.Logger.Printf("[Initialization] // function of the Beatmap Submission system, as      //\n")
	helpers.Logger.Printf("[Initialization] // forum posts will not be created when not using     //\n")
	helpers.Logger.Printf("[Initialization] // waffle-web, and the ranking queue will not be run  //\n")
	helpers.Logger.Printf("[Initialization] // without waffle-web.                                //\n")
	helpers.Logger.Printf("[Initialization] // silence this warning.                              //\n")
	helpers.Logger.Printf("[Initialization] //                                                    //\n")
	helpers.Logger.Printf("[Initialization] // using_waffle_web: Whether the original waffle-web  //\n")
	helpers.Logger.Printf("[Initialization] //                   is used, as opposed to a custom  //\n")
	helpers.Logger.Printf("[Initialization] //                   made frontend to go around it.   //\n")
	helpers.Logger.Printf("[Initialization] //                   it used around beatmap forum     //\n")
	helpers.Logger.Printf("[Initialization] //                   posts and ranking queue.         //\n")
	helpers.Logger.Printf("[Initialization] //                                                    //\n")
	helpers.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")

	displayedWarnings["using_waffle_web"] = true
}

func FfmpegMissing() {
	if HostIrcSsl == "false" {
		return
	}

	_, runAlready := displayedWarnings["ffmpeg_executable"]

	if runAlready {
		return
	}

	helpers.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")
	helpers.Logger.Printf("[Initialization] //////////////////  Attention!!!!!!!  //////////////////\n")
	helpers.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")
	helpers.Logger.Printf("[Initialization] // FFMPEG executable path unset!                      //\n")
	helpers.Logger.Printf("[Initialization] // Beatmaps uploaded by the Beatmap Submission System //\n")
	helpers.Logger.Printf("[Initialization] // will not be generating mp3 previews heard on the   //\n")
	helpers.Logger.Printf("[Initialization] // Website and inside osu!direct.                     //\n")
	helpers.Logger.Printf("[Initialization] //                                                    //\n")
	helpers.Logger.Printf("[Initialization] // ffmpeg_executable: Path to the ffmpeg executable   //\n")
	helpers.Logger.Printf("[Initialization] //                    including extension on Windows  //\n")
	helpers.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")

	displayedWarnings["ffmpeg_executable"] = true
}
