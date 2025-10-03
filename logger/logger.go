package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// NewLog логирование ошибок и событий в файл
func NewLog() (*logrus.Logger, *os.File) {
	logger := logrus.New()
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Fatal("Не удалось открыть файл для логирования:", err)
	}
	logger.SetLevel(logrus.InfoLevel)
	logger.SetOutput(io.MultiWriter(os.Stdout, file))
	return logger, file
}
