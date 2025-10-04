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
	LogLevel   string
}

// NewConfig функция читает файл конфигурации и на основании ключей заполняет поля структуры Config
func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, errors.New("файл с конфигурацией не найден")
	}
	return &Config{
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT"),
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbName:     os.Getenv("DB_NAME"),
		LogLevel:   os.Getenv("LOG_LEVEL"),
	}, nil
}
