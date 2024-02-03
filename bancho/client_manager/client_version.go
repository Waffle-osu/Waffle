package client_manager

type ClientVersion int32

const (
	ClientVersionHiddenOsu ClientVersion = 0
	ClientVersionOsuB1815  ClientVersion = 1
	ClientVersionOsuIrc    ClientVersion = 2
	ClientVersionIrc       ClientVersion = 3
)
