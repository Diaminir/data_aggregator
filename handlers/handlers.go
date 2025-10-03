package handlers

import (
	"junior_effectivemobile/db"

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
