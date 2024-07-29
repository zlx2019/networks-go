// @Title client.go
// @Description 基于Go net标准库 实现UDP客户端
// @Author Zero - 2023/9/5 14:49:39

package main

import (
	"fmt"
	"net"
)

var sendCounter int = 0

func main() {
	udpAddr, _ := net.ResolveUDPAddr("udp", ":8889")
	conn, _ := net.DialUDP("udp", nil, udpAddr)
	defer conn.Close()
	buf := make([]byte, 1024)
	for i := range buf {
		buf[i] = 255
	}
	// 循环向服务端写入数据
	for i := 0; i < 1000; i++ {
		n, err := conn.Write(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		sendCounter += n
		fmt.Printf("发送 %d 字节成功，共发送了 %d 字节 \n", n, sendCounter)
	}
}
