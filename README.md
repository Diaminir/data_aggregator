# Сервис для агрегации данных

Этот проект представляет собой REST-сервис написанный на Golang, который принимает название подписки, ее стоимость, id пользователя (в формате uuid) который владеет подпиской, дату начала и конца подписки, сохраняет их в базу данных PostgreSQL и предоставляет CRUDL-методы для взаимодействия с записями и метод который подсчитывает суммарную стоимость всех подписок за выбранный период с фильтрацией по id пользователя и названию подписки.

## Требования

Для запуска этого проекта вам потребуется следующее:
- Docker
- Создать и заполнить по данному шаблону .env файл:

### Пример файла ```.env```
```
#DB
DB_HOST=postgres
DB_PORT=5432
DB_USER=qwert
DB_PASSWORD=12345
DB_NAME=subscriptions

#api/log level = debug/info
API_PORT=8080
LOG_LEVEL=debug
```

## Запуск
Для запуска выполните в терминале команду ```docker-compose up -d```, после чего сервер будет запущен на указанном вами порту (далее localhost:8080).
Для остановки сервера нужно прописать команду ```docker-compose down```.

## REST Методы
Сервис предоставляет следующие REST-методы:

## SWAGGER
```
http://localhost:8080/swagger/index.html#/
```

### Добавление новой записи:

- Метод: POST
- URL: /subscriptions
- Пример запроса:
```
curl --location 'localhost:8080/subscriptions' \
--header 'Content-Type: application/json' \
--data '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "2025-06-25",
    "end_date": "2027-06-25"
}
```
#### Возвращает JSON
```
{
    "id": 1,
    "subRecord": {
        "service_name": "Yandex Plus",
        "price": 400,
        "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
        "start_date": "2025-06-25",
        "end_date": "2027-06-25"
    }
}
```
#### Пример ошибки
```
{
    "message": "Неверный запрос: введите дату окончания подписки",
    "time": "2025-10-04 13:25:11"
}
```

### Удаление по id:

- Метод: DELETE
- URL: /subscriptions/id
- Пример запроса:
```
curl --location --request DELETE 'localhost:8080/subscriptions/1'
```
#### Возвращает JSON
```
{
    "message": "Запись успешно удалена",
    "time": "2025-10-04 13:04:45"
}
```
#### Пример ошибки:
```
{
    "message": "Ошибка работы с БД: id записи не существует: %!s(<nil>)",
    "time": "2025-10-04 13:28:14"
}
```

### Получение записи по id:

- Метод: GET
- URL: /subscriptions/id
- Пример запроса:
```
curl --location 'localhost:8080/subscriptions/2'
```

#### Возвращает JSON:
```
{
    "id": 1,
    "subRecord": {
        "service_name": "Yandex Plus",
        "price": 400,
        "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
        "start_date": "2025-06-25",
        "end_date": "2027-06-25"
    }
}
```
#### Пример ошибки:
```
{
    "message": "Ошибка работы с БД: id записи не существует: %!s(<nil>)",
    "time": "2025-10-04 13:30:24"
}
```

### Получение всех записей:

- Метод: GET
- URL: /subscriptions
- Пример запроса:
```
curl --location 'localhost:8080/subscriptions'
```
#### Возвращает JSON:
```
[
    {
        "id": 1,
        "subRecord": {
            "service_name": "Yandex Plus",
            "price": 400,
            "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
            "start_date": "2025-06-25",
            "end_date": "2027-06-25"
        }
    },
    {
        "id": 2,
        "subRecord": {
            "service_name": "Nintendo Switch",
            "price": 909,
            "user_id": "b5c3d4e5-f6f7-8901-bcde-f12345678901",
            "start_date": "2027-06-01",
            "end_date": "2029-12-31"
        }
    }
]
```
### Частично или полностью изменяет запись по id:

- Метод: PATCH
- URL: /subscriptions/id
- Пример запроса:
```
curl --location --request PATCH 'localhost:8080/subscriptions/2' \
--header 'Content-Type: application/json' \
--data '{
    "service_name": "Nintendo Switch",
    "price": 909,
    "user_id": "b5c3d4e5-f6f7-8901-bcde-f12345678901",
    "start_date": "2027-06-01",
    "end_date": "2029-12-31"
}'
```
#### Возвращает JSON:
```
{
    "id": 1,
    "subRecord": {
        "service_name": "Nintendo Switch",
        "price": 909,
        "user_id": "b5c3d4e5-f6f7-8901-bcde-f12345678901",
        "start_date": "2027-06-01",
        "end_date": "2029-12-31"
    }
}
```
#### Пример ошибки:
```
{
    "message": "Неверный запрос: хотя бы одно из полей должно быть заполнено",
    "time": "2025-10-04 13:37:19"
}
```

### Подсчет суммарной стоимости всех подписок за выбранный период с фильтрацией по id пользователя и названию подписки:

- Метод: GET
- URL: /subscriptions/cost
- Пример запроса:
```
curl --location 'localhost:8080/subscriptions/cost?user_id=b5c3d4e5-f6f7-8901-bcde-f12345678901&service_name=Nintendo%20Switch&start_period=2023-06-25&end_period=2027-12-31'
```
Поддерживает данные фильтрации:
- user_id (в формате uuid) 
- service_name
- start_period (в формате YYYY-MM-DD)
- end_period (в формате YYYY-MM-DD)
#### Возвращает JSON:
```
{
    "total_cost": 909,
    "query_param": {
        "user_id": "b5c3d4e5-f6f7-8901-bcde-f12345678901",
        "service_name": "Nintendo Switch",
        "start_period": "2023-06-25",
        "end_period": "2027-12-31"
    }
}
```
#### Пример ошибки:
```
{
    "message": "Неверный запрос: введите дату начала периода поиска",
    "time": "2025-10-04 13:40:36"
}
```
