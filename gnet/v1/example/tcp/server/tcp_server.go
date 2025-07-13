package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pkg/pool/goroutine"
	"log"
	"time"
)

// @Title tcp_server.go
// @Author Zero - 2024/8/11 09:57:54
// @Description 使用 Gnet 实现的 TCP Server

// 服务端事件处理器
type echoServerEventHandle struct {
	*gnet.Server
	pool *goroutine.Pool
}

// OnInitComplete 服务器初始化完成事件
func (es *echoServerEventHandle) OnInitComplete(server gnet.Server) (action gnet.Action) {
	fd, _ := server.DupFd()
	fmt.Printf("[ServerInitComplete] Tcp server is listening on %s threads-core: %v event-loop: %d fd: %d \n",
		server.Addr.String(), server.Multicore, server.NumEventLoop, fd)
	return
}

// OnShutdown 服务器关闭事件
func (es *echoServerEventHandle) OnShutdown(server gnet.Server) {
	fmt.Printf("[ServerShutdown] Tcp server shutdown %s \n", server.Addr.String())
	return
}

// OnOpened 连接建立事件
func (es *echoServerEventHandle) OnOpened(conn gnet.Conn) (out []byte, action gnet.Action) {
	fmt.Printf("[ConnectionOpen] Accept new connection from %s \n", conn.RemoteAddr().String())
	return
}

// OnClosed 连接关闭事件
func (es *echoServerEventHandle) OnClosed(conn gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		fmt.Printf("[Close] Connection %s closed with error: %v\n", conn.RemoteAddr().String(), err)
	} else {
		fmt.Printf("[Close] Connection %s closed. \n", conn.RemoteAddr().String())
	}
	return
}

// PreWrite 数据写回客户端之前回调
func (es *echoServerEventHandle) PreWrite(conn gnet.Conn) {
	fmt.Printf("[PreWrite] Preparing to write to %s\n", conn.RemoteAddr().String())
	return
}

// AfterWrite 数据写回客户端之后回调
func (es *echoServerEventHandle) AfterWrite(conn gnet.Conn, b []byte) {
	fmt.Printf("[AfterWrite] Data written to %s, bytes: %d\n", conn.RemoteAddr().String(), len(b))
	return
}

// React 收到客户端数据事件
func (es *echoServerEventHandle) React(in []byte, conn gnet.Conn) (out []byte, action gnet.Action) {
	// 将数据转换为大写后写回
	line := string(in)
	if line == "squit\n" {
		action = gnet.Close
		return
	}
	fmt.Printf("[React] %s: %s", conn.RemoteAddr().String(), string(in))
	out = bytes.ToUpper(in)
	return
}

// Tick 定时任务事件
func (es *echoServerEventHandle) Tick() (delay time.Duration, action gnet.Action) {
	fmt.Printf("[Tick] Tick triggered at %s \n", time.Now().Format("2006-01-02 15:04:05"))
	delay = time.Second * 30 // 每 30 秒回调一次
	return
}

var port int
var multicore bool

func init() {
	flag.IntVar(&port, "p", 9090, "--p {port}")
	flag.BoolVar(&multicore, "m", false, "--m true")
}

func main() {
	flag.Parse()
	handle := &echoServerEventHandle{pool: goroutine.Default()}
	addr := fmt.Sprintf("tcp://0.0.0.0:%d", port)
	if err := gnet.Serve(handle, addr, gnet.WithMulticore(multicore), gnet.WithTicker(true)); err != nil {
		log.Fatal(err)
	}

}
