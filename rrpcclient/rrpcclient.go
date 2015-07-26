// package rrpcclient runs in go proper, it provides the reverse RPC client
package rrpcclient

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"
	"sync"

	"github.com/tardisgo/tgoremote"
)

// Rrpc holds information for a reverse RPC
type Rrpc struct {
	mutex      sync.Mutex
	nextCallID uint32
	waiting    chan *Call
	processing map[uint32]*Call // TODO add timeout
}

// Ping is a request from the reverse RPC server for work.
func (r Rrpc) Ping(ping int, cmd *Call) {
	w := <-r.waiting
	*cmd = *w
	//fmt.Printf("DEBUG rRPCclient Ping Rrpc=%#v Call=%#v\n", r, cmd)
	r.mutex.Lock()
	r.processing[cmd.ID] = w
	r.mutex.Unlock()
}

// Pong is a response from the reverse RPC server
func (r Rrpc) Pong(reply Call, pong *int) {
	//fmt.Printf("DEBUG rRPCclient Pong Rrpc=%#v Call=%#v\n", r, reply)
	r.mutex.Lock()
	call, ok := r.processing[reply.ID]
	delete(r.processing, reply.ID)
	r.mutex.Unlock()
	if ok {
		call.ErrorMsg = reply.ErrorMsg
		if call.ErrorMsg == "" {
			networkBack := bytes.NewBuffer(reply.Reply)
			dec := gob.NewDecoder(networkBack)
			err := dec.Decode(call.replyIface)
			if err != nil {
				call.ErrorMsg = err.Error()
				call.Error = err
			}
		} else {
			call.Error = errors.New(call.ErrorMsg)
		}
		call.Done <- call
	} else {
		log.Println("Pong unknown call ID")
	}
}

// Dial an end-point for a reverse RPC connection.
func Dial(name string) *Rrpc {
	con := &Rrpc{
		waiting:    make(chan *Call, 10),
		processing: make(map[uint32]*Call),
	}
	tgoremote.RegisterName(name, con)
	return con
}

// Call provides the information to make and get results from a call.
type Call struct {
	Valid         bool        // Is this a valid call.
	ID            uint32      // ID of this call.
	ServiceMethod string      // The name of the service and method to call.
	Args          []byte      // Encoded arguments to the function (struct).
	Reply         []byte      // Encoded reply from the function (*struct).
	replyIface    interface{} // Where to put the result using gob.
	ErrorMsg      string      // After completion, the error status.
	Error         error       // Proper error for regular Go
	Done          chan *Call  // Strobes when call is complete.
}

// Go should have the same behaviour as net/rpc/Go().
func (r Rrpc) Go(serviceMethod string, args, reply interface{}, done chan *Call) *Call {
	r.mutex.Lock()
	r.nextCallID++
	thisID := r.nextCallID
	r.mutex.Unlock()

	if done == nil {
		done = make(chan *Call, 1)
	} else {
		// Comment from: http://golang.org/src/net/rpc/client.go
		// If caller passes done != nil, it must arrange that
		// done has enough buffer for the number of simultaneous
		// RPCs that will be using that channel.  If the channel
		// is totally unbuffered, it's best not to run at all.
		if cap(done) == 0 {
			log.Panic("rpc: done channel is unbuffered")
		}
	}

	call := &Call{
		Valid:         true,
		ID:            thisID,
		ServiceMethod: serviceMethod,
		replyIface:    reply,
		Done:          done,
	}
	ca := bytes.NewBuffer(call.Args)
	enc := gob.NewEncoder(ca)
	err := enc.Encode(args)
	if err != nil {
		call.ErrorMsg = err.Error()
		call.Error = err
		done <- call
		return call
	}
	call.Args = ca.Bytes()
	cr := bytes.NewBuffer(call.Reply)
	enc2 := gob.NewEncoder(cr)
	err = enc2.Encode(reply)
	if err != nil {
		call.ErrorMsg = err.Error()
		call.Error = err
		done <- call
		return call
	}
	call.Reply = cr.Bytes()
	//fmt.Printf("DEBUG c.waiting <- call Rrpc=%#v Call=%#v\n", r, call)
	r.waiting <- call

	return call
}

// Call should have the same behaviour as net/rpc/Call().
func (c Rrpc) Call(serviceMethod string, args, reply interface{}) error {
	call := <-c.Go(serviceMethod, args, reply, make(chan *Call, 1)).Done
	return call.Error
}
