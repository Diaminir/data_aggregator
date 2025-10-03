package dto

import (
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
