package logger

import (
	"io"
	"junior_effectivemobile/config"
	"os"

	"github.com/sirupsen/logrus"
)

// NewLog логирование ошибок и событий в файл
func NewLog(cfg *config.Config) (*logrus.Logger, *os.File) {
	logger := logrus.New()
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Fatal("Не удалось открыть файл для логирования:", err)
	}
	switch cfg.LogLevel {
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	}
	logger.SetOutput(io.MultiWriter(os.Stdout, file))
	return logger, file
}
