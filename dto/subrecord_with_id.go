package dto

import "github.com/google/uuid"

// SubRecordWithIdDTO структура для вывода информации об записях, содержит id записи и саму запись
type SubRecordWithIdDTO struct {
	ID        int          `json:"id"`
	SubRecord SubRecordDTO `json:"subRecord"`
}

// NewSubRecordWithIdDTO создает экземпляр структуры SubRecordWithIdDTO и заполняет ее
func NewSubRecordWithIdDTO(id int, serviceName string, price int, userID uuid.UUID, startDate string, endDate string) SubRecordWithIdDTO {
	return SubRecordWithIdDTO{
		ID: id,
		SubRecord: SubRecordDTO{
			ServiceName: serviceName,
			Price:       price,
			UserID:      userID,
			StartDate:   startDate,
			EndDate:     endDate,
		},
	}
}
