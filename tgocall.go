package tgoremote

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"
)

var registry = make(map[string]interface{})
var registryMutex sync.Mutex

func RegisterName(name string, rcvr interface{}) error {
	//fmt.Printf("DEBUG %s %T\n", name, rcvr)
	registryMutex.Lock()
	registry[name] = rcvr
	registryMutex.Unlock()
	return nil
}

func Register(rcvr interface{}) error {
	name := reflect.TypeOf(rcvr).Name()
	if name == "" {
		return errors.New("un-named type")
	}
	RegisterName(name, rcvr)
	return nil
}

// CallFunc is called from tgohttp, tgolocal
// arg[0] string containing the serviceMethod
// arg[1] string containing the base64 and gob-encoded args
// arg[2] string containing the base64 and gob-encoded reply
// returns:
// []interface{}[0] string containing the base64 and gob-encoded reply
// []interface{}[1] error string
func CallFunc(arg interface{}) interface{} {
	argA, ok := arg.([]interface{})
	if !ok {
		msg := fmt.Sprintf(
			"the arguments for CallFunc() are not []interface{} = %v : %T",
			arg, arg)
		log.Println(msg)
		return []interface{}{"", msg}
	}
	name, ok := (argA[0]).(string)
	if !ok {
		msg := "the expected remote function name is not a string for CallFunc()"
		log.Println(msg)
		return []interface{}{"", msg}
	}
	//fmt.Printf("DEBUG un-decoded args %v:%T\n", argA[1], argA[1])
	argDat, err := base64.StdEncoding.DecodeString(argA[1].(string))
	if err != nil {
		msg := err.Error()
		log.Println(msg)
		return []interface{}{"", msg}
	}

	//fmt.Printf("DEBUG un-decoded return %v:%T\n", argA[2], argA[2])
	retDat, err := base64.StdEncoding.DecodeString(argA[2].(string))
	if err != nil {
		msg := err.Error()
		log.Println(msg)
		return []interface{}{"", msg}
	}

	reply, errMsg := CallFuncDirect(name, argDat, retDat)

	if errMsg != "" {
		log.Println(errMsg)
		return []interface{}{"", errMsg}
	}
	return []interface{}{base64.StdEncoding.EncodeToString(reply), ""}

}

func CallFuncDirect(name string, argDat, retDat []byte) ([]byte, string) {
	parts := strings.Split(name, ".")
	if len(parts) != 2 {
		msg := "RPC serviceMethod should be of the form A.B rather than: " + name
		log.Println(msg)
		return []byte{}, msg
	}
	registryMutex.Lock()
	obj, ok := registry[parts[0]]
	registryMutex.Unlock()
	if !ok {
		msg := "could not find object " + parts[0] + " in CallFunc()"
		log.Println(msg)
		return []byte{}, msg
	}
	//fmt.Println("DEBUG found object", parts[0])
	meth := reflect.ValueOf(obj).MethodByName(parts[1])
	if meth == (reflect.Value{}) {
		msg := "could not find method " + parts[0] + "." + parts[1] + " in CallFunc()"
		log.Println(msg)
		return []byte{}, msg
	}
	//fmt.Println("DEBUG found method", parts[1])

	argBuf := bytes.NewBuffer(argDat)
	argDec := gob.NewDecoder(argBuf)

	argVal := reflect.New(meth.Type().In(0)).Interface()
	err := argDec.Decode(argVal)
	if err != nil {
		msg := err.Error()
		log.Println(msg)
		return []byte{}, msg
	}
	//fmt.Printf("DEBUG decoded args %v:%T\n", argVal, argVal)
	retBuf := bytes.NewBuffer(retDat)
	retDec := gob.NewDecoder(retBuf)
	retVal := reflect.New(meth.Type().In(1)).Interface()
	err = retDec.Decode(retVal)
	if err != nil {
		msg := err.Error()
		log.Println(msg)
		return []byte{}, msg
	}
	//fmt.Printf("DEBUG decoded return %v:%T\n", retVal, retVal)

	rets := meth.Call(
		[]reflect.Value{
			reflect.ValueOf(argVal).Elem(),
			reflect.ValueOf(retVal).Elem(),
		})
	//fmt.Printf("DEBUG return error %v:%T\n", err, err)
	if len(rets) > 0 {
		err, ok = rets[0].Interface().(error)
		if ok {
			msg := err.Error()
			log.Println(msg)
			return []byte{}, msg
		}
	}

	var results bytes.Buffer
	resEnc := gob.NewEncoder(&results)
	err = resEnc.Encode(retVal)
	if err != nil {
		msg := err.Error()
		log.Println(msg)
		return []byte{}, msg
	}

	return results.Bytes(), ""
}
