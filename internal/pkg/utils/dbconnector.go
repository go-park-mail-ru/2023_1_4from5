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

func GetConnectionStringAuth() (string, error) {
	key, flag := os.LookupEnv("DATABASE_URL_AUTH")
	if !flag {
		return "", errors.New("connection string not found")
	}
	return key, nil
}

func GetConnectionStringUser() (string, error) {
	key, flag := os.LookupEnv("DATABASE_URL_USER")
	if !flag {
		return "", errors.New("connection string not found")
	}
	return key, nil
}

func GetConnectionStringCreator() (string, error) {
	key, flag := os.LookupEnv("DATABASE_URL_CREATOR")
	if !flag {
		return "", errors.New("connection string not found")
	}
	return key, nil
}
