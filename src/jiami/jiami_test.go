package jiami

import (
	"os"
	"testing"
)

const initplaintext = "1111111111111111"

func TestAES0(t *testing.T) {
	key := LoadAesKeyFile("./key")
	aes := NewAES(key, nil)
	ciphertext, err := aes.Encrypt([]byte(initplaintext))
	if err != nil {
		t.Error(err)
	}
	t.Log(ciphertext, err)
	t.Log(len(ciphertext))
	plaintext, err := aes.Decrypt(ciphertext)
	t.Log(string(plaintext))
	t.Log(len(string(plaintext)))
	if string(plaintext) != initplaintext ||
		len(plaintext) != len(initplaintext) {
		t.Error("error")
	}
}

func TestAES1(t *testing.T) {
	file, err := os.OpenFile("./test.file", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		t.Error(err)
	}
	key := LoadAesKeyFile("./key")
	aes := NewAES(key, file)
	_, err = aes.Write([]byte(initplaintext))
	if err != nil {
		t.Error(err)
	}
	aes.Close()
}

func TestAES2(t *testing.T) {
	file, err := os.Open("./test.file")
	if err != nil {
		t.Error(err)
	}
	key := LoadAesKeyFile("./key")
	aes := NewAES(key, file)
	buf, err := aes.Read()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(buf))
	aes.Close()
}
