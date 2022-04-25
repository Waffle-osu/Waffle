package main

import (
	"Waffle/bancho"
	"time"
)

func main() {
	bancho := bancho.CreateBancho()
	go bancho.RunBancho()

	for {
		time.Sleep(2 * time.Second)
	}
}
