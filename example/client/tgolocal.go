package main

import "C"
import "github.com/tardisgo/tgoremote/tgolocal"

//export TgoCall
func TgoCall(cstr *C.char) *C.char {
	// TODO deal with C.string garbage
	return C.CString(tgolocal.TgoCall(C.GoString(cstr)))
}
