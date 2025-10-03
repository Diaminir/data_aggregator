package db

import (
	"context"
	"fmt"
	"junior_effectivemobile/dto"

	"github.com/jackc/pgx/v5"
)

type Postgres struct {
	db *pgx.Conn
}

func NewConPostgres() (*Postgres, error) {
	connFig := "postgres://qwert:12345@localhost:8081/subscriptions"

	conn, err := pgx.Connect(context.Background(), connFig)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к Postgres: %s", err)
	}
	return &Postgres{
		db: conn,
	}, nil
}

func (pg *Postgres) DbClose() {
	pg.db.Close(context.Background())
}

func (pg *Postgres) PostNewSubRecord(rec dto.SubRecordDTO) (dto.SubRecordWithIdDTO, error) {

	return dto.SubRecordWithIdDTO{}, nil
}

func (pg *Postgres) GetSubRecord(id int) (dto.SubRecordWithIdDTO, error) {

	return dto.SubRecordWithIdDTO{}, nil
}

func (pg *Postgres) GetListSubRecords() ([]dto.SubRecordWithIdDTO, error) {

	return []dto.SubRecordWithIdDTO{}, nil
}

func (pg *Postgres) UpdateSubRecord(id int, updateData dto.UpdateSubRecordDTO) (dto.SubRecordWithIdDTO, error) {

	return dto.SubRecordWithIdDTO{}, nil
}

func (pg *Postgres) DeleteSubRecord(id int) error {

	return nil
}

func (pg *Postgres) CalculateCost(queryParam dto.CostSummaryReqDTO) (dto.CostSummaryRespDTO, error) {

	return dto.CostSummaryRespDTO{}, nil
}
