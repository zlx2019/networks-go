// 基于 GNet v2 的 TCP server
// 实现简单的 echo 功能

package main

import (
	"bytes"
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"log"
	"strings"
)

type echoServerHandle struct {
	gnet.BuiltinEventEngine
}

// OnBoot 服务端引擎就绪事件
func (es *echoServerHandle) OnBoot(eng gnet.Engine) gnet.Action {
	fmt.Println("Service is ready, starting to accept connections")
	return gnet.None
}

// OnShutdown 服务引擎关闭事件
func (es *echoServerHandle) OnShutdown(eng gnet.Engine) {
	fmt.Println("Service closed.")
}

// OnOpen 新的客户端连接事件
func (*echoServerHandle) OnOpen(conn gnet.Conn) (out []byte, action gnet.Action) {
	fmt.Printf("Accept connection from %s\n", conn.RemoteAddr().String())
	out = []byte("Welcome!\n")
	return
}

// OnClose 客户端关闭连接
func (es *echoServerHandle) OnClose(conn gnet.Conn, err error) (action gnet.Action) {
	fmt.Printf("Connection closed from %s\n", conn.RemoteAddr().String())
	return
}

// OnTraffic 数据可读事件
func (es *echoServerHandle) OnTraffic(c gnet.Conn) (action gnet.Action) {
	buf, err := c.Next(-1)
	if err != nil {
		fmt.Printf("conn read error: %v \n", err)
		_ = c.Close()
		return
	}
	input := string(buf)
	if strings.TrimSuffix(input, "\n") == "exit" {
		return gnet.Close
	}
	_, _ = c.Write(bytes.ToUpper(buf))
	return
}

func main() {
	server := &echoServerHandle{}
	if err := gnet.Run(server, "tcp://:9501", gnet.WithMulticore(true)); err != nil {
		log.Fatal(err)
	}
}
