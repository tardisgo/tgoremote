package main

// #include <stdio.h>
// #include <stdlib.h>
import "C"
import "github.com/tardisgo/tgoremote/tgolocal"
import "unsafe"

//export TgoCall
func TgoCall(cstr *C.char) *C.char {
	return C.CString(tgolocal.TgoCall(C.GoString(cstr)))
}

//export TgoFree
func TgoFree(cstr *C.char) {
	// called to deal with C.CString allocation created above
	C.free(unsafe.Pointer(cstr))
}
