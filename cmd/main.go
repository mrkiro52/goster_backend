package main

import (
	"fmt"
	"goster/config"
	logging "goster/library/logger"
)

func init() {
	// Загружаем переменные окружения
	config.GetEnvVariables()

	// Конфигурируем логгер
	logging.ConfigureLogger()

	logger := logging.GetLogger("init")
	logger.Info("Env variables are loaded!")
}

func main() {
	fmt.Print("Hello")
}
