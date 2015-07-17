package main

/*
// NOTE libraries below are for the test mac, your libraries WILL be different!
// TODO make it easier to use!

#cgo LDFLAGS: -stdlib=libstdc++ tgo/tardis/cpplib/libGo.a /usr/lib/haxe/lib/hxcpp/3,2,102/lib/Mac64/libstd.a /usr/lib/haxe/lib/hxcpp/3,2,102/lib/Mac64/libzlib.a /usr/lib/haxe/lib/hxcpp/3,2,102/lib/Mac64/libregexp.a
// /usr/lib/haxe/lib/hxcpp/3,2,102/lib/Mac64/libsqlite.a /usr/lib/haxe/lib/hxcpp/3,2,102/lib/Mac64/libmysql5.a
extern int hxmain();
*/
import "C"
import "github.com/tardisgo/tgoremote/example"

func main() {
	example.Register()
	C.hxmain()
}
