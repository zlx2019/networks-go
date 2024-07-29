// @Description 使用netpoll 实现一个tcp客户端(NIO)
// @Author Zero - 2023/8/12 17:04:00

package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/cloudwego/netpoll"
	"log"
	"os"
	"time"
)

func main() {
	network, addr := "tcp", "127.0.0.1:8989"
	// 1. 连接服务端
	conn, err := netpoll.DialConnection(network, addr, time.Second)
	if err != nil {
		panic("conn server failed: " + err.Error())
	}
	log.Println("conn server successful...")

	// 2.注册事件处理函数
	_ = conn.SetOnRequest(OnRequest)
	// 注册连接关闭hook函数
	_ = conn.AddCloseCallback(OnClose)

	// 3. 获取连接的写入器,写入器通过NIO模型向服务端写数据
	writer := conn.Writer()

	// 创建终端读取器，以终端输入作为发送的数据内容
	reader := bufio.NewReader(os.Stdin)
	for {
		// 获取终端输入数据，写入到连接中
		// TODO 为了解决TCP粘包与拆包，将数据的长度以uint64格式(8byte) 写在前8个byte位置，方便服务端读取;
		// TODO [数据长度(8byte)|数据体]
		line, _, _ := reader.ReadLine()
		// 数据长度
		dataLen := len(line)
		// 分配缓冲区，容量为: 8 + 数据长度
		buf, _ := writer.Malloc(8 + dataLen)
		// 标记数据长度，将数据长度写入到缓冲区索引0~8的位置
		binary.BigEndian.PutUint64(buf[:8], uint64(dataLen))
		// 将数据体写入到缓冲区索引8~最末尾的位置
		copy(buf[8:], line)
		// 将缓冲区内的数据，发送给连接
		_ = writer.Flush()
	}
}

// OnRequest 服务端事件处理函数
func OnRequest(ctx context.Context, connection netpoll.Connection) error {
	reader := connection.Reader()
	defer reader.Release()
	bytes, err := reader.Next(reader.Len())
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))
	return nil
}

// OnClose 服务端连接关闭hook处理函数
func OnClose(connection netpoll.Connection) error {
	_ = connection.Close()
	fmt.Println("服务器已关闭连接...")
	return nil
}
