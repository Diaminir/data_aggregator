package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// Config cтруктура для входящей конфигурации с .env файла
type Config struct {
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
}

// NewConfig функция читает файл конфигурации и на основании ключей заполняет поля структуры Config
func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, errors.New("файл с конфигурацией не найден")
	}
	return &Config{
		DbHost:     getEnv("DB_HOST"),
		DbPort:     getEnv("DB_PORT"),
		DbUser:     getEnv("DB_USER"),
		DbPassword: getEnv("DB_PASSWORD"),
		DbName:     getEnv("DB_NAME"),
	}, nil
}

// getEnv получает значение из файла по ключу
func getEnv(key string) string {
	return os.Getenv(key)
}
