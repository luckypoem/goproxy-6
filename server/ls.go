package server

import (
	"github.com/yuya008/goproxy/jiami"
	"log"
	"net"
)

type LocalService struct {
	// 本地地址
	Host string
	// 远程地址
	RemoteHost string
	// aes加密密钥路径
	AESKeyPath string
	AESkey     []byte
}

// 初始化本地服务
func InitLocalServer(ls *LocalService) error {
	// 初始化AES加密机制
	ls.AESkey = jiami.LoadAesKeyFile(ls.AESKeyPath)
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
	stream := jiami.NewAES(ls.AESkey, remoteConn)
	go LocalReader(ls, b, stream)
	go LocalWriter(ls, b, stream)
}
