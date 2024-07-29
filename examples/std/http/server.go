// @Title server.go
// @Description
// @Author Zero - 2023/9/5 19:38:45

package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {

	// 注册一个http处理函数
	http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("hello"))
	})
	// 注册一个http处理函数，返回json格式数据
	http.HandleFunc("/json", func(writer http.ResponseWriter, request *http.Request) {
		// 设置响应码
		writer.WriteHeader(200)
		// 设置响应数据格式
		writer.Header().Set("Content-Type", "application/json")

		// 响应数据
		resp := map[string]any{"name": "张三", "age": 18}
		bytes, _ := json.Marshal(resp)
		_, _ = writer.Write(bytes)
	})

	http.HandleFunc("/timeTask", func(writer http.ResponseWriter, request *http.Request) {
		log.Println("定时任务被调用..")
	})

	// 启动http服务
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
