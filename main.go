package main

import (
	"Waffle/waffle"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	bancho := waffle.CreateBancho()
	go bancho.RunBancho()

	for {
		time.Sleep(2 * time.Second)
	}
}
