package server

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"strconv"
)

const (
	// 数据体的长度
	DATA_LEN = 1024
)

var (
	// sock5协议解析错误
	SOCK5_PROTO_ERR = errors.New("sock5协议解析错误")
)

func Pack(b []byte) {
	if len(b) <= 0 {
		return
	}
	for i := 0; i < len(b); i++ {
		b[i]++
	}
}

func Unpack(b []byte) {
	if len(b) <= 0 {
		return
	}
	for i := 0; i < len(b); i++ {
		b[i]--
	}
}

func LocalReader(ls *LocalService, b net.Conn, r net.Conn) {
	defer func() {
		b.Close()
		r.Close()
	}()

	data := make([]byte, DATA_LEN)
	for {
		n, err := b.Read(data)
		if err != nil {
			return
		}
		Pack(data[:n])
		r.Write(data[:n])
	}
}

func LocalWriter(ls *LocalService, b net.Conn, r net.Conn) {
	defer func() {
		b.Close()
		r.Close()
	}()
	data := make([]byte, DATA_LEN)
	for {
		n, err := r.Read(data)
		if err != nil {
			return
		}
		Unpack(data[:n])
		b.Write(data[:n])
	}
}

func generalRead(c net.Conn) ([]byte, error) {
	buffer := make([]byte, DATA_LEN)
	n, err := c.Read(buffer)
	if err != nil {
		return nil, err
	}
	return buffer[:n], nil
}

func generalWrite(c net.Conn, buf []byte) error {
	_, err := c.Write(buf)
	return err
}

func Sock5(conn net.Conn) (string, error) {
	// step1
	buf0, err0 := generalRead(conn)
	if err0 != nil {
		return "", SOCK5_PROTO_ERR
	}
	Unpack(buf0)
	if buf0[0] != 0x05 || buf0[1] != 0x01 ||
		buf0[2] != 0x00 {
		return "", SOCK5_PROTO_ERR
	}
	// step2
	buf2 := []byte{0x05, 0x00}
	Pack(buf2)
	err1 := generalWrite(conn, buf2)
	if err1 != nil {
		return "", SOCK5_PROTO_ERR
	}
	// step3
	buf3, err := generalRead(conn)
	if err != nil {
		return "", SOCK5_PROTO_ERR
	}
	Unpack(buf3)
	if buf3[0] != 0x05 || buf3[1] != 0x01 || buf3[2] != 0x00 {
		return "", SOCK5_PROTO_ERR
	}

	var host_port string
	if buf3[3] == 0x03 { // 主机名
		hostlen, n5 := binary.Uvarint(buf3[4:5])
		if n5 <= 0 {
			return "", SOCK5_PROTO_ERR
		}
		host := string(buf3[5 : 5+hostlen])
		var port int16
		err := binary.Read(
			bytes.NewBuffer(buf3[5+hostlen:5+hostlen+2]),
			binary.BigEndian,
			&port)
		if err != nil {
			return "", SOCK5_PROTO_ERR
		}
		host_port = net.JoinHostPort(host, strconv.Itoa(int(port)))
	} else if buf3[3] == 0x01 { // IP
		ip := net.IPv4(buf3[4], buf3[5], buf3[6], buf3[7])
		port, n := binary.Uvarint(buf3[8:10])
		if n <= 0 {
			return "", SOCK5_PROTO_ERR
		}
		host_port = net.JoinHostPort(ip.String(), strconv.Itoa(int(port)))
	} else {
		return "", SOCK5_PROTO_ERR
	}
	// step4
	rep := []byte{
		0x05, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	Pack(rep)

	err = generalWrite(conn, rep)
	if err != nil {
		return "", SOCK5_PROTO_ERR
	}
	return host_port, nil
}

func RemoteRead(conn net.Conn, targetconn net.Conn) {
	defer func() {
		conn.Close()
		targetconn.Close()
	}()

	data := make([]byte, DATA_LEN)
	for {
		n, err := conn.Read(data)
		if err != nil {
			return
		}
		Unpack(data[:n])
		targetconn.Write(data[:n])
	}
}

func RemoteWriter(conn net.Conn, targetconn net.Conn) {
	defer func() {
		conn.Close()
		targetconn.Close()
	}()
	data := make([]byte, DATA_LEN)
	for {
		n, err := targetconn.Read(data)
		if err != nil {
			if err == io.EOF {
				return
			}
			return
		}
		Pack(data[:n])
		conn.Write(data[:n])
	}
}
