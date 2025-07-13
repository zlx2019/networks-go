package main

import (
	"bufio"
	"fmt"
	"github.com/panjf2000/gnet"
	"os"
	"os/signal"
	"syscall"
)

// @Title tcp_client.go
// @Description
// @Author Zero - 2024/8/11 09:57:54

type clientHandle struct {
	*gnet.EventServer
}

// OnClosed 服务端连接关闭事件
func (ec *clientHandle) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	fmt.Printf("Server connection closed")
	sig <- syscall.SIGINT
	return
}

// React 服务端响应事件
func (ec *clientHandle) React(in []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	fmt.Printf("Server reply: %s \n", string(in))
	return
}

var sig = make(chan os.Signal, 1)

func main() {
	// 创建客户端
	client, err := gnet.NewClient(&clientHandle{})
	if err != nil {
		panic(err)
	}
	// 连接服务端
	conn, err := client.Dial("tcp", "127.0.0.1:9090")
	if err != nil {
		panic(err)
	}
	// 启动客户端
	err = client.Start()
	if err != nil {
		panic(err)
	}
	// 向服务端写入数据
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			line, _ := reader.ReadString('\n')
			if err != nil {
				break
			}
			// 退出客户端
			if line == "cquit\n" {
				_ = conn.Close()
				return
			}
			if err = conn.SendTo([]byte(line)); err != nil {
				fmt.Println("send err:", err)
			}
		}
	}()

	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sig
	// 关闭客户端
	if err = client.Stop(); err != nil {
		fmt.Printf("client stop error: %v \n", err)
	}

}
