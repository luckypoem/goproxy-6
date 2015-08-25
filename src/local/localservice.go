package main

import (
	"connpool"
	"flag"
	"fmt"
	"log"
	"server"
)

// 标准输出输出启动参数
func printArgs(ls *server.LocalService) {
	fmt.Printf("本地服务监听主机及端口配置:%s\n", ls.Host)
	fmt.Printf("远程服务监听主机及端口配置:%s\n", ls.RemoteHost)
	fmt.Printf("连接池大小:%d\n", ls.ConnecterN)
}

// 分析参数
func parseArgs(ls *server.LocalService) {
	flag.StringVar(&ls.Host, "h", "127.0.0.1:1324", "本地服务监听主机及端口配置")
	flag.StringVar(&ls.RemoteHost, "r", "127.0.0.1:44444", "远程服务监听主机及端口配置")
	flag.IntVar(&ls.ConnecterN, "c", 100, "远程服务连接池的连接数目配置")
	flag.Parse()
	printArgs(ls)
}

func main() {
	ls := &server.LocalService{}
	// 分析启动参数
	parseArgs(ls)
	// 初始化长连接连接池

	log.Println("连接远程服务器，初始化连接池")
	ls.Pool = connpool.New(ls.RemoteHost, uint32(ls.ConnecterN))
	log.Println("初始化连接池完成")

	log.Println("启动初始化本地服务器")
	server.InitLocalServer(ls)
}
