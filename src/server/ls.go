package server

import (
	"log"
	"net"
)

type LocalService struct {
	// 本地地址
	Host string
	// 远程地址
	RemoteHost string
}

// 初始化本地服务
func InitLocalServer(ls *LocalService) error {
	listener, err := net.Listen("tcp", ls.Host)
	if err != nil {
		log.Println(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go localserveracceptProc(conn, ls)
	}
}

// accpet处理函数
func localserveracceptProc(b net.Conn, ls *LocalService) {
	remoteConn, err := net.Dial("tcp", ls.RemoteHost)
	if err != nil {
		b.Close()
		return
	}
	go LocalReader(ls, b, remoteConn)
	go LocalWriter(ls, b, remoteConn)
}
