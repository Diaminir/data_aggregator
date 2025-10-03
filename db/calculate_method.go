package db

import (
	"context"
	"fmt"
	"junior_effectivemobile/dto"
	"net/http"
)

// CalculateCost производит подсчет суммарной стоимости всех подписок с применением фильтров из структуры queryParam и возвращает структуру с суммарной стоимость и фильтрами
func (pg *Postgres) CalculateCost(ctx context.Context, queryParam dto.CostSummaryReqDTO) (dto.CostSummaryRespDTO, *ErrorDB) {
	pg.log.Debug("Данные от хэндлера получены", queryParam)
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

	if err := pg.db.QueryRow(ctx, sql, args...).Scan(&totalCost); err != nil {
		return dto.CostSummaryRespDTO{}, NewErrorDB(fmt.Errorf("ошибка получения суммарной стоимости подписок из Postgres: %s", err), http.StatusInternalServerError)
	}

	pg.log.Debug("Формирование ответа с результатами расчета")
	result := dto.NewCostSummaryRespDTO(totalCost, queryParam)
	return result, nil
}
