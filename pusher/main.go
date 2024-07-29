package main

import (
	"github.com/zserge/lorca"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 打开一个窗口页面
	// 页面加载资源为百度页面
	// 窗口尺寸为600 * 800
	windows, err := lorca.New("https://www.google.com/", "", 800, 600, "--remote-allow-origins=*")
	if err != nil {
		// 错误
		log.Fatalln(err)
	}
	defer windows.Close()

	// 程序优雅退出
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-windows.Done():
		// 窗口主动关闭
	case <-stop:
		// 程序关闭
	}
}
