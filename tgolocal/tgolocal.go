/*

usage from main package if running on client device:
//export TgoCall
func TgoCall(cstr *C.char) *C.char {
	return C.CString(tgolocal.TgoCall(C.GoString(cstr)))
}
*/
package tgolocal

import "C"
import "github.com/tardisgo/haxeremote"
import "github.com/tardisgo/tgoremote"

func TgoCall(serialized string) string {
	args, _, err := haxeremote.Unserialize([]byte(serialized))
	if err != nil {
		panic(err)
	}
	res := tgoremote.CallFunc(args)
	return haxeremote.Serialize(res)
}
