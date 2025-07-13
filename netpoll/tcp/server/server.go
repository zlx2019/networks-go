// @Title server.go
// @Description 使用netpoll 实现一个tcp服务端(NIO)
// @Author Zero - 2023/8/12 15:48:33

package main

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/cloudwego/netpoll"
	"log"
	"math/rand"
	"strings"
	"time"
)

func main() {
	// 1. 创建tcp监听器
	listener, err := netpoll.CreateListener("tcp", "127.0.0.1:8989")
	if err != nil {
		panic("create tcp listener failed: " + err.Error())
	}
	// 2. 创建事件循环调度器,这是一个真正的NIO服务,负责连接管理、事件处理等
	// 注册事件处理函数，以及一些hook函数，如连接初始化时、连接完成后;
	// 设置连接的读取超时时间
	loop, err := netpoll.NewEventLoop(
		OnRequest,
		netpoll.WithOnPrepare(OnPrepare),
		netpoll.WithOnConnect(OnConnect),
		netpoll.WithReadTimeout(time.Second*3),
	)
	if err != nil {
		panic("create netpoll listener failed")
	}
	log.Println("server start successful...")

	// 优雅关闭服务
	//go func() {
	//	time.Sleep(time.Second * 3)
	//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	//	defer cancel()
	//	loop.Shutdown(ctx)
	//}()

	// 3.运行服务
	if err = loop.Serve(listener); err != nil {
		panic("start nio tcp server failed")
	}
}

// OnPrepare 连接初始化时的执行hook函数
// 用于在连接初始化时注入自定义准备，这是可选的,但在某些情况下很重要。
// 返回的context上下文将成为 OnConnect和OnRequest的参数。
func OnPrepare(connection netpoll.Connection) context.Context {
	log.Printf("[%s] 连接初始化... \n", connection.RemoteAddr().String())

	// 可以在这里向下传递一些K-V参数
	ctx := context.WithValue(context.Background(), "connkey", rand.Int())
	return ctx
}

// OnConnect 创建连接之后的执行hook函数
func OnConnect(ctx context.Context, connection netpoll.Connection) context.Context {
	log.Printf("[%s] 连接创建完成... \n", connection.RemoteAddr().String())
	// 为连接注册关闭hook函数
	_ = connection.AddCloseCallback(OnClose)
	return ctx
}

// OnClose 连接关闭之后执行的hook函数
func OnClose(connection netpoll.Connection) error {
	if !connection.IsActive() {
		if err := connection.Close(); err != nil {
			return err
		}
		log.Printf("[%s] 连接关闭... \n", connection.RemoteAddr().String())
	}
	return nil
}

// OnRequest 调度器的事件处理
func OnRequest(ctx context.Context, connection netpoll.Connection) error {
	fmt.Println("zero")
	// 判断连接是否活跃
	if connection.IsActive() {
		// 获取连接读取器，通过读取器以NIO模型进行交互
		reader := connection.Reader()
		// 获取读取器此时可读的字节数量
		readLen := reader.Len()
		fmt.Println("可读取的数据的总长度: ", readLen)

		// TODO 粘包与拆包处理
		// TODO 读取前8个byte位置的数据，获取数据长度
		// 这里只是读取，并不会移除数据
		lengthBytes, _ := reader.Peek(8)
		length := binary.BigEndian.Uint64(lengthBytes)
		fmt.Println("要读取的数据长度: ", length)

		// TODO 如果当前读取器可读字节量不足数据长度，表示数据不足，无法解析，需要等待更多数据的到来
		if uint64(readLen-8) < length {
			return nil
		}

		// 跳过前8个byte(数据头部)
		reader.Skip(8)

		// 根据数据长度，读取具体的数据
		bytes, err := reader.Next(int(length))
		if err != nil {
			return err
		}
		// 释放资源
		reader.Release()

		// 截取末尾的\n
		message := strings.TrimSuffix(string(bytes), "\n")
		fmt.Println(message)
		return nil
	}
	return errors.New("连接已关闭..")
}
