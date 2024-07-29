// @Title client.go
// @Description 基于Go net标准库 实现TCP客户端
// @Author Zero - 2023/9/5 14:49:39

package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	// 要连接的服务端addr
	tcpServerAddr, _ := net.ResolveTCPAddr("tcp", ":8888")
	//tcpServerAddr, err := net.ResolveTCPAddr("tcp", "47.94.213.132:8888")
	conn, err := net.DialTCP("tcp", nil, tcpServerAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	// 创建一个读取的终端Reader
	reader := bufio.NewReader(os.Stdout)

	// 接收服务端的数据
	go func(c net.Conn) {
		for {
			buf := make([]byte, 1024)
			_, err := c.Read(buf)
			switch err {
			case io.EOF:
				fmt.Println("服务端关闭了连接")
				os.Exit(1)
			}
			//fmt.Printf("服务端响应信息(%d字节): %s", size, string(buf))
		}
	}(conn)

	for {
		fmt.Println("请输入:")
		// 读取终端输入的每一行内容
		line, _ := reader.ReadString('\n')
		line = strings.ReplaceAll(line, "\n", "")

		// 如果输入的是exit 就退出
		if line == "exit" {
			fmt.Println("TCP Client Exit.")
			break
		}
		// 将数据发送到服务端
		size, _ := conn.Write([]byte(line))
		fmt.Printf("发送数据成功,共%d字节。\n", size)
		//buf := make([]byte, 5000)
		//for i := range buf {
		//	buf[i] = 255
		//}
		//conn.Write(buf)

	}

}
