// @Title server.go
// @Description 基于Go net标准库 实现TCP服务端
// @Author Zero - 2023/9/5 14:49:36

package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	// 创建tcp协议的地址对象
	tcpAddr, _ := net.ResolveTCPAddr("tcp", ":8888")
	//创建tcp服务监听器
	listener, _ := net.ListenTCP("tcp", tcpAddr)
	fmt.Println("TCP Server Start Success...")
	for {
		// 等待tcp客户端的连接
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			break
		}
		// 启用协程 处理连接发来的数据
		go process(conn)
	}
}

func process(conn *net.TCPConn) {
	for {
		// 获取客户端的信息
		fmt.Printf("%s 连接到服务器...\n", conn.RemoteAddr().String())

		// 读取客户端传来的数据
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		switch err {
		case io.EOF:
			fmt.Println("客户端断开连接!")
			conn.Close()
			return
		default:
			break
		}
		fmt.Printf("接收到了客户端一个数据包，共 %d 个字节\n", n)
	}
}
