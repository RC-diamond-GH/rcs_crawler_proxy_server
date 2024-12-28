package main

import (
	"fmt"
	"os"

	"rcs_crawler_proxy_server/proxy"
	"rcs_crawler_proxy_server/util"
)

func main() {
	// 读取配置文件
	util.ReadConfig("./config.json")
	// 开启日志
	file, err := os.OpenFile("proxy.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Error while initial log file: %s", err.Error())
	}
	defer file.Close()

	util.InitLogger(file)

	// 初始化 Proxy
	proxy.ParseHTTPProxyList()
	proxy.PingProxyServer()

	// 初始化缓存数据库
	util.InitCacheDatabase()
	defer util.Cache.Close()

	// 启动 HTTP 代理服务器
	go proxy.HttpListener()
	console()
	fmt.Println("Bye~")
}

func console() {
	var cmd string = ""
	for cmd != "exit" {
		print("> ")
		fmt.Scan(&cmd)
	}
}
