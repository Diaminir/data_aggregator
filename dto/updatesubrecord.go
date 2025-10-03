package dto

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// UpdateSubRecordDTO структура для изменения данных записи, хотя бы одно поле должно быть заполнено
type UpdateSubRecordDTO struct {
	ServiceName *string    `json:"service_name,omitempty"`
	Price       *int       `json:"price,omitempty"`
	UserID      *uuid.UUID `json:"user_id,omitempty"`
	StartDate   *string    `json:"start_date,omitempty"`
	EndDate     *string    `json:"end_date,omitempty"`
}

// ValidateUpdateData функция для валидации входящих данных при редактировании данных
func (upd *UpdateSubRecordDTO) ValidateUpdateData() error {
	if upd.ServiceName == nil && upd.Price == nil && upd.UserID == nil && upd.StartDate == nil && upd.EndDate == nil {
		return errors.New("хотя бы одно из полей должно быть заполнено")
	}
	if upd.ServiceName != nil && *upd.ServiceName == "" {
		return errors.New("введите название сервиса")
	}
	if upd.Price != nil && *upd.Price <= 0 {
		return errors.New("введите стоимость подписки, стоимость не может быть равной или меньше нуля")
	}
	if upd.UserID != nil && *upd.UserID == uuid.Nil {
		return errors.New("введите uuid пользователя")
	}
	if upd.StartDate != nil {
		if _, err := time.Parse("2006-01-02", *upd.StartDate); err != nil {
			return errors.New("неверный формат даты, используйте YYYY-MM-DD")
		}
	}
	if upd.EndDate != nil {
		if _, err := time.Parse("2006-01-02", *upd.EndDate); err != nil {
			return errors.New("неверный формат даты, используйте YYYY-MM-DD")
		}
	}
	if upd.StartDate != nil && upd.EndDate != nil {
		start, _ := time.Parse("2006-01-02", *upd.StartDate)
		end, _ := time.Parse("2006-01-02", *upd.EndDate)
		if end.Before(start) {
			return errors.New("дата окончания должна быть позже даты начала")
		}
	}
	return nil
}
