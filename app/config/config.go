package config

import (
	"os"
)

var (
	// JWT секрет
	JWTsecret string

	// Уровень логгирования
	LogLevel string

	// URI для подключения к бд
	DBuri string
)

func GetEnvVariables() {
	// Функция для загрузки переменных окружения

	JWTsecret = os.Getenv("JWT")
	if JWTsecret == "" {
		panic("JWTsecret variable is required")
	}

	LogLevel = os.Getenv("LOG_LEVEL")
	if LogLevel == "" {
		LogLevel = "info"
	}

	DBuri = os.Getenv("DB_URI")
	if DBuri == "" {
		panic("DBuri variable is required")
	}
}
