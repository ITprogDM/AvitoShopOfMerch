Avito Shop



Микросервис для управления балансом пользователей, покупок, инвентаря и перевода монет между пользователями.

Используемые технологии

PostgreSQL (в качестве Базы Данных)

Docker (для запуска сервиса)

Gin (веб-фреймворк)

golang-migrate/migrate (для миграций БД)

pgx (драйвер для работы с PostgreSQL)

golang/mock, testify (для тестирования, в некоторых местах subs)

Сервис был написан с использованием Clean Architecture, что позволяет легко расширять функционал и тестировать сервис. Также реализован Graceful Shutdown для корректного завершения работы.

Getting Started

Для запуска сервиса необходимо:

Заполнить .env файл на основе .env.example.

Применить миграции БД:

make migrate

Запустить сервис:

make start

Для полного тестирования микросервиса, сначала нужно запустить контейнеры!

Usage

API Эндпоинты

Аутентификация

POST /api/auth — логин или регистрация (JWT-токен)

Пользователь

GET /api/info — информация о пользователе (баланс, инвентарь, транзакции)

Покупки

GET /api/buy/:item — покупка товара

Транзакции

POST /api/sendCoin — перевод монет другому пользователю

Тестирование и линтинг

Запуск тестов: make test

Запуск тестов с покрытием: make cover

Запуск линтера: make lint

Примеры запросов

Регистрация

curl --location --request POST 'http://localhost:8080/api/auth' \
--header 'Content-Type: application/json' \
--data-raw '{ "username": "user", "password": "password" }'

Пример ответа:

{ "token": "your_jwt_token" }

Получение информации о пользователе

curl --location --request GET 'http://localhost:8080/api/info' \
--header 'Authorization: Bearer your_jwt_token'

Пример ответа:

{
"coins": 1000,
"inventory": [
{ "type": "t-shirt", "quantity": 2 },
{ "type": "book", "quantity": 1 }
],
"coinHistory": {
"received": [{ "fromUser": "friend", "amount": 200 }],
"sent": [{ "toUser": "shop", "amount": 50 }]
}
}

Покупка товара

curl --location --request GET 'http://localhost:8080/api/buy/t-shirt' \
--header 'Authorization: Bearer your_jwt_token'

Пример ответа:

{ "message": "Purchase successful" }

Перевод монет

curl --location --request POST 'http://localhost:8080/api/sendCoin' \
--header 'Authorization: Bearer your_jwt_token' \
--header 'Content-Type: application/json' \
--data-raw '{ "toUser": "receiver", "amount": 100 }'

Пример ответа:

{ "message": "Transaction successful" }