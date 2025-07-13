package main

import (
	"bytes"
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"log"
)

// @Title example.go
// @Author Zero - 2024/8/11 09:57:54
// @Description

type echoServerHandle struct {
	gnet.BuiltinEventEngine
	eng gnet.Engine
}

// OnBoot 服务端初始化完毕事件
func (es *echoServerHandle) OnBoot(eng gnet.Engine) gnet.Action {
	return gnet.None
}

// OnTraffic 数据可读事件
func (es *echoServerHandle) OnTraffic(c gnet.Conn) (action gnet.Action) {
	buf, err := c.Next(-1)
	if err != nil {
		fmt.Printf("conn read error: %v \n", err)
		_ = c.Close()
		return
	}
	_, _ = c.Write(bytes.ToUpper(buf))
	return
}

func main() {
	server := &echoServerHandle{}
	if err := gnet.Run(server, "tcp://:9090", gnet.WithMulticore(true)); err != nil {
		log.Fatal(err)
	}
}
