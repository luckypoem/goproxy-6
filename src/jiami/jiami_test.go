package jiami

import (
	_ "os"
	"testing"
)

var initplaintext []byte = []byte{84, 32, 47, 32, 72, 84, 84, 80, 47, 49, 46, 49, 13, 10, 13, 10}

func TestAES0(t *testing.T) {
	key := LoadAesKeyFile("../../aes.key")
	aes := NewAES(key, nil)
	ciphertext, err := aes.Encrypt(initplaintext)
	if err != nil {
		t.Error(err)
	}
	t.Log(ciphertext, err)
	t.Log(len(ciphertext))
	plaintext, err := aes.Decrypt(ciphertext)
	t.Log(string(plaintext))
	t.Log(len(string(plaintext)))
}

func TestAES1(t *testing.T) {
	// file, err := os.OpenFile("./test.file", os.O_RDWR|os.O_CREATE, 0777)
	// if err != nil {
	// 	t.Error(err)
	// }
	// key := LoadAesKeyFile("./key")
	// aes := NewAES(key, file)
	// _, err = aes.Write([]byte(initplaintext))
	// if err != nil {
	// 	t.Error(err)
	// }
	// aes.Close()
}

func TestAES2(t *testing.T) {
	// file, err := os.Open("./test.file")
	// if err != nil {
	// 	t.Error(err)
	// }
	// key := LoadAesKeyFile("./key")
	// aes := NewAES(key, file)
	// buf, err := aes.Read()
	// if err != nil {
	// 	t.Error(err)
	// }
	// t.Log(string(buf))
	// aes.Close()
}
