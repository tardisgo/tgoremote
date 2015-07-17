# tgoremote
Provide RPC handlers for normal Go to be called by code written in TARDISgo.

A work in progress. Presents a sub-set of the go RPC API, to make remote calls from TARDISgo to normal go.

Simple use is as a component in a Go web-server, built on top of the haxeremote package (but only using the Haxe string serialization method). 

Complex use is linked with generated TARDISgo code on the client device. This uses the same protocol, but in memory. (In future it could be optimized to use a mechanism that does not require so much encoding/decoding.) See the 
