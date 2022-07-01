package irc_messages

import (
	"fmt"
	"strings"
)

type IRCCode int32

type Message struct {
	Source       string
	Command      string
	NumCommand   IRCCode
	Params       []string
	Trailing     string
	SkipUsername bool
}

const (
	ErrNoSuchNick        IRCCode = 401
	ErrNoSuchServer      IRCCode = 402
	ErrNoSuchChannel     IRCCode = 403
	ErrCannotSendToChan  IRCCode = 404
	ErrToManyChannels    IRCCode = 405
	ErrWasNoSuchNick     IRCCode = 406
	ErrTooManyTargets    IRCCode = 407
	ErrNoOrigin          IRCCode = 409
	ErrNoRecipient       IRCCode = 411
	ErrNoTextToSend      IRCCode = 412
	ErrNoTopLevel        IRCCode = 413
	ErrWildTopLevel      IRCCode = 414
	ErrUnknownCommand    IRCCode = 421
	ErrNoMotd            IRCCode = 422
	ErrNoAdminInfo       IRCCode = 423
	ErrFileError         IRCCode = 424
	ErrNoNicknameGiven   IRCCode = 431
	ErrErroneusNickname  IRCCode = 432
	ErrNicknameInUse     IRCCode = 433
	ErrNickCollision     IRCCode = 436
	ErrUserNotInChannel  IRCCode = 441
	ErrNotOnChannel      IRCCode = 442
	ErrUserOnChannel     IRCCode = 443
	ErrNoLogin           IRCCode = 444
	ErrSummonDisabled    IRCCode = 445
	ErrUsersDisabled     IRCCode = 446
	ErrNotRegistered     IRCCode = 451
	ErrNeedMoreParams    IRCCode = 461
	ErrAlreadyRegistered IRCCode = 462
	ErrNoPermForHost     IRCCode = 463
	ErrPasswdMissmatch   IRCCode = 464
	ErrYoureBannedCreep  IRCCode = 465
	ErrKeySet            IRCCode = 467
	ErrChannelIsFull     IRCCode = 471
	ErrUnknownMode       IRCCode = 472
	ErrInviteOnlyChan    IRCCode = 473
	ErrBannedFromChan    IRCCode = 474
	ErrBadChannelKey     IRCCode = 475
	ErrNoPrivileges      IRCCode = 481
	ErrChanNoPrivsNeeded IRCCode = 482
	ErrCantKillServer    IRCCode = 483
	ErrNoOperHost        IRCCode = 491
	ErrUModeUnknownFlag  IRCCode = 501
	ErrUsersDontMatch    IRCCode = 502

	RplNone            IRCCode = 300
	RplAway            IRCCode = 301
	RplUserHost        IRCCode = 302
	RplIson            IRCCode = 303
	RplUnaway          IRCCode = 305
	RplNowArray        IRCCode = 306
	RplWhoIsUser       IRCCode = 311
	RplWhoIsServer     IRCCode = 312
	RplWhoIsOperator   IRCCode = 313
	RplWhoWasUser      IRCCode = 314
	RplWhoIsIdle       IRCCode = 317
	RplEndOfWhoIs      IRCCode = 318
	RplWhoIsChannels   IRCCode = 319
	RplEndOfWhoWas     IRCCode = 369
	RplListStart       IRCCode = 321
	RplList            IRCCode = 322
	RplListEnd         IRCCode = 323
	RplChannelModeIs   IRCCode = 324
	RplNoTopic         IRCCode = 331
	RplTopic           IRCCode = 332
	RplInviting        IRCCode = 341
	RplSummoning       IRCCode = 342
	RplVersion         IRCCode = 351
	RplWhoReply        IRCCode = 352
	RplEndOfWho        IRCCode = 315
	RplNameReply       IRCCode = 353
	RplEndOfNames      IRCCode = 366
	RplLinks           IRCCode = 364
	RplEndOfLinks      IRCCode = 365
	RplBanList         IRCCode = 367
	RplEndOfBanList    IRCCode = 368
	RplInfo            IRCCode = 371
	RplEndOfInfo       IRCCode = 374
	RplMotdStart       IRCCode = 375
	RplMotd            IRCCode = 372
	RplEndOfMotd       IRCCode = 376
	RplYoureOper       IRCCode = 381
	RplRehashing       IRCCode = 382
	RplTime            IRCCode = 391
	RplUsersStart      IRCCode = 392
	RplUsers           IRCCode = 393
	RplEndOfUsers      IRCCode = 394
	RplNoUsers         IRCCode = 395
	RplTraceLink       IRCCode = 200
	RplTraceConnecting IRCCode = 201
	RplTraceHandshake  IRCCode = 202
	RplTraceUnknwon    IRCCode = 203
	RplTraceOperator   IRCCode = 204
	RplTraceUser       IRCCode = 205
	RplTraceServer     IRCCode = 206
	RplTraceNewType    IRCCode = 208
	RplTraceLog        IRCCode = 261
	RplStatsLinkInfo   IRCCode = 211
	RplStatsCommands   IRCCode = 212
	RplStatsCLine      IRCCode = 213
	RplStatsNLine      IRCCode = 214
	RplStatsILine      IRCCode = 215
	RplStatsKLine      IRCCode = 216
	RplStatsYLine      IRCCode = 218
	RplEndOfStats      IRCCode = 219
	RplStatsLLine      IRCCode = 241
	RplStatsUptime     IRCCode = 242
	RplStatsOLine      IRCCode = 243
	RplStatsHLine      IRCCode = 244
	RplUModeIs         IRCCode = 221
	RplLUserClient     IRCCode = 251
	RplLUserOp         IRCCode = 252
	RplLUserUnknown    IRCCode = 253
	RplLUserChannels   IRCCode = 254
	RplLUserMe         IRCCode = 255
	RplAdminMe         IRCCode = 256
	RplAdminLoc1       IRCCode = 257
	RplAdminLoc2       IRCCode = 258
	RplAdminEmail      IRCCode = 259
)

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

func (message Message) FormatMessage(username string) (formatted string, formatErr string) {
	if len(message.Params) == 0 && message.Trailing == "" {
		return "", "Either Parameters or Trailing has to be set!"
	}

	source := "irc.waffle.nya"

	if message.Source != "" {
		source = message.Source
	}

	returnString := ""

	if message.Command == "" {
		returnString += fmt.Sprintf(":%s %03d", source, message.NumCommand)

		if !message.SkipUsername {
			returnString = fmt.Sprintf("%s %s", returnString, username)
		}
	} else {
		returnString = fmt.Sprintf(":%s %s", source, message.Command)

		if !message.SkipUsername {
			returnString = fmt.Sprintf("%s %s", returnString, username)
		}
	}

	if len(message.Params) != 0 {
		returnString = fmt.Sprintf("%s %s", returnString, strings.Join(message.Params, " "))
	}

	if len(message.Trailing) != 0 {
		returnString = fmt.Sprintf("%s :%s", returnString, message.Trailing)
	}

	return returnString + "\n", ""
}
