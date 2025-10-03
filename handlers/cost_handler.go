package handlers

import (
	"junior_effectivemobile/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CalculateCost подсчет суммарной стоимости всех подписок за выбранный период и фильтрами
// @Summary Подсчет суммарной стоимости подписок
// @Description Подсчет суммарной стоимости всех подписок за выбранный период и опциональными фильтрами
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "Фильтр User ID"
// @Param service_name query string false "Фильтр названия подписки"
// @Param start_period query string true "Начальная дата (YYYY-MM-DD)" format(string)
// @Param end_period query string true "Конечная дата (YYYY-MM-DD)" format(string)
// @Success 200 {object} dto.CostSummaryRespDTO
// @Failure 400 {object} dto.MessageDTO
// @Failure 500 {object} dto.MessageDTO
// @Router /subscriptions/cost [get]
func (app *HandlersApp) CalculateCost(c *gin.Context) {
	queryParam, err := dto.NewQueryParam(c)
	if err != nil {
		app.log.WithError(err).Error(ErrorInvalidRequest)
		c.JSON(http.StatusBadRequest, dto.NewMessageDTO(ErrorInvalidRequest, err))
		return
	}
	costSummary, errDB := app.db.CalculateCost(c.Request.Context(), queryParam)
	if errDB != nil {
		app.log.WithError(errDB.Err).Error(ErrorDB)
		c.JSON(errDB.Code, dto.NewMessageDTO(ErrorDB, err))
		return
	}
	c.JSON(http.StatusOK, costSummary)
}
