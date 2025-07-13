// @Title server.go
// @Description 基于Go net标准库 实现UDP服务端
// @Author Zero - 2023/9/5 14:49:36

package main

import (
	"fmt"
	"net"
)

// 接收的数据总和 byte
var receiveCounter int = 0

func main() {
	// 定义 UDP 服务地址
	tcpAddr, _ := net.ResolveUDPAddr("udp", ":8889")
	// 启动UDP 服务
	listener, _ := net.ListenUDP("udp", tcpAddr)
	defer listener.Close()
	fmt.Println("UDP Server Start Success...")
	buf := make([]byte, 1024)
	for {
		// 循环读取发送过来的所有客户端数据
		n, err := listener.Read(buf)
		if err != nil {
			fmt.Println(err)
			continue
		}
		receiveCounter += n
		fmt.Printf("接收到数据包, %d 个字节，已接收到的数据总和: %d \n", n, receiveCounter)
	}
}
