package main

import (
	"fmt"
	"haxegoruntime"

	"github.com/tardisgo/tgoremote"
	"github.com/tardisgo/tgoremote/rrpcserver"
	"github.com/tardisgo/tgoremote/tgorpc"
)

type Tfoo struct {
	A, B int
}

type Tbar struct {
	A, B, C string
}

type Tdoodah struct {
	A []byte
}

func main() {
	haxegoruntime.BrowserMain(func() {
		endpt := "http://localhost:8088/_haxeRPC_"

		conn := tgorpc.Dial(endpt)
		var tfooRes int
		err := conn.Call("TfooBase.Add", Tfoo{A: 1, B: 2}, &tfooRes)
		fmt.Println("tfoo:", tfooRes, err)
		var tbarRes string
		err = conn.Call("TbarBase.Join", Tbar{"it", "works", "!"}, &tbarRes)
		fmt.Println("tbar:", tbarRes, err)
		var tdoodahRes []byte
		err = conn.Call("TdoodahBase.Upper",
			Tdoodah{[]byte("The 世界 is my oyster!\n")}, &tdoodahRes)
		fmt.Println("tdoodah:", string(tdoodahRes), err)
		var tdoodahResBad []byte
		err = conn.Call("TdoodahBase.Upper",
			"The 世界 is my whelk!\n", &tdoodahResBad)
		fmt.Println("tdoodah bad call:", string(tdoodahResBad), err)

		var f Float
		tgoremote.Register(f)
		rrpcserver.Serve(endpt, "_RRPC_")

	}, 10 /* ms between each wake-up */, 10000 /* max scheduler runs per wake-up */)
}

type Float struct {
}

func (Float) Square(of float64, gives *float64) {
	*gives = of * of
	fmt.Println("DEBUG Square of ", of, " gives ", *gives)
}
