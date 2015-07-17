// Package tgorpc provides the TARDISgo end of the RPC connection
package tgorpc

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"net/url"
	"sync"

	"github.com/tardisgo/tardisgo/haxe/hx"
)

func init() {
	hx.Source("TgoConnect", `
package tardis;

#if (cpp&&static_link)
	   @:cppFileCode('extern "C" char *TgoCall(char *args);')
#end

class TgoConnect {
	   #if !(cpp&&static_link)
	   	var cnx:haxe.remoting.HttpAsyncConnection;
	   #end
	   public function new(url:String){
	   	#if !(cpp&&static_link)
	   		cnx = haxe.remoting.HttpAsyncConnection.urlConnect(url);
	   	#end
	   }
	   public function setErrorHandler(errFn:Dynamic->Void):Void{
	   	#if !(cpp&&static_link)
	   		cnx.setErrorHandler(errFn);
	   	#end
	   }
	   public function call(args:Array<Dynamic>,result:Dynamic->Void):Void
	   {
	   	#if (cpp&&static_link)
	   		//trace("DEBUG LocCon args", args);
	   		var eargs = haxe.Serializer.run(args);
	   		var cstr = cpp.NativeString.c_str(eargs);
	   		var rcstr:cpp.ConstPointer<cpp.Char> = untyped __cpp__("TgoCall(cstr)");
	   		var rstr = cpp.NativeString.fromPointer(rcstr);
	   		var res = haxe.Unserializer.run(rstr);
	   		//trace("DEBUG LocCon res", res);
	   		result(res);
	   	#else
	   		cnx._TgoRPC_.call(args,result);
	   	#end
	   }
}
`)
}

type Conn struct {
	conn  uintptr
	mutex sync.Mutex // one-at-a-time because error handing is per conn
}

type Call struct {
	ServiceMethod string      // The name of the service and method to call.
	Args          interface{} // The argument to the function (*struct).
	Reply         interface{} // The reply from the function (*struct).
	Error         error       // After completion, the error status.
	Done          chan *Call  // Strobes when call is complete.
}

func Dial(endpoint string) *Conn {
	return &Conn{
		conn: hx.New("", "TgoConnect", 1, endpoint),
	}
}

func (c Conn) Go(serviceMethod string, args, reply interface{}, done chan *Call) *Call {
	call := &Call{
		ServiceMethod: serviceMethod,
		Args:          args,
		Reply:         reply,
		Error:         nil,
		Done:          done,
	}
	c.mutex.Lock()
	hx.Meth("", c.conn, "TgoConnect", "setErrorHandler", 1, func(e uintptr) {
		call.Error = errors.New(hx.CallString("", "Std.string", 1, e))
		done <- call
		c.mutex.Unlock()
	})
	//fmt.Printf("DEBUG args=%v:%T\n", args, args)
	var networkOut bytes.Buffer
	enc := gob.NewEncoder(&networkOut)
	err := enc.Encode(args)
	if err != nil {
		call.Error = err
		done <- call
		c.mutex.Unlock()
		return call
	}
	var networkOut2 bytes.Buffer
	enc2 := gob.NewEncoder(&networkOut2)
	err = enc2.Encode(reply)
	if err != nil {
		call.Error = err
		done <- call
		c.mutex.Unlock()
		return call
	}

	ifa := []interface{}{serviceMethod,
		base64.StdEncoding.EncodeToString(networkOut.Bytes()),
		base64.StdEncoding.EncodeToString(networkOut2.Bytes())}

	//fmt.Printf("DEBUG ifa %#v:%T\n", ifa, ifa)

	haxeArgs := interfaceToDynamic(ifa)
	//fmt.Println("DEBUG haxe ifa", hx.CallString("", "Std.string", 1, haxeArgs))

	hx.Meth("", c.conn, "TgoConnect", "call", 2, haxeArgs, func(r uintptr) {
		c.mutex.Unlock() // another call can begin before processing what's below
		rA, ok := dynamicToInterface(r).([]interface{})
		if !ok {
			call.Error = errors.New("returned value not an []interface{} in tgocall")
		} else {
			errMsg, ok := rA[1].(string)
			if !ok {
				call.Error = errors.New("returned error message not a string in tgocall")
			} else {
				if errMsg != "" {
					call.Error = errors.New(errMsg)
				} else {
					back64, ok := rA[0].(string)
					if !ok {
						call.Error = errors.New("returned encoded data not a string in tgocall")
					} else {
						backBuf, err64 := base64.StdEncoding.DecodeString(back64)
						if err64 != nil {
							call.Error = err64
						} else {
							networkBack := bytes.NewBuffer(backBuf)
							dec := gob.NewDecoder(networkBack)
							err = dec.Decode(call.Reply)
							if err != nil {
								call.Error = err
							}
						}
					}
				}
			}
		}
		done <- call
	})
	return call
}

