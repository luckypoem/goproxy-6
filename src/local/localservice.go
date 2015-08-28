package main

import (
	"flag"
	"fmt"
	"log"
	"server"
)

// 标准输出输出启动参数
func printArgs(ls *server.LocalService) {
	fmt.Printf("本地服务监听主机及端口配置:%s\n", ls.Host)
	fmt.Printf("远程服务监听主机及端口配置:%s\n", ls.RemoteHost)
}

// 分析参数
func parseArgs(ls *server.LocalService) {
	flag.StringVar(&ls.Host, "h", "127.0.0.1:1324", "本地服务监听主机及端口配置")
	flag.StringVar(&ls.RemoteHost, "r", "127.0.0.1:44444", "远程服务监听主机及端口配置")
	flag.Parse()
	printArgs(ls)
}

func main() {
	ls := &server.LocalService{}
	// 分析启动参数
	parseArgs(ls)
	// 初始化长连接连接池
	log.Println("启动初始化本地服务器")
	server.InitLocalServer(ls)
}
