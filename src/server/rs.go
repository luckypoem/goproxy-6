package server

import (
	"log"
	"net"
)

type RemoteService struct {
	Host string
}

func InitRemoteServer(rs *RemoteService) {
	listener, err := net.Listen("tcp", rs.Host)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			conn.Close()
			return
		}
		go remoteserveracceptProc(conn)
	}
}

func remoteserveracceptProc(conn net.Conn) {
	// sock5沟通
	host, err := Sock5(conn)
	if err != nil {
		conn.Close()
		return
	}
	log.Println("连接到目标服务器", host)
	// 连接目标服务器
	targetconn, err := net.Dial("tcp", host)
	if err != nil {
		conn.Close()
		return
	}
	// 数据传输
	go RemoteRead(conn, targetconn)
	go RemoteWriter(conn, targetconn)
}
