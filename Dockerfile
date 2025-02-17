# Этап сборки
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Установка зависимостей и Goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Копирование исходного кода и миграций
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Финальный этап
FROM alpine:latest

WORKDIR /app

# Копирование бинарного файла и миграций из этапа сборки
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Установка необходимых зависимостей
RUN apk --no-cache add ca-certificates

# Создаем скрипт для запуска
COPY scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"] 