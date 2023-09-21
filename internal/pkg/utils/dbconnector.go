package utils

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

func Init() {
	_ = godotenv.Load(".env")
}

func GetConnectionString() (string, error) {
	key, flag := os.LookupEnv("DATABASE_URL")
	if !flag {
		return "", errors.New("connection string not found")
	}
	return key, nil
}
