package main

import (
	"Waffle/waffle"
	"time"
)

func main() {
	bancho := waffle.CreateBancho()
	go bancho.RunBancho()

	for {
		time.Sleep(2 * time.Second)
	}
}
