package tgohttp

import (
	"net/http"

	"github.com/tardisgo/haxeremote"
)
import "github.com/tardisgo/tgoremote"

func Setup() {
	haxeremote.AddFunc("_TgoRPC_", tgoremote.CallFunc)
	http.HandleFunc("/_haxeRPC_", haxeremote.HttpHandler)
}
