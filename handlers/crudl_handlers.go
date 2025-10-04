package handlers

import (
	"junior_effectivemobile/db"
	"junior_effectivemobile/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// HandlersApp cтруктура для наследования методов
type HandlersApp struct {
	db  *db.Postgres
	log *logrus.Logger
}

// NewHandlersApp cоздает экземпляр структуры HandlerApp
func NewHandlersApp(db *db.Postgres, log *logrus.Logger) *HandlersApp {
	return &HandlersApp{
		db:  db,
		log: log,
	}
}

// NewSubRecord создает новую запись о подписке
// @Summary Создание новой записи
// @Description Создает новую запись о подписке пользователя
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body dto.SubRecordDTO true "Данные о подписке"
// @Success 201 {object} dto.SubRecordWithIdDTO
// @Failure 400 {object} dto.MessageDTO
// @Failure 500 {object} dto.MessageDTO
// @Router /subscriptions [post]
func (app *HandlersApp) NewSubRecord(c *gin.Context) {
	app.log.Debug("Начало создания новой записи подписки")
	var subRecord dto.SubRecordDTO

	if err := c.BindJSON(&subRecord); err != nil {
		app.log.WithError(err).Error(ErrorInvalidType)
		c.JSON(http.StatusBadRequest, dto.NewMessageDTO(ErrorInvalidType, err))
		return
	}

	app.log.Debug("Валидация входных данных")
	if err := subRecord.ValidateInputData(); err != nil {
		app.log.WithError(err).Error(ErrorInvalidRequest)
		c.JSON(http.StatusBadRequest, dto.NewMessageDTO(ErrorInvalidRequest, err))
		return
	}

	app.log.Debug("Отправка данных в базу данных для создания записи")
	subRec, errDB := app.db.PostNewSubRecord(c.Request.Context(), subRecord)
	if errDB != nil {
		app.log.WithError(errDB.Err).Error(ErrorDB)
		c.JSON(errDB.Code, dto.NewMessageDTO(ErrorDB, errDB.Err))
		return
	}
	app.log.Infof("Данные успешно записаны: %v", subRec)
	c.JSON(http.StatusCreated, subRec)
}

// GetUserSubRecord получает запись об подписке по ее id
// @Summary Запись по ID
// @Description Получает запись о подписке по ее ID
// @Tags subscriptions
// @Produce json
// @Param id path int true "ID записи"
// @Success 200 {object} dto.SubRecordWithIdDTO
// @Failure 400 {object} dto.MessageDTO
// @Failure 404 {object} dto.MessageDTO
// @Failure 500 {object} dto.MessageDTO
// @Router /subscriptions/{id} [get]
func (app *HandlersApp) GetUserSubRecord(c *gin.Context) {
	app.log.Debug("Начало получения записи подписки по ID")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		app.log.WithError(err).Error(ErrorInvalidRecordId)
		c.JSON(http.StatusBadRequest, dto.NewMessageDTO(ErrorInvalidRecordId, err))
		return
	}

	app.log.Debugf("Запрос записи из базы данных по ID: id = %d", id)
	subRec, errDB := app.db.GetSubRecord(c.Request.Context(), id)
	if errDB != nil {
		app.log.WithError(errDB.Err).Error(ErrorDB)
		c.JSON(errDB.Code, dto.NewMessageDTO(ErrorDB, errDB.Err))
		return
	}
	app.log.Infof("Запись успешно получена: %v", subRec)
	c.JSON(http.StatusOK, subRec)
}

// ListAllSubRecords получает все записи о подписках
// @Summary Лист всех записей
// @Description Получает все записи об подписках
// @Tags subscriptions
// @Produce json
// @Success 200 {object} dto.SubRecordWithIdDTO
// @Failure 500 {object} dto.MessageDTO
// @Router /subscriptions [get]
func (app *HandlersApp) ListAllSubRecords(c *gin.Context) {
	app.log.Debug("Начало получения всех записей подписок")
	subRecs, errDB := app.db.GetListSubRecords(c.Request.Context())
	if errDB != nil {
		app.log.WithError(errDB.Err).Error(ErrorDB)
		c.JSON(errDB.Code, dto.NewMessageDTO(ErrorDB, errDB.Err))
		return
	}
	app.log.Infof("Записи успешно получены: %v", subRecs)
	c.JSON(http.StatusOK, subRecs)
}

// UpdateSubRecord обновляет данные о подписке
// @Summary Изменение записи
// @Description Изменение данных о подписке в записе
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "ID записи"
// @Param subscription body dto.UpdateSubRecordDTO true "Данные для обновление записи"
// @Success 200 {object} dto.SubRecordWithIdDTO
// @Failure 400 {object} dto.MessageDTO
// @Failure 404 {object} dto.MessageDTO
// @Failure 500 {object} dto.MessageDTO
// @Router /subscriptions/{id} [put]
func (app *HandlersApp) UpdateSubRecord(c *gin.Context) {
	app.log.Debug("Начало обновления записи подписки")
	var updateSubRecord dto.UpdateSubRecordDTO
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		app.log.WithError(err).Error(ErrorInvalidRecordId)
		c.JSON(http.StatusBadRequest, dto.NewMessageDTO(ErrorInvalidRecordId, err))
		return
	}
	if err := c.BindJSON(&updateSubRecord); err != nil {
		app.log.WithError(err).Error(ErrorInvalidType)
		c.JSON(http.StatusBadRequest, dto.NewMessageDTO(ErrorInvalidType, err))
		return
	}
	app.log.Debug("Валидация данных для обновления")
	if err := updateSubRecord.ValidateUpdateData(); err != nil {
		app.log.WithError(err).Error(ErrorInvalidRequest)
		c.JSON(http.StatusBadRequest, dto.NewMessageDTO(ErrorInvalidRequest, err))
		return
	}
	app.log.Debugf("Отправка данных в базу данных для обновления записи: id = %d", id)
	subRec, errDB := app.db.UpdateSubRecord(c.Request.Context(), id, updateSubRecord)
	if errDB != nil {
		app.log.WithError(errDB.Err).Error(ErrorDB)
		c.JSON(errDB.Code, dto.NewMessageDTO(ErrorDB, errDB.Err))
		return
	}
	app.log.Infof("Обновление записи подписки успешно завершено: %v", subRec)
	c.JSON(http.StatusOK, subRec)
}

// DeleteSubRecord удаляет запись о подписке
// @Summary Удаление записи о подписке
// @Description Удаление записи о подписке по ID
// @Tags subscriptions
// @Param id path int true "ID записи"
// @Success 200 {object} dto.MessageDTO
// @Failure 400 {object} dto.MessageDTO
// @Failure 404 {object} dto.MessageDTO
// @Failure 500 {object} dto.MessageDTO
// @Router /subscriptions/{id} [delete]
func (app *HandlersApp) DeleteSubRecord(c *gin.Context) {
	app.log.Debug("Начало удаления записи подписки")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		app.log.WithError(err).Error(ErrorInvalidRecordId)
		c.JSON(http.StatusBadRequest, dto.NewMessageDTO(ErrorInvalidRecordId, err))
		return
	}
	app.log.Debugf("Отправка запроса на удаление в базу данных: id = %d", id)
	if errDB := app.db.DeleteSubRecord(c.Request.Context(), id); errDB != nil {
		app.log.WithError(errDB.Err).Error(ErrorDB)
		c.JSON(errDB.Code, dto.NewMessageDTO(ErrorDB, errDB.Err))
		return
	}
	app.log.Info("Запись успешно удалена")
	c.JSON(http.StatusOK, dto.NewMessageDTO("Запись успешно удалена", nil))
}
