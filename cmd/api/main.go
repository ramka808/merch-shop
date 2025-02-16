package main

import (
	"log"
	"os"

	"github.com/avito/internal/app"
	"github.com/joho/godotenv"
)

func init() {
	// Загрузка переменных окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found")
	}
}

func main() {
	// Инициализация приложения
	app, err := app.NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	// Получение порта из переменных окружения или использование порта по умолчанию
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Запуск сервера
	if err := app.Run(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
