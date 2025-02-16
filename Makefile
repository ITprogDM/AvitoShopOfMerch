
# Загружаем переменные из .env
include .env

# Определяем URL базы данных в зависимости от окружения (контейнер или локально)
ifeq ($(shell docker ps -q -f name=shopavito-db-1),)
    DB_URL := $(PG_URL_LOCALHOST)
else
    DB_URL := $(PG_URL_CONTAINER)
endif

# Контейнеры
DOCKER_COMPOSE := docker-compose
GOLANGCI_LINT := golangci-lint

.PHONY: build run stop restart clean test lint fmt migrate db-reset db-up db-down integration-test


#Запустить приложение(с Docker)
start:
	@echo "Starting containers..."
	$(DOCKER_COMPOSE) up --build -d

# Остановить контейнеры
stop:
	@echo "Stopping containers..."
	$(DOCKER_COMPOSE) down

# Перезапустить контейнеры
restart: stop
	@echo "Restarting containers..."
	$(DOCKER_COMPOSE) up --build -d

# Удалить скомпилированные файлы
clean:
	@echo "Cleaning up..."
	rm -rf shop

# Запуск всех тестов (unit + integration)
test:
	@echo "Running tests..."
	go test ./tests/... -cover

# Запуск unit-тестов отдельно
unit-test:
	@echo "Running unit tests..."
	go test ./tests/unit/... -cover

# Запуск integration-тестов отдельно (с поднятием контейнеров)
integration-test: db-up
	@echo "Running integration tests..."
	go test ./tests/integration/... -cover

# Запуск линтера
lint:
	@echo "Running golangci-lint..."
	$(GOLANGCI_LINT) run

# Применить миграции (
migrate:
	@echo "Applying migrations..."
	migrate -path ./migrations -database "$(DB_URL)" up

# Откатить миграции
migrate-down:
	@echo "Rolling back migrations..."
	migrate -path ./migrations -database "$(DB_URL)" down



