package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTP     HTTPConfig
	Postgres PostgresConfig
	JWT      JWTConfig
}

type HTTPConfig struct {
	Port            string
	ShutdownTimeout time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	SecretKey string
	TTL       time.Duration
}

func New() (*Config, error) {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	shutdownTimeout, err := strconv.Atoi(os.Getenv("HTTP_SHUTDOWN_TIMEOUT"))
	if err != nil {
		shutdownTimeout = 5
	}

	readTimeout, err := strconv.Atoi(os.Getenv("HTTP_READ_TIMEOUT"))
	if err != nil {
		readTimeout = 5
	}

	writeTimeout, err := strconv.Atoi(os.Getenv("HTTP_WRITE_TIMEOUT"))
	if err != nil {
		writeTimeout = 5
	}

	jwtTTL, err := strconv.Atoi(os.Getenv("JWT_TTL"))
	if err != nil {
		jwtTTL = 24
	}

	return &Config{
		HTTP: HTTPConfig{
			Port:            httpPort,
			ShutdownTimeout: time.Duration(shutdownTimeout) * time.Second,
			ReadTimeout:     time.Duration(readTimeout) * time.Second,
			WriteTimeout:    time.Duration(writeTimeout) * time.Second,
		},
		Postgres: PostgresConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DBName:   os.Getenv("POSTGRES_DB"),
			SSLMode:  os.Getenv("POSTGRES_SSL_MODE"),
		},
		JWT: JWTConfig{
			SecretKey: os.Getenv("JWT_SECRET_KEY"),
			TTL:       time.Duration(jwtTTL) * time.Hour,
		},
	}, nil
}

func (c *PostgresConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}
