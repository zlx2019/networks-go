package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"log"
)

/// gNet v2 udp server

type echoServer struct {
	gnet.BuiltinEventEngine
}

func (es *echoServer) OnTraffic(conn gnet.Conn) gnet.Action {
	buf, _ := conn.Next(-1)
	_, _ = conn.Write(bytes.ToUpper(buf))
	return gnet.None
}

func main() {
	var port int
	var multiCore, reusePort bool
	flag.IntVar(&port, "p", 9501, "listen port")
	flag.BoolVar(&multiCore, "m", false, "enable multi core processing")
	flag.BoolVar(&reusePort, "r", false, "reuse port")
	flag.Parse()
	es := new(echoServer)
	log.Fatal(gnet.Run(es, fmt.Sprintf("udp://:%d", port), gnet.WithMulticore(multiCore), gnet.WithReusePort(reusePort)))
}
