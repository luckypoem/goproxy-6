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
	Keysize    = 32
	Vectorsize = aes.BlockSize
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
	_, err = self.iohandler.Read(buffer)
	if err != nil {
		return nil, err
	}
	// log.Println("Read", pkglen, buffer)
	res, err := self.Decrypt(buffer)
	log.Println("Read d", res)
	return res, err
}

func (self *aesSupport) Write(d []byte) (int, error) {
	log.Println("Write d", d)
	ciphertextbuffer, err := self.Encrypt(d)
	if err != nil {
		return -1, err
	}
	var pkglen int16 = int16(len(ciphertextbuffer))
	err = binary.Write(self.iohandler, binary.BigEndian, &pkglen)
	if err != nil {
		return -1, err
	}
	_, err = self.iohandler.Write(ciphertextbuffer)
	if err != nil {
		return -1, err
	}
	// log.Println("Write", pkglen, ciphertextbuffer)
	return len(d), nil
}

func (self *aesSupport) Close() error {
	return self.iohandler.Close()
}

func (self *aesSupport) Encrypt(src []byte) ([]byte, error) {
	src = Padding(src, aes.BlockSize)
	if len(src)%aes.BlockSize != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	encryptText := make([]byte, aes.BlockSize+len(src))
	iv := encryptText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cblock, err := aes.NewCipher(self.key)
	if err != nil {
		log.Panicln("aes.NewCipher: " + err.Error())
	}
	mode := cipher.NewCBCEncrypter(cblock, iv)
	mode.CryptBlocks(encryptText[aes.BlockSize:], src)
	return encryptText, nil
}

func (self *aesSupport) Decrypt(decryptText []byte) ([]byte, error) {
	if len(decryptText) < aes.BlockSize {
		return nil, errors.New("crypto/cipher: ciphertext too short")
	}
	iv := decryptText[:aes.BlockSize]
	decryptText = decryptText[aes.BlockSize:]
	if len(decryptText)%aes.BlockSize != 0 {
		return nil, errors.New("crypto/cipher: ciphertext is not a multiple of the block size")
	}
	cblock, err := aes.NewCipher(self.key)
	if err != nil {
		log.Panicln("aes.NewCipher: " + err.Error())
	}
	mode := cipher.NewCBCDecrypter(cblock, iv)
	mode.CryptBlocks(decryptText, decryptText)
	return UnPadding(decryptText), nil
}
