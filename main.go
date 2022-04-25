package main

import (
	"Waffle/bancho"
	"Waffle/web"
	"time"
)

func main() {
	go bancho.RunBancho()
	go web.RunOsuWeb()

	for {
		time.Sleep(2 * time.Second)
	}
}
