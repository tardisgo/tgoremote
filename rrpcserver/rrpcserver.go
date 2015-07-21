// Package rrpcserver runs in TARDISgo,
// it provides a server for reverse RPC calls
package rrpcserver

import (
	"log"
	"runtime"

	"github.com/tardisgo/tgoremote"
	"github.com/tardisgo/tgoremote/rrpcclient"
	"github.com/tardisgo/tgoremote/tgorpc"
)

// Serve provides the server for reverse RPC calls.
// TODO allow for parallel execution and timout of calls.
func Serve(endpoint, name string) {
	conn := tgorpc.Dial(endpoint)
	for {
		var i int
		var cmd rrpcclient.Call
		//log.Println("DEBUG send Ping")
		err := conn.Call(name+".Ping", i, &cmd)
		//log.Printf("DEBUG rRPCserver Ping reply cmd=%#v\n", cmd)
		if err != nil {
			log.Printf("rRPCserver Ping Error=%s\ncmd=%#v\n", err.Error(), cmd)
		} else {
			if cmd.Valid {
				cmd.Reply, cmd.ErrorMsg = tgoremote.CallFuncDirect(
					cmd.ServiceMethod,
					cmd.Args,
					cmd.Reply,
				)
				err := conn.Call(name+".Pong", cmd, &i)
				if err != nil {
					log.Println("Pong " + err.Error())
				}
			} else {
				//log.Println("DEBUG Ping invalid reply")
			}
		}
		runtime.Gosched()
	}
}
