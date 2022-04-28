package main

import (
	"Waffle/bancho"
	"Waffle/web"
	"os"
	"time"
)

func EnsureDirectoryExists(name string) {
	_, err := os.Stat(name)

	if err == nil {
		return
	}

	_ = os.Mkdir(name, os.ModeDir)
}

func main() {
	EnsureDirectoryExists("screenshots")

	go bancho.RunBancho()
	go web.RunOsuWeb()

	for {
		time.Sleep(2 * time.Second)
	}
}
