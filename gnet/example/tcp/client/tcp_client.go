package main

import (
	"github.com/panjf2000/gnet"
	"time"
)

// @Title tcp_client.go
// @Description
// @Author Zero - 2024/8/11 09:57:54

type clientHandle struct {
}

func (ec *clientHandle) OnInitComplete(server gnet.Server) (action gnet.Action) {
	//TODO implement me
	panic("implement me")
}

func (ec *clientHandle) OnShutdown(server gnet.Server) {
	//TODO implement me
	panic("implement me")
}

func (ec *clientHandle) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	//TODO implement me
	panic("implement me")
}

func (ec *clientHandle) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	//TODO implement me
	panic("implement me")
}

func (ec *clientHandle) PreWrite(c gnet.Conn) {
	//TODO implement me
	panic("implement me")
}

func (ec *clientHandle) AfterWrite(c gnet.Conn, b []byte) {
	//TODO implement me
	panic("implement me")
}

func (ec *clientHandle) Tick() (delay time.Duration, action gnet.Action) {
	//TODO implement me
	panic("implement me")
}

func (ec *clientHandle) React(in []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	return
}

func main() {
	client, err := gnet.NewClient(&clientHandle{})
	if err != nil {
		panic(err)
	}
	_, err = client.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		panic(err)
	}
}
