/* Package tgolocal is run in a normal go environment

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
	results := tgoremote.CallFunc(args)
	reply := haxeremote.Serialize(results)
	//fmt.Printf("DEBUG TgoCall results: %v serialized-reply: %s\n", results, reply)
	return reply
}
