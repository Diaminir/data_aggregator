package db

import (
	"context"
	"fmt"
	"junior_effectivemobile/dto"
	"net/http"
	"strings"
)

func (pg *Postgres) PostNewSubRecord(ctx context.Context, rec dto.SubRecordDTO) (dto.SubRecordWithIdDTO, *ErrorDB) {
	pg.log.Debug("Начало создания новой записи подписки: ", rec)

	sql := `INSERT INTO subscriptions(service_name, price, user_id, start_date, end_date) 
			VALUES ($1,$2,$3,$4,$5)
			RETURNING *`

	rowRec, err := pg.db.Query(ctx, sql, rec.ServiceName, rec.Price, rec.UserID, rec.StartDate, rec.EndDate)
	if err != nil {
		return dto.SubRecordWithIdDTO{}, NewErrorDB(fmt.Errorf("ошибка вставки новой записи в Postgres: %s", err), http.StatusInternalServerError)
	}
	defer rowRec.Close()

	subRec, errDB := pg.rowsInSlice(rowRec)
	if errDB != nil {
		return dto.SubRecordWithIdDTO{}, errDB
	}

	pg.log.Debug("Данные в бд записаны: ", subRec[0])
	return subRec[0], nil
}

// GetSubRecord получает запись по ее id в базе данных
func (pg *Postgres) GetSubRecord(ctx context.Context, id int) (dto.SubRecordWithIdDTO, *ErrorDB) {
	pg.log.Debugf("Получение записи от базы данных по id = %d", id)

	if errDB := pg.validateId(ctx, id); errDB != nil {
		return dto.SubRecordWithIdDTO{}, errDB
	}

	sql := `SELECT id, service_name, price, user_id, start_date, end_date 
			FROM subscriptions 
			WHERE id = $1`

	rowRec, err := pg.db.Query(ctx, sql, id)
	if err != nil {
		return dto.SubRecordWithIdDTO{}, NewErrorDB(fmt.Errorf("ошибка получения данных из Postgres: %s", err), http.StatusInternalServerError)
	}
	defer rowRec.Close()

	subRec, errDB := pg.rowsInSlice(rowRec)
	if errDB != nil {
		return dto.SubRecordWithIdDTO{}, errDB
	}

	pg.log.Debug("Данные из бд получены: ", subRec[0])
	return subRec[0], nil
}

// GetListSubRecords получает все записи из базы данных
func (pg *Postgres) GetListSubRecords(ctx context.Context) ([]dto.SubRecordWithIdDTO, *ErrorDB) {
	pg.log.Debug("Получение списка записей от базы данных")

	sql := `SELECT id, service_name, price, user_id, start_date, end_date 
			FROM subscriptions`

	rowList, err := pg.db.Query(ctx, sql)
	if err != nil {
		return []dto.SubRecordWithIdDTO{}, NewErrorDB(fmt.Errorf("ошибка получения всех данных из Postgres: %s", err), http.StatusInternalServerError)
	}
	defer rowList.Close()

	subRec, errDB := pg.rowsInSlice(rowList)
	if errDB != nil {
		return []dto.SubRecordWithIdDTO{}, errDB
	}

	pg.log.Debugf("Данные из бд получены: количество записей = %d", len(subRec))
	return subRec, nil
}

// UpdateSubRecord позволяет редактировать запись по id, меняет данные на данные в структуре updateData и возвращает измененную запись
func (pg *Postgres) UpdateSubRecord(ctx context.Context, id int, updateData dto.UpdateSubRecordDTO) (dto.SubRecordWithIdDTO, *ErrorDB) {
	pg.log.Debugf("Начало редактирования записи в базе данных id = %d, updateData = %v", id, updateData)

	if errDB := pg.validateId(ctx, id); errDB != nil {
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

	rowUpdate, err := pg.db.Query(ctx, sql, args...)
	if err != nil {
		return dto.SubRecordWithIdDTO{}, NewErrorDB(fmt.Errorf("ошибка получения измененных данных из Postgres: %s", err), http.StatusInternalServerError)
	}
	defer rowUpdate.Close()

	subRec, errDB := pg.rowsInSlice(rowUpdate)
	if errDB != nil {
		return dto.SubRecordWithIdDTO{}, errDB
	}

	pg.log.Debugf("Редактирование записи завершено успешно id = %d, обновленная задача = %v", id, subRec[0])
	return subRec[0], nil
}

// DeleteSubRecord удаляет запись из базы данных по ее id
func (pg *Postgres) DeleteSubRecord(ctx context.Context, id int) *ErrorDB {
	pg.log.Debugf("Начало удаления записи из базы данных id = %d", id)

	if err := pg.validateId(ctx, id); err != nil {
		return err
	}

	sql := `DELETE FROM subscriptions 
			WHERE id = $1`

	if _, err := pg.db.Exec(ctx, sql, id); err != nil {
		return NewErrorDB(fmt.Errorf("ошибка удаления данных из Postgres: %s", err), http.StatusInternalServerError)
	}

	pg.log.Debugf("Запись успешно удалена id = %d", id)
	return nil
}
