package example

import "bytes"
import "github.com/tardisgo/tgoremote"

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
