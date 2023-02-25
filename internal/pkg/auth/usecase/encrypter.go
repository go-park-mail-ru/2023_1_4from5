package usecase

import (
	"crypto/sha256"
	"fmt"
	"os"
)

type Encrypter struct {
	salt string
}

func NewEncrypter() *Encrypter {
	salt := os.Getenv("SECRET")
	return &Encrypter{salt: salt}
}

func (ec *Encrypter) EncryptPswd(pswd string) string {
	Encryptedpswd := sha256.New()
	_, err := Encryptedpswd.Write([]byte(pswd))
	if err != nil {
		return ""
	}
	_, err = Encryptedpswd.Write([]byte(ec.salt))
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", Encryptedpswd.Sum(nil))
}
