package main

import (
	"flag"
	"fmt"
	"log"
	"server"
)

// 标准输出输出启动参数
func printArgs(rs *server.RemoteService) {
	fmt.Printf("远程服务监听主机及端口配置:%s\n", rs.Host)
	fmt.Printf("AES加密Key文件路径:%s\n", rs.AESKeyPath)
}

// 分析参数
func parseArgs(rs *server.RemoteService) {
	flag.StringVar(&rs.Host, "h", "127.0.0.1:44444", "远程服务监听主机及端口配置")
	flag.StringVar(&rs.AESKeyPath, "k", "", "AES加密Key文件")
	flag.Parse()
	printArgs(rs)
}

func main() {
	rs := &server.RemoteService{}
	parseArgs(rs)
	log.Println("初始化并启动远程服务")
	server.InitRemoteServer(rs)
}
