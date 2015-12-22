package jiami

type CryptoStream interface {
	Read() ([]byte, error)
	Write(d []byte) (int, error)
	Close() error
	Encrypt(src []byte) ([]byte, error)
	Decrypt(decryptText []byte) ([]byte, error)
}
