package dto

import (
	"fmt"
	"time"
)

// MessageDTO структура сообщения, может сообщать как об ошибке, так и об успешном завершении операции
type MessageDTO struct {
	Message string `json:"message"`
	Time    string `json:"time"`
}

// Создает экземпляр структуры MessageDTO и заполняет ее
func NewMessageDTO(message string, err error) *MessageDTO {
	var str string
	if err != nil {
		str = fmt.Sprintf("%s: %s", message, err)
	} else {
		str = message
	}
	return &MessageDTO{
		Message: str,
		Time:    time.Now().Format("2006-01-02 15:04:05"),
	}
}
