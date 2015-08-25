package server

import (
	"connpool"
	"log"
	"net"
)

type LocalService struct {
	// 本地地址
	Host string
	// 远程地址
	RemoteHost string
	// 连接池连接数目
	ConnecterN int
	// 连接池
	Pool connpool.ConnPool
}

type Brower struct {
	// 浏览器连接
	BrowerConn net.Conn
	// 远程服务器连接
	RemoteConn net.Conn
	// LocalService 结构
	Ls *LocalService
}

// 初始化本地服务
func InitLocalServer(ls *LocalService) error {
	listener, err := net.Listen("tcp", ls.Host)
	if err != nil {
		log.Println(err)
	}
	for {
		conn, err := listener.Accept()
		log.Println("一个连接进入")
		if err != nil {
			log.Fatal(err)
		}
		brower := &Brower{
			BrowerConn: conn,
			Ls:         ls,
		}
		go localserveracceptProc(brower)
	}
}

// accpet处理函数
func localserveracceptProc(brower *Brower) {
	remoteConn, err := brower.Ls.Pool.Get()
	log.Println("取得一个远程连接")
	if err != nil {
		log.Fatal(err)
	}
	brower.RemoteConn = remoteConn
	go LocalReader(brower)
	go LocalWriter(brower)
}
