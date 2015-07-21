package main

import (
	"log"
	"net/http"

	"github.com/tardisgo/tgoremote/example"
	"github.com/tardisgo/tgoremote/tgohttp"
)

func main() {
	tgohttp.Setup("/_haxeRPC_")
	http.Handle("/client/", http.StripPrefix("/client/", http.FileServer(http.Dir("../client"))))
	example.Register()
	go example.RRPC()
	println("Tgo remote webserver running on port 8088")
	log.Fatal(http.ListenAndServe("localhost:8088", nil))
}
