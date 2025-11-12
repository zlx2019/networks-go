package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"github.com/zlx2019/go-net-v2/simple_protocol/protocol"
	"io"
	"log"
	"math"
	"sync/atomic"
)

type simpleServer struct {
	gnet.BuiltinEventEngine
	eng           gnet.Engine
	network       string
	addr          string
	multicore     bool
	connected     int32
	disconnected  int32
	maxBatchReads int
}

func (s *simpleServer) ServeAddr() string {
	return fmt.Sprintf("%s://%s", s.network, s.addr)
}

// OnBoot 服务端引擎就绪事件
func (s *simpleServer) OnBoot(eng gnet.Engine) gnet.Action {
	fmt.Println("Service is ready, starting to accept connections")
	s.eng = eng
	return gnet.None
}

// OnShutdown 服务引擎关闭事件
func (s *simpleServer) OnShutdown(eng gnet.Engine) {
	fmt.Println("Service closed.")
}

// OnOpen 新的客户端连接事件
func (s *simpleServer) OnOpen(conn gnet.Conn) (out []byte, action gnet.Action) {
	fmt.Printf("Accept connection from %s\n", conn.RemoteAddr().String())
	atomic.AddInt32(&s.connected, 1)
	// 为连接设置编解码器
	conn.SetContext(new(protocol.SimpleCodec))
	out = []byte("Welcome!\n")
	return
}

// OnClose 客户端关闭连接
func (s *simpleServer) OnClose(conn gnet.Conn, err error) (action gnet.Action) {
	if errors.Is(err, io.EOF) {
		fmt.Printf("Connection closed from %s\n", conn.RemoteAddr().String())
	} else {
		fmt.Printf("Connection %s Closed due to error: %v \n", conn.RemoteAddr().String(), err)
	}
	_ = atomic.AddInt32(&s.disconnected, 1)
	_ = atomic.AddInt32(&s.connected, -1)
	//if connected == 0 {
	// fmt.Printf("all %d connections are closed, shut it down \n", disconnected)
	// action = gnet.Shutdown
	//}
	return
}

// OnTraffic 数据可读事件
func (s *simpleServer) OnTraffic(c gnet.Conn) (action gnet.Action) {
	codec, ok := c.Context().(*protocol.SimpleCodec)
	if !ok {
		fmt.Printf("Connect context codec not found \n")
		return gnet.Close
	}
	var packets [][]byte
	// 批次读取数据包
	for i := 0; i < s.maxBatchReads; i++ {
		data, err := codec.Decode(c)
		if err != nil {
			if errors.Is(err, protocol.ErrIncompletePacket) {
				break
			}
			fmt.Printf("invalid packet: %v \n", err)
			return gnet.Close
		}
		packet, _ := codec.Encode(bytes.ToUpper(data))
		packets = append(packets, packet)
	}

	// 写回响应数据
	readN := len(packets)
	if readN > 1 {
		_, _ = c.Writev(packets)
	} else if readN == 1 {
		_, _ = c.Write(packets[0])
	}
	if readN == s.maxBatchReads && c.InboundBuffered() > 0 {
		// 本次已经读取到最大上限数据包大小，但缓冲区依然还有数据
		if err := c.Wake(nil); err != nil {
			fmt.Printf("failed to wake up the connection: %v \n", err)
			return gnet.Close
		}
	}
	return
}

func main() {
	var (
		port          int
		multicore     bool
		maxBatchReads int
	)
	flag.IntVar(&port, "π", 9601, "server port")
	flag.BoolVar(&multicore, "m", false, "enable multicore")
	flag.IntVar(&maxBatchReads, "b", 100, "read batch size")
	flag.Parse()
	if maxBatchReads <= 0 {
		maxBatchReads = math.MaxInt32
	}
	server := &simpleServer{
		network:       "tcp",
		addr:          fmt.Sprintf(":%d", port),
		multicore:     multicore,
		maxBatchReads: maxBatchReads,
	}
	err := gnet.Run(server, server.ServeAddr(), gnet.WithMulticore(multicore), gnet.WithReusePort(true))
	if err != nil {
		log.Fatal(err)
	}
}
