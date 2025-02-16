# Этап сборки
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Установка зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Финальный этап
FROM alpine:latest

WORKDIR /app

# Копирование бинарного файла из этапа сборки
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

# Установка необходимых зависимостей
RUN apk --no-cache add ca-certificates

EXPOSE 8080

CMD ["./main"] 