package main

import (
	"log"
	"net/http"

	"github.com/tardisgo/tgoremote/example"
	"github.com/tardisgo/tgoremote/tgohttp"
)

func main() {
	example.Register()

	tgohttp.Setup()
	http.Handle("/client/", http.StripPrefix("/client/", http.FileServer(http.Dir("../client"))))

	println("Tgo remote webserver running on port 8088")
	log.Fatal(http.ListenAndServe("localhost:8088", nil))
}
