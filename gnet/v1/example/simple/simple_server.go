package main

import (
	"fmt"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pkg/pool/goroutine"
	"log"
	"log/slog"
	"time"
)

// @Title simple_server.go
// @Author Zero - 2024/8/11 09:57:54
// @Description 使用 gnet-v1 实现一个TCP服务器

type EchoServer struct {
	*gnet.EventServer
	pool *goroutine.Pool
}

// OnInitComplete 当 Server 初始化完成之后调用。
func (es *EchoServer) OnInitComplete(server gnet.Server) (action gnet.Action) {
	fmt.Printf("server initial complete. \n")
	return
}

// OnOpened 当新的连接被打开时回调
func (es *EchoServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	fmt.Printf("client connected to %s\n", c.RemoteAddr().String())
	return
}

// OnClosed 当连接关闭时回调
func (es *EchoServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	fmt.Printf("client closed %s\n", c.RemoteAddr().String())
	return
}

// Tick 定时事件, 服务器启动的时候会调用一次，之后每间隔 delay 调用一次。
func (es *EchoServer) Tick() (delay time.Duration, action gnet.Action) {
	slog.Info("server interval tick")
	delay = time.Second * 10
	return
}

// React EventHandler 有数据可读事件
func (es *EchoServer) React(in []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	fmt.Printf("client [%s]: %s", c.RemoteAddr().String(), string(in))
	//out = in
	data := append([]byte{}, in...)

	// 如果处理逻辑中可能会有IO阻塞，那么应该提交至协程池中执行
	_ = es.pool.Submit(func() {
		time.Sleep(time.Second)
		if err := c.AsyncWrite(data); err != nil {
			fmt.Printf("async write fail: %v\n", err)
		}
	})
	return
}

// PreWrite Server 将数据写回 Client 之前回调
// 将一些记录/ 计数/ 报告/ 报告码或任何前额操作列出。
func (es *EchoServer) PreWrite(c gnet.Conn) {

}

// AfterWrite Server 将数据写回 Client 之后回调
// 此事件功能通常是将 数据缓冲区 归还到内存池。
func (es *EchoServer) AfterWrite(c gnet.Conn, b []byte) {

}

func main() {
	server := &EchoServer{
		pool: goroutine.Default(),
	}
	log.Fatal(gnet.Serve(server, "tcp://:9000", gnet.WithMulticore(true), gnet.WithTicker(true)))
}
