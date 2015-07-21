# tgoremote
Provide RPC handlers to and from regular Go and TARDISgo.

A work in progress proof-of-concept. 

Presents a sub-set of the go RPC API, to make remote calls from TARDISgo to normal go, and from normal go to TARDISgo (reverse RPC).

Simple use is as a component in a Go web-server, built on top of the haxeremote package (but only using the Haxe string serialization method). 

Complex use is linked with generated TARDISgo code on the client device. This uses the same protocol, but in memory.

There is significant work to do to make this code production ready.


