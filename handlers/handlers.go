package handlers

import (
	"junior_effectivemobile/db"
	"junior_effectivemobile/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HandlersApp struct {
	db *db.Postgres
}

func NewHandlersApp(db *db.Postgres) *HandlersApp {
	return &HandlersApp{
		db: db,
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
		c.JSON(http.StatusBadRequest, dto.NewMessageDTO("Неверный тип данных в запросе", err))
		return
	}
	subRec, err := app.db.PostNewSubRecord(subRecord)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageDTO("Ошибка работы с БД", err))
		return
	}
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
		c.JSON(http.StatusBadRequest, dto.NewMessageDTO("Неверный id записи", err))
		return
	}
	subRec, err := app.db.GetSubRecord(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageDTO("Ошибка работы с БД", err))
		return
	}
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
		c.JSON(http.StatusInternalServerError, dto.NewMessageDTO("Ошибка работы с БД", err))
		return
	}
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
		c.JSON(http.StatusBadRequest, dto.NewMessageDTO("Неверный id записи", err))
		return
	}
	if err := c.BindJSON(&updateSubRecord); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewMessageDTO("Неверный тип данных в запросе", err))
		return
	}
	subRec, err := app.db.UpdateSubRecord(id, updateSubRecord)
	if err != nil {
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
		c.JSON(http.StatusBadRequest, dto.NewMessageDTO("Неверный id записи", err))
		return
	}
	if err := app.db.DeleteSubRecord(id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewMessageDTO("Ошибка работы с БД", err))
		return
	}
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

}
