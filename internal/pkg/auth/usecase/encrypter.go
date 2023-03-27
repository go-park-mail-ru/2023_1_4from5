package usecase

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
)

type Encrypter struct {
	salt string
}

func NewEncryptor() (*Encrypter, error) {
	salt, flag := os.LookupEnv("ENCRYPTER_SECRET")
	if !flag {
		return &Encrypter{}, errors.New("NoSecretKey")
	}
	return &Encrypter{salt: salt}, nil
}

func (ec *Encrypter) EncryptPswd(ctx context.Context, pswd string) string {
	encryptedPswd := sha256.New()
	_, err := encryptedPswd.Write([]byte(pswd))
	if err != nil {
		return ""
	}
	_, err = encryptedPswd.Write([]byte(ec.salt))
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", encryptedPswd.Sum(nil))
}
