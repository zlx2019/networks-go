package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/panjf2000/gnet/v2"
)

type echoServer struct {
	gnet.BuiltinEventEngine
}

func (e *echoServer) OnTraffic(conn gnet.Conn) gnet.Action {
	buf, _ := conn.Next(-1)
	n, _ := conn.Write(bytes.ToUpper(buf))
	fmt.Printf("echo: %s \n", string(buf[:n]))
	return gnet.None
}

func main() {
	var path string
	var multicore bool
	flag.StringVar(&path, "s", "echo.sock", "unix socket file path")
	flag.BoolVar(&multicore, "m", false, "enable multicore")
	flag.Parse()
	if err := gnet.Run(new(echoServer), fmt.Sprintf("unix://%s", path), gnet.WithMulticore(multicore)); err != nil {
		panic(err)
	}
}
