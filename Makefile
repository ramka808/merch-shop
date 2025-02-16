.PHONY: build run test migrate-up migrate-down

# Сборка приложения
build:
	go build -o bin/api cmd/api/main.go

# Запуск приложения
run:
	go run cmd/api/main.go

# Запуск тестов
test:
	go test -v ./...

# Запуск интеграционных тестов
test-integration:
	go test -v -tags=integration ./...

# Запуск E2E тестов
test-e2e:
	go test -v -tags=e2e ./...

# Применение миграций
migrate-up:
	goose -dir migrations postgres "$(DB_DSN)" up

# Откат миграций
migrate-down:
	goose -dir migrations postgres "$(DB_DSN)" down

# Запуск линтера
lint:
	golangci-lint run

# Генерация Swagger документации
swagger:
	swag init -g cmd/api/main.go -o docs

# Запуск в Docker
docker-run:
	docker-compose up --build -d

# Остановка Docker контейнеров
docker-stop:
	docker-compose down

# Очистка
clean:
	rm -rf bin/
	go clean 