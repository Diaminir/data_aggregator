package db

import (
	"context"
	"fmt"
	"junior_effectivemobile/dto"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// validateId проверяет наличие введенного id записи, в случает отсутствия возвращает ошибку
func (pg *Postgres) validateId(ctx context.Context, id int) *ErrorDB {
	pg.log.Debugf("Проверка наличия записи: id = %d", id)
	sqlValidation := `SELECT id 
					  FROM subscriptions WHERE id = $1`
	tag, err := pg.db.Exec(ctx, sqlValidation, id)
	if err != nil || tag.RowsAffected() == 0 {
		return NewErrorDB(fmt.Errorf("id записи не существует: %s", err), http.StatusNotFound)
	}
	pg.log.Debugf("Запись обнаружена id = %d", id)
	return nil
}

// rowsInSlice получает ответ от базы данных в формате pgx.Rows, вычитывает записи из него и заполняет слайс структур для вывода записей пользователю
func (pg *Postgres) rowsInSlice(rows pgx.Rows) ([]dto.SubRecordWithIdDTO, *ErrorDB) {
	pg.log.Debug("Обработка полученныйх строк и запись в слайс")
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
			return []dto.SubRecordWithIdDTO{}, NewErrorDB(fmt.Errorf("ошибка получения данных из Postgres: %s", err), http.StatusInternalServerError)
		}
		userIDConv, err := uuid.FromBytes(userID.Bytes[:])
		if err != nil {
			return []dto.SubRecordWithIdDTO{}, NewErrorDB(fmt.Errorf("ошибка преобразования pgtype.UUID в uuid.UUID: %s", err), http.StatusInternalServerError)
		}
		data := dto.NewSubRecordWithIdDTO(int(id), serviceName, int(price), userIDConv, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
		subRecSlice = append(subRecSlice, data)
	}
	pg.log.Debug("Запись полученых строк в слайс прошла успешно")
	return subRecSlice, nil
}
