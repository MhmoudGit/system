package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppAddr       string
	DbHost        string
	DbPort        string
	DbUser        string
	DbPassword    string
	DbName        string
	JWTSecret     string
	EmailFrom     string
	EmailPassword string
}

func LoadConfig() *Config {
	appAddr := os.Getenv("APP_ADDR")

	if appAddr == "" {
		appAddr = "localhost:8000"
	}

	return &Config{
		AppAddr:       appAddr,
		DbHost:        os.Getenv("DB_HOST"),
		DbPort:        os.Getenv("DB_PORT"),
		DbUser:        os.Getenv("DB_USER"),
		DbPassword:    os.Getenv("DB_PASSWORD"),
		DbName:        os.Getenv("DB_NAME"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		EmailFrom:     os.Getenv("EMAIL_FROM"),
		EmailPassword: os.Getenv("EMAIL_PASSWORD"),
	}
}

func InitEnv() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}
	return nil
}
