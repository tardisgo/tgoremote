package example

import (
	"bytes"
	"fmt"
)

import "github.com/tardisgo/tgoremote"

import "github.com/tardisgo/tgoremote/rrpcclient"

type Tfoo struct {
	A, B int
}

type TfooBase struct{}

func (tfb TfooBase) Add(args Tfoo, ans *int) error {
	*ans = args.A + args.B
	println("DEBUG tfoo", args.A, args.B, *ans)
	return nil
}

type Tbar struct {
	A, B, C string
}
type TbarBase struct{}

func (tbb TbarBase) Join(args Tbar, ans *string) error {
	*ans = args.A + " " + args.B + args.C
	println("DEBUG tbar", args.A, args.B, args.C, *ans)
	return nil
}

type Tdoodah struct {
	A []byte
}

type TdoodahBase struct{}

func (tddb TdoodahBase) Upper(args Tdoodah, ans *[]byte) error {
	*ans = bytes.ToUpper(args.A)
	println("DEBUG tdoodah", string(args.A), string(*ans))
	return nil
}

var tfb TfooBase
var tbb TbarBase
var tddb TdoodahBase

func Register() {
	tgoremote.Register(tfb)
	tgoremote.Register(tbb)
	tgoremote.Register(tddb)
}

func RRPC() {
	var res float64
	rrpc := rrpcclient.Dial("_RRPC_")
	err := rrpc.Call("Float.Square", float64(4), &res)
	fmt.Println("YAY! Float.Square=", res, err)
	var ires int64
	err = rrpc.Call("Float.Square", int64(4), &ires)
	fmt.Println("NAY! Float.Square=", ires, err)
}
