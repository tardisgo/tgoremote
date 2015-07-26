package tgohttp

import (
	"net/http"

	"github.com/tardisgo/haxeremote/hxrhttp"
	"github.com/tardisgo/tgoremote"
)

func Setup(haxeEndpoint string) {
	hxrhttp.AddFunc("_TgoRPC_", tgoremote.CallFunc)
	http.HandleFunc(haxeEndpoint, hxrhttp.HttpHandler)
}
