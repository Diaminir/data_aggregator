package db

import (
	"context"
	"fmt"
	"junior_effectivemobile/dto"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

	sql := `CREATE TABLE IF NOT EXISTS subscriptions
(
	id serial PRIMARY KEY,
	service_name varchar(255) NOT NULL,
    price int NOT NULL CHECK (price > 0),
    user_id uuid NOT NULL,
    start_date date NOT NULL,
    end_date date NULL
)`

	if _, err := conn.Exec(context.Background(), sql); err != nil {
		return nil, err
	}

	return &Postgres{
		db: conn,
	}, nil
}

func (pg *Postgres) DbClose() {
	pg.db.Close(context.Background())
}

func (pg *Postgres) PostNewSubRecord(rec dto.SubRecordDTO) (dto.SubRecordWithIdDTO, error) {
	sql := `INSERT INTO subscriptions(service_name, price, user_id, start_date, end_date) 
			VALUES ($1,$2,$3,$4,$5)
			RETURNING *`
	rowRec, err := pg.db.Query(context.Background(), sql, rec.ServiceName, rec.Price, rec.UserID, rec.StartDate, rec.EndDate)
	if err != nil {
		return dto.SubRecordWithIdDTO{}, err
	}
	defer rowRec.Close()
	subRec, err := pg.rowsInSlice(rowRec)
	if err != nil {
		return dto.SubRecordWithIdDTO{}, err
	}
	return subRec[0], nil
}

func (pg *Postgres) GetSubRecord(id int) (dto.SubRecordWithIdDTO, error) {
	sql := `SELECT id, service_name, price, user_id, start_date, end_date 
			FROM subscriptions 
			WHERE id = $1`
	rowRec, err := pg.db.Query(context.Background(), sql, id)
	if err != nil {
		return dto.SubRecordWithIdDTO{}, err
	}
	defer rowRec.Close()
	subRec, err := pg.rowsInSlice(rowRec)
	if err != nil {
		return dto.SubRecordWithIdDTO{}, err
	}
	return subRec[0], nil
}

func (pg *Postgres) GetListSubRecords() ([]dto.SubRecordWithIdDTO, error) {
	sql := `SELECT id, service_name, price, user_id, start_date, end_date 
			FROM subscriptions`
	rowList, err := pg.db.Query(context.Background(), sql)
	if err != nil {
		return []dto.SubRecordWithIdDTO{}, err
	}
	defer rowList.Close()
	subRec, err := pg.rowsInSlice(rowList)
	if err != nil {
		return []dto.SubRecordWithIdDTO{}, err
	}
	return subRec, nil
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

func (pg *Postgres) rowsInSlice(rows pgx.Rows) ([]dto.SubRecordWithIdDTO, error) {
	subRecSlice := []dto.SubRecordWithIdDTO{}
	for rows.Next() {
		var (
			id          int32
			serviceName string
			price       int32
			userID      pgtype.UUID
			startDate   time.Time
			endDate     time.Time
		)
		if err := rows.Scan(&id, &serviceName, &price, &userID, &startDate, &endDate); err != nil {
			return []dto.SubRecordWithIdDTO{}, err
		}
		userIDConv, err := uuid.FromBytes(userID.Bytes[:])
		if err != nil {
			return []dto.SubRecordWithIdDTO{}, err
		}
		data := dto.NewSubRecordWithIdDTO(int(id), serviceName, int(price), userIDConv, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
		subRecSlice = append(subRecSlice, data)
	}

	return subRecSlice, nil
}
