package db

import (
	"context"
	"fmt"
	"junior_effectivemobile/config"
	"junior_effectivemobile/dto"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sirupsen/logrus"
)

type Postgres struct {
	db  *pgx.Conn
	log *logrus.Logger
}

func NewConPostgres(log *logrus.Logger, cfg *config.Config) (*Postgres, error) {
	log.Debug("Подключение к серверу Postgres...")
	connFig := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName)

	conn, err := pgx.Connect(context.Background(), connFig)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к Postgres: %s", err)
	}

	log.Debug("Начало миграции")
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
	log.Debug("Конец миграции")
	return &Postgres{
		db:  conn,
		log: log,
	}, nil
}

func (pg *Postgres) DbClose() {
	pg.log.Debug("Отключение от сервера Postgres...")
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
	pg.log.Debug("Получение списка записей от базы данных")
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
	pg.log.Debug("Редактиворание записи в базе данных по ID...", id, updateData)

	if errDB := pg.validateId(id); errDB != nil {
		return dto.SubRecordWithIdDTO{}, errDB
	}
	var sql string
	var conditions []string
	var args []interface{}
	argPos := 1
	if updateData.ServiceName != nil {
		conditions = append(conditions, fmt.Sprintf("service_name = $%d", argPos))
		args = append(args, *updateData.ServiceName)
		argPos++
	}
	if updateData.Price != nil {
		conditions = append(conditions, fmt.Sprintf("price = $%d", argPos))
		args = append(args, *updateData.Price)
		argPos++
	}
	if updateData.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argPos))
		args = append(args, *updateData.UserID)
		argPos++
	}
	if updateData.StartDate != nil {
		conditions = append(conditions, fmt.Sprintf("start_date = $%d", argPos))
		args = append(args, *updateData.StartDate)
		argPos++
	}
	if updateData.EndDate != nil {
		conditions = append(conditions, fmt.Sprintf("end_date = $%d", argPos))
		args = append(args, *updateData.EndDate)
		argPos++
	}
	args = append(args, id)
	end := fmt.Sprintf(" WHERE id = $%d RETURNING *", argPos)
	if len(conditions) > 0 {
		sql += "UPDATE subscriptions SET " + strings.Join(conditions, ", ") + end
	}
	rowUpdate, err := pg.db.Query(context.Background(), sql, args...)
	if err != nil {
		return dto.SubRecordWithIdDTO{}, err
	}
	defer rowUpdate.Close()
	subRec, errDB := pg.rowsInSlice(rowUpdate)
	if errDB != nil {
		return dto.SubRecordWithIdDTO{}, errDB
	}
	return subRec[0], nil
}

func (pg *Postgres) DeleteSubRecord(id int) error {
	pg.log.Debug("Удаление записи из базы данных по ID...", id)

	if err := pg.validateId(id); err != nil {
		return err
	}
	sql := `DELETE FROM subscriptions 
			WHERE id = $1`

	if _, err := pg.db.Exec(context.Background(), sql, id); err != nil {
		return err
	}
	return nil
}

func (pg *Postgres) CalculateCost(queryParam dto.CostSummaryReqDTO) (dto.CostSummaryRespDTO, error) {
	sql := `SELECT COALESCE(SUM(price), 0)
       		FROM subscriptions
        	WHERE start_date <= $1 
        	AND start_date >= $2`
	args := []interface{}{queryParam.EndPeriod, queryParam.StartPeriod}
	argPos := 3
	if queryParam.UserID != nil {
		sql += fmt.Sprintf(" AND user_id = $%d", argPos)
		args = append(args, *queryParam.UserID)
		argPos++
	}
	if queryParam.ServiceName != nil {
		sql += fmt.Sprintf(" AND service_name = $%d", argPos)
		args = append(args, *queryParam.ServiceName)
	}
	var totalCost int
	if err := pg.db.QueryRow(context.Background(), sql, args...).Scan(&totalCost); err != nil {
		return dto.CostSummaryRespDTO{}, err
	}
	return dto.NewCostSummaryRespDTO(totalCost, queryParam), nil
}

func (pg *Postgres) validateId(id int) error {
	pg.log.Debug("Проверка наличия записи...", id)

	sqlValidation := `SELECT id 
					  FROM subscriptions WHERE id = $1`

	tag, err := pg.db.Exec(context.Background(), sqlValidation, id)
	if err != nil || tag.RowsAffected() == 0 {
		return err
	}

	return nil
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