func (c Conn) Call(serviceMethod string, args, reply interface{}) error {
	call := <-c.Go(serviceMethod, args, reply, make(chan *Call, 1)).Done
	return call.Error
}

// interfaceToDynamic - limited runtime conversion of Go->Haxe types
// TODO consider making public
func interfaceToDynamic(a interface{}) uintptr {
	//fmt.Printf("DEBUG interfaceToDynamic a= %v:%T\n", a, a)
	if a == nil {
		return hx.Null()
	}
	switch a.(type) {
	case []interface{}:
		//println("DEBUG []interface{}")
		ret := hx.New("", "Array<Dynamic>", 0)
		for _, aa := range a.([]interface{}) {
			//fmt.Printf("DEBUG aa= %v:%T\n", aa, aa)
			val := interfaceToDynamic(aa)
			//fmt.Println("DEBUG val=" + hx.CallString("", "Std.string", 1, val))
			hx.Code("", "_a.param(0).val.push(_a.param(1).val);", ret, val)
		}
		return ret
	case bool, string, int, float64:
		return hx.CodeDynamic("", "_a.param(0).val;", a)
	case []byte:
		return hx.CodeDynamic("", "Slice.toBytes(cast(_a.param(0).val,Slice));", a)
	default:
		panic("Unhandled Go interface{} to Haxe Dynamic: " + hx.CallString("", "Std.string", 1, a))
	}
	return hx.Null()
}

// dynamicToInterface - limited runtime conversion of Haxe->Go types
// with URL un-escaping of strings
// TODO consider making public
func dynamicToInterface(dyn uintptr) interface{} {
	switch {
	case hx.CodeBool("", "Std.is(_a.param(0).val,Array);", dyn):
		l := hx.CodeInt("", "_a.param(0).val.length;", dyn)
		ret := make([]interface{}, l)
		for i := 0; i < l; i++ {
			ret[i] =
				dynamicToInterface(hx.CodeDynamic("",
					"_a.param(0).val[_a.param(1).val];", dyn, i))
		}
		return ret
	case hx.IsNull(dyn):
		return nil
	case hx.CodeBool("", "Std.is(_a.param(0).val,Bool);", dyn):
		return hx.CodeBool("", "_a.param(0).val;", dyn)
	case hx.CodeBool("", "Std.is(_a.param(0).val,Int);", dyn):
		return hx.CodeInt("", "_a.param(0).val;", dyn)
	case hx.CodeBool("", "Std.is(_a.param(0).val,Float);", dyn):
		return hx.CodeFloat("", "_a.param(0).val;", dyn)
	case hx.CodeBool("", "Std.is(_a.param(0).val,String);", dyn):
		raw := hx.CodeString("", "_a.param(0).val;", dyn)
		clean, err := url.QueryUnescape(raw)
		if err == nil {
			return clean
		}
		return raw
	case hx.CodeBool("", "Std.is(_a.param(0).val,haxe.io.Bytes);", dyn):
		return hx.CodeIface("", "[]byte",
			"Slice.fromBytes(_a.param(0).val);", dyn)
	default:
		panic("unhandled haxe Dynamic to Go interface{} :" +
			hx.CallString("", "Std.string", 1, dyn))
	}
	return nil
}
