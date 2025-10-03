package handlers

import (
	"junior_effectivemobile/db"
	"junior_effectivemobile/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HandlersApp struct {
	db  *db.Postgres
	log *logrus.Logger
}

func NewHandlersApp(db *db.Postgres, log *logrus.Logger) *HandlersApp {
	return &HandlersApp{
		db:  db,
		log: log,
	}
}

// Создание записи о подписке
/*
паттерн: /subscriptions
метод: POST
получает: JSON в теле запроса

успех:
  - код: 201 Create
  - тело ответа: JSON структура записи с ее id в бд

ошибки:
  - код: 400, 500
  - тело ответа: JSON с ошибкой + время
*/
func (app *HandlersApp) NewSubRecord(c *gin.Context) {
	var subRecord dto.SubRecordDTO

	if err := c.BindJSON(&subRecord); err != nil {
		app.log.WithError(err).Error("Неверный тип данных в запросе")
		c.JSON(http.StatusBadRequest, dto.NewMessageDTO("Неверный тип данных в запросе", err))
		return
	}
	subRec, err := app.db.PostNewSubRecord(subRecord)
	if err != nil {
		app.log.WithError(err).Error("Ошибка работы с БД")
		c.JSON(http.StatusInternalServerError, dto.NewMessageDTO("Ошибка работы с БД", err))
		return
	}
	app.log.Infof("Данные успешно записаны: %v", subRec)
	c.JSON(http.StatusCreated, subRec)
}

// Получение записи о подписке по id в бд
/*
паттерн: /subscriptions/:id
метод: GET
получает: id в паттерне

успех:
  - код: 200 ok
  - тело ответа: JSON структура записи с ее id в бд

ошибки:
  - код: 400, 404, 500
  - тело ответа: JSON с ошибкой + время
*/
func (app *HandlersApp) GetUserSubRecord(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		app.log.WithError(err).Error("Неверный id записи")
		c.JSON(http.StatusBadRequest, dto.NewMessageDTO("Неверный id записи", err))
		return
	}
	subRec, err := app.db.GetSubRecord(id)
	if err != nil {
		app.log.WithError(err).Error("Ошибка работы с БД")
		c.JSON(http.StatusInternalServerError, dto.NewMessageDTO("Ошибка работы с БД", err))
		return
	}
	app.log.Infof("Запись успешно получена: %v", subRec)
	c.JSON(http.StatusOK, subRec)
}

// Получение списка подписок
/*
паттерн: /subscriptions/
метод: GET
получает: ничего

успех:
  - код: 200 ok
  - тело ответа: JSON слайс структур записей с их id в бд

ошибки:
  - код: 500
  - тело ответа: JSON с ошибкой + время
*/
func (app *HandlersApp) ListAllSubRecords(c *gin.Context) {
	subRecs, err := app.db.GetListSubRecords()
	if err != nil {
		app.log.WithError(err).Error("Ошибка работы с БД")
		c.JSON(http.StatusInternalServerError, dto.NewMessageDTO("Ошибка работы с БД", err))
		return
	}
	app.log.Infof("Записи успешно получены: %v", subRecs)
	c.JSON(http.StatusOK, subRecs)
}

// Изменение записи об подписке
/*
паттерн: /subscriptions/:id
метод: PUT
получает: JSON с данными для обновления + id задачи в паттерне

успех:
  - код: 200 ok
  - тело ответа: JSON структура измененной записи с ее id в бд

ошибки:
  - код: 400, 404, 500
  - тело ответа: JSON с ошибкой + время
*/
func (app *HandlersApp) UpdateSubRecord(c *gin.Context) {
	var updateSubRecord dto.UpdateSubRecordDTO
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		app.log.WithError(err).Error("Неверный id записи")
		c.JSON(http.StatusBadRequest, dto.NewMessageDTO("Неверный id записи", err))
		return
	}
	if err := c.BindJSON(&updateSubRecord); err != nil {
		app.log.WithError(err).Error("Неверный тип данных в запросе")
		c.JSON(http.StatusBadRequest, dto.NewMessageDTO("Неверный тип данных в запросе", err))
		return
	}
	subRec, err := app.db.UpdateSubRecord(id, updateSubRecord)
	if err != nil {
		app.log.WithError(err).Error("Ошибка работы с БД")
		c.JSON(http.StatusInternalServerError, dto.NewMessageDTO("Ошибка работы с БД", err))
		return
	}
	c.JSON(http.StatusOK, subRec)
}

// Удаляет запись об подписке
/*
паттерн: /subscriptions/:id
метод: DELETE
получает: id в паттерне

успех:
  - код: 200 ok
  - тело ответа: JSON об успешном удалении

ошибки:
  - код: 400, 404, 500
  - тело ответа: JSON с ошибкой + время
*/
func (app *HandlersApp) DeleteSubRecord(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		app.log.WithError(err).Error("Неверный id записи")
		c.JSON(http.StatusBadRequest, dto.NewMessageDTO("Неверный id записи", err))
		return
	}
	if err := app.db.DeleteSubRecord(id); err != nil {
		app.log.WithError(err).Error("Ошибка работы с БД")
		c.JSON(http.StatusInternalServerError, dto.NewMessageDTO("Ошибка работы с БД", err))
		return
	}
	app.log.Info("Запись успешно удалена")
	c.JSON(http.StatusOK, dto.NewMessageDTO("Запись успешно удалена", nil))
}

// Подсчет суммарной стоимости всех подписок за выбранный период и фильтрами
/*
паттерн: /subscriptions/cost
метод: GER
получает: query параметры (id пользователя, название подписки, начало и конец периода)

успех:
  - код: 200 ok
  - тело ответа: JSON структура измененной записи с ее id в бд

ошибки:
  - код: 400, 404, 500
  - тело ответа: JSON с ошибкой + время
*/
func (app *HandlersApp) CalculateCost(c *gin.Context) {
	queryParam, err := dto.NewQueryParam(c)
	if err != nil {
		app.log.WithError(err).Error("Неверный тип данных в запросе")
		c.JSON(http.StatusBadRequest, dto.NewMessageDTO("Неверный тип данных в запросе", err))
		return
	}
	costSummary, err := app.db.CalculateCost(queryParam)
	if err != nil {
		app.log.WithError(err).Error("Ошибка работы с БД")
		c.JSON(http.StatusInternalServerError, dto.NewMessageDTO("Ошибка работы с БД", err))
		return
	}
	c.JSON(http.StatusOK, costSummary)
}
