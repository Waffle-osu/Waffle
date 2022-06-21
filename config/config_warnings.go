package config

import "Waffle/helpers"

func MySqlSettingsIncompleteError() string {
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

	return "mysql"
}

func TokenFormatWarning() string {
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

	return "token_format"
}

func IrcSslCertsWarning() string {
	if SSLSilenceWarning == "true" {
		return "ssl"
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

	return "ssl"
}
