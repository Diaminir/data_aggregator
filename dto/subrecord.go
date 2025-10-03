package dto

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// SubRecordDTO структура для информации о записи
type SubRecordDTO struct {
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date"`
}

// ValidateInputData функция для валидации входящих данных при создании записи
func (sr SubRecordDTO) ValidateInputData() error {
	if sr.ServiceName == "" && sr.Price == 0 && sr.UserID == uuid.Nil && sr.StartDate == "" && sr.EndDate == "" {
		return errors.New("введите данные")
	}
	if sr.ServiceName == "" {
		return errors.New("введите название сервиса")
	}
	if sr.Price <= 0 {
		return errors.New("введите стоимость подписки, стоимость не может быть равной или меньше нуля")
	}
	if sr.UserID == uuid.Nil {
		return errors.New("введите uuid пользователя")
	}
	if sr.StartDate == "" {
		return errors.New("введите дату начала подписки")
	}
	if sr.EndDate == "" {
		return errors.New("введите дату окончания подписки")
	}
	if _, err := time.Parse("2006-01-02", sr.StartDate); err != nil {
		return errors.New("неверный формат даты, используйте YYYY-MM-DD")
	}
	if _, err := time.Parse("2006-01-02", sr.EndDate); err != nil {
		return errors.New("неверный формат даты, используйте YYYY-MM-DD")
	}
	start, _ := time.Parse("2006-01-02", sr.StartDate)
	end, _ := time.Parse("2006-01-02", sr.EndDate)
	if end.Before(start) || end.Equal(start) {
		return errors.New("дата окончания должна быть позже даты начала")
	}
	return nil
}
