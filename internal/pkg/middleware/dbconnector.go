package middleware

import (
	"errors"
	"os"
)

func GetConnectionString() (string, error) {
	key, flag := os.LookupEnv("DATABASE_URL")
	if !flag {
		return "", errors.New("connection string not found")
	}
	return key, nil
}
