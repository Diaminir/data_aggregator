package dto

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MessageDTO struct {
	Message string `json:"message"`
	Time    string `json:"time"`
}

type UpdateSubRecordDTO struct {
	ServiceName *string    `json:"service_name,omitempty"`
	Price       *int       `json:"price,omitempty"`
	UserID      *uuid.UUID `json:"user_id,omitempty"`
	StartDate   *string    `json:"start_date,omitempty"`
	EndDate     *string    `json:"end_date,omitempty"`
}

type SubRecordDTO struct {
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date"`
}

type SubRecordWithIdDTO struct {
	ID        int          `json:"id"`
	SubRecord SubRecordDTO `json:"subRecord"`
}

type CostSummaryReqDTO struct {
	UserID      *uuid.UUID `json:"user_id,omitempty"`
	ServiceName *string    `json:"service_name,omitempty"`
	StartPeriod string     `json:"start_period"`
	EndPeriod   string     `json:"end_period"`
}

type CostSummaryRespDTO struct {
	TotalCost  int               `json:"total_cost"`
	QueryParam CostSummaryReqDTO `json:"query_param"`
}

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

func NewQueryParam(c *gin.Context) (CostSummaryReqDTO, error) {
	var queryParam CostSummaryReqDTO
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return CostSummaryReqDTO{}, errors.New("неверный формат userId")
		}
		queryParam.UserID = &userID
	}
	if serviceName := c.Query("service_name"); serviceName != "" {
		queryParam.ServiceName = &serviceName
	}
	if startPeriodStr := c.Query("start_period"); startPeriodStr != "" {
		if _, err := time.Parse("2006-01-02", startPeriodStr); err != nil {
			return CostSummaryReqDTO{}, errors.New("неверный формат даты, используйте YYYY-MM-DD")
		}
		queryParam.StartPeriod = startPeriodStr
	} else {
		return CostSummaryReqDTO{}, errors.New("введите дату начала периода поиска")
	}
	if endPeriodStr := c.Query("end_period"); endPeriodStr != "" {
		if _, err := time.Parse("2006-01-02", endPeriodStr); err != nil {
			return CostSummaryReqDTO{}, errors.New("неверный формат даты, используйте YYYY-MM-DD")
		}
		queryParam.EndPeriod = endPeriodStr
	} else {
		return CostSummaryReqDTO{}, errors.New("введите дату конца периода поиска")
	}
	return queryParam, nil
}

func NewCostSummaryRespDTO(totalCost int, queryParam CostSummaryReqDTO) CostSummaryRespDTO {
	return CostSummaryRespDTO{
		TotalCost:  totalCost,
		QueryParam: queryParam,
	}
}

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
