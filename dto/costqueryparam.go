package dto

import (
	"errors"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CostSummaryReqDTO структура query-параметров при подсчете общей стоимости подписок за период, поля дат обязательный, остальные опциональны
type CostSummaryReqDTO struct {
	UserID      *uuid.UUID `json:"user_id,omitempty"`
	ServiceName *string    `json:"service_name,omitempty"`
	StartPeriod string     `json:"start_period"`
	EndPeriod   string     `json:"end_period"`
}

// CostSummaryRespDTO структура для вывода информации об общей стоимости, представляет из себя саму общую стоимость и фильтры
type CostSummaryRespDTO struct {
	TotalCost  int               `json:"total_cost"`
	QueryParam CostSummaryReqDTO `json:"query_param"`
}

// NewQueryParam создает экземпляр структуры CostSummaryReqDTO, проверяет входящие данные и заполняет ее
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

// NewCostSummaryRespDTO создает экземпляр структуры CostSummaryRespDTO и заполняет ее
func NewCostSummaryRespDTO(totalCost int, queryParam CostSummaryReqDTO) CostSummaryRespDTO {
	return CostSummaryRespDTO{
		TotalCost:  totalCost,
		QueryParam: queryParam,
	}
}
