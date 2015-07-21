package tgohttp

import (
	"net/http"

	"github.com/tardisgo/haxeremote"
	"github.com/tardisgo/tgoremote"
)

func Setup(haxeEndpoint string) {
	haxeremote.AddFunc("_TgoRPC_", tgoremote.CallFunc)
	http.HandleFunc(haxeEndpoint, haxeremote.HttpHandler)
}
