package jiami

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"os"
)

const (
	// AES-256
	Keysize = 32
)

type aesSupport struct {
	key       []byte
	block     cipher.Block
	iohandler io.ReadWriteCloser
}

func LoadAesKeyFile(filename string) []byte {
	keybuf := make([]byte, 32)
	file, err := os.Open(filename)
	if err != nil {
		log.Panicln(filename, err)
	}
	n, err := io.ReadFull(file, keybuf)
	if err != nil {
		log.Panicln("readFull", err)
	}
	if n != Keysize {
		log.Panicln(n, "key size error")
	}
	return keybuf
}

func NewAES(key []byte, ioh io.ReadWriteCloser) CryptoStream {
	block, e := aes.NewCipher(key)
	if e != nil {
		log.Panicln(e.Error())
	}
	return &aesSupport{
		key:       key,
		block:     block,
		iohandler: ioh,
	}
}

func (self *aesSupport) Read() ([]byte, error) {
	var pkglen int16
	err := binary.Read(self.iohandler, binary.BigEndian, &pkglen)
	if pkglen <= 0 {
		return nil, err
	}
	buffer := make([]byte, pkglen)
	var c int16
	for pkglen > c {
		n, err := self.iohandler.Read(buffer[c:])
		if err != nil {
			return nil, err
		}
		c += int16(n)
	}
	return self.Decrypt(buffer)
}

func (self *aesSupport) Write(d []byte) (int, error) {
	ciphertextbuffer, err := self.Encrypt(d)
	if err != nil {
		return -1, err
	}
	var pkglen int16 = int16(len(ciphertextbuffer))
	if pkglen <= 0 {
		return -1, errors.New("Encrypt error")
	}
	err = binary.Write(self.iohandler, binary.BigEndian, &pkglen)
	if err != nil {
		return -1, err
	}
	_, err = self.iohandler.Write(ciphertextbuffer)
	if err != nil {
		return -1, err
	}
	return len(d), nil
}

func (self *aesSupport) Close() error {
	return self.iohandler.Close()
}

func (self *aesSupport) Encrypt(data []byte) ([]byte, error) {
	// 进行数据填充
	data = Padding(data, aes.BlockSize)
	// 判断密文是否是加密块整数倍
	if len(data)%aes.BlockSize != 0 {
		return nil, errors.New("jiami: input not full blocks")
	}
	// 构造一个密文+初始向量的空间
	encryptText := make([]byte, aes.BlockSize+len(data))
	iv := encryptText[:aes.BlockSize]
	// 获取一个随机的初始向量
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	// 设置加密模式
	mode := cipher.NewCBCEncrypter(self.block, iv)
	// 加密
	mode.CryptBlocks(encryptText[aes.BlockSize:], data)
	return encryptText, nil
}

func (self *aesSupport) Decrypt(data []byte) ([]byte, error) {
	// 提取初始向量
	iv := data[:aes.BlockSize]
	// 提取密文
	ciphertext := data[aes.BlockSize:]
	// 判断密文是否小于加密块大小
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("jiami: ciphertext too short")
	}
	// 判断密文是否是加密块整数倍
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("jiami: input not full blocks")
	}
	// 解密模式
	mode := cipher.NewCBCDecrypter(self.block, iv)
	ret := make([]byte, len(ciphertext))
	// 解密
	mode.CryptBlocks(ret, ciphertext)
	// 取消填充
	return UnPadding(ret), nil
}
