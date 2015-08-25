package server

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

const (
	// 协议状态头的长度
	STATUS_HEAD_LEN = 1
	// 数据体的长度
	DATA_LEN = 1024
	// 协议缓冲区长度
	BUFFER_LEN = STATUS_HEAD_LEN + DATA_LEN
)

// sock5协议解析错误
var SOCK5_PROTO_ERR = errors.New("sock5协议解析错误")

func Pack(b []byte, n int) {
	if len(b) < n {
		return
	}
	for i := 0; i < n; i++ {
		b[i]++
	}
}

func Unpack(b []byte, n int) {
	if len(b) < n {
		return
	}
	for i := 0; i < n; i++ {
		b[i]--
	}
}

func LocalReader(b *Brower) {
	log.Println("LocalReader")
	defer func() {
		log.Println("LocalReader 关闭连接")
		b.BrowerConn.Close()
		if !b.Ls.Pool.IsMember(b.RemoteConn) {
			b.Ls.Pool.Put(b.RemoteConn)
			log.Println("LocalReader 连接放回")
		}
	}()

	var buffer [BUFFER_LEN]byte
	data := buffer[STATUS_HEAD_LEN:]

	for {
		n0, err0 := b.BrowerConn.Read(data)
		if err0 != nil {
			if err0 == io.EOF {
				buffer[0] = 1 // 完成标志
				_, err1 := b.RemoteConn.Write(buffer[:])
				if err1 != nil {
					log.Println(err1)
				}
				return
			}
			log.Println(err0)
			return
		}
		Pack(data, n0)
		buffer[0] = 0 // 继续标志
		_, err2 := b.RemoteConn.Write(buffer[:n0+STATUS_HEAD_LEN])
		if err2 != nil {
			log.Println(err2)
			return
		}
	}
}

func LocalWriter(b *Brower) {
	log.Println("LocalReader")
	defer func() {
		log.Println("LocalReader 关闭连接")
		b.BrowerConn.Close()
		if !b.Ls.Pool.IsMember(b.RemoteConn) {
			b.Ls.Pool.Put(b.RemoteConn)
			log.Println("LocalReader 连接放回")
		}
	}()

	var buffer [BUFFER_LEN]byte
	data := buffer[STATUS_HEAD_LEN:]

	for {
		n0, err0 := b.RemoteConn.Read(buffer[:])
		if err0 != nil {
			log.Println(err0)
			return
		}
		// 关闭标志
		if buffer[0] == 1 {
			return
		}
		Unpack(data, n0)
		_, err1 := b.BrowerConn.Write(data[:n0])
		if err1 != nil {
			log.Println(err1)
			return
		}
	}
}

func sock5Read(conn net.Conn, buflen int, flag bool) ([]byte, int, error) {
	if flag {
		buflen += STATUS_HEAD_LEN
	}
	buf := make([]byte, buflen)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, n, err
	}
	if flag {
		if buf[0] == 1 {
			return nil, n, io.EOF
		}
		Unpack(buf[STATUS_HEAD_LEN:], n)
		return buf[STATUS_HEAD_LEN:], n - STATUS_HEAD_LEN, nil
	} else {
		Unpack(buf, n)
		return buf, n, nil
	}
}

func sock5Write(conn net.Conn, buf []byte, n int) (int, error) {
	Pack(buf, n)
	return conn.Write(buf[:n])
}

func Sock5(conn net.Conn) (string, error) {
	// step1
	buf0, _, err0 := sock5Read(conn, 3, true)
	if err0 != nil {
		return "", err0
	}
	if buf0[0] != 0x05 || buf0[1] != 0x01 ||
		buf0[2] != 0x00 {
		log.Println(buf0)
		return "", errors.New("协议错误1")
	}
	// step2
	_, err1 := sock5Write(conn, []byte{0x05, 0x00}, 2)
	if err1 != nil {
		return "", err1
	}
	// step3
	buf2, _, err2 := sock5Read(conn, 4, true)
	if err2 != nil {
		return "", err2
	}
	if buf2[0] != 0x05 || buf2[1] != 0x01 || buf2[2] != 0x00 {
		return "", errors.New("协议错误2")
	}
	var host_port string
	if buf2[3] == 0x03 {
		hostlenbuf, _, err4 := sock5Read(conn, 1, false)
		if err4 != nil {
			return "", err4
		}
		hostlen, n5 := binary.Uvarint(hostlenbuf)
		if n5 <= 0 {
			return "", errors.New("协议错误3")
		}
		hostbuf, _, err6 := sock5Read(conn, int(hostlen), false)
		if err6 != nil {
			return "", err6
		}

		host := string(hostbuf)

		portbuf, _, err7 := sock5Read(conn, 2, false)
		if err7 != nil {
			return "", err7
		}
		port, n6 := binary.Uvarint(portbuf)
		if n6 <= 0 {
			return "", errors.New("协议错误4")
		}
		host_port = fmt.Sprintf("%s:%d", host, port)
	} else if buf2[3] == 0x01 {
		hostbuffer, n, err := sock5Read(conn, 4, false)
		if err != nil || n != 4 {
			return "", err
		}

		ip0, n0 := binary.Uvarint(hostbuffer[0:1])
		ip1, n1 := binary.Uvarint(hostbuffer[1:2])
		ip2, n2 := binary.Uvarint(hostbuffer[2:3])
		ip3, n3 := binary.Uvarint(hostbuffer[3:4])

		if n0 <= 0 || n1 <= 0 || n2 <= 0 || n3 <= 0 {
			return "", errors.New("协议错误5")
		}

		host := fmt.Sprintf("%d.%d.%d.%d", ip0, ip1, ip2, ip3)

		portbuf, _, err := sock5Read(conn, 2, false)
		if err != nil {
			return "", err
		}
		port, n := binary.Uvarint(portbuf)
		if n <= 0 {
			return "", errors.New("协议错误6")
		}
		host_port = fmt.Sprintf("%s:%d", host, port)
	} else {
		return "", errors.New("协议错误1")
	}
	// step4
	_, err := sock5Write(conn, []byte{0x05, 0x00, 0x00, 0x01}, 4)
	if err != nil {
		return "", err
	}
	b2 := [6]byte{}
	_, err = sock5Write(conn, b2[:], 6)
	if err != nil {
		return "", err
	}
	return host_port, nil
}

func RemoteRead(conn net.Conn, targetconn net.Conn) {
	defer func() {
		targetconn.Close()
		remoteserveracceptProc(conn)
	}()
	buffer := make([]byte, BUFFER_LEN)
	data := buffer[STATUS_HEAD_LEN:]
	for {
		n, err := conn.Read(buffer)
		if err == nil {
			log.Println(err)
			return
		}
		if buffer[0] == 1 {
			return
		}
		Unpack(data, n-STATUS_HEAD_LEN)
		_, err = targetconn.Write(data)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func RemoteWriter(conn net.Conn, targetconn net.Conn) {
	defer func() {
		targetconn.Close()
	}()
	buffer := make([]byte, BUFFER_LEN)
	data := buffer[STATUS_HEAD_LEN:]
	for {
		n, err := targetconn.Read(data)
		if err != nil {
			log.Println(err)
			return
		}
		Pack(data, n)
		buffer[0] = 0
		_, err = conn.Write(buffer)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func SayBye(conn net.Conn) {
	var buffer [BUFFER_LEN]byte
	buffer[0] = 1
	conn.Write(buffer[:])
}
