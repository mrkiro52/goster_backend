package main

import (
	"context"
	"errors"
	"fmt"
	"goster/applications/user"
	"goster/config"
	infra "goster/infrastructure"
	controllers "goster/infrastructure/controllers"
	gormRepo "goster/infrastructure/gorm"
	logging "goster/library/logger"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	logger := logging.GetLogger("main")

	// Подключаемся к БД
	db := infra.ConnectToDb()

	// Выполняем миграции
	if err := infra.AutoMigrate(db); err != nil {
		logger.Fatalf("Failed to migrate database: %v", err)
	}
	logger.Info("Database migrations completed")

	// Инициализируем репозитории
	userRepo := gormRepo.NewUserRepository(db)

	// Инициализируем фабрику транзакций
	txFactory := func(ctx context.Context) *gormRepo.Context {
		return gormRepo.WithTransaction(ctx, db)
	}

	// Инициализируем сервисы
	userService := user.NewService(userRepo, txFactory)

	// Инициализируем контроллеры
	userController := controllers.NewUserController(userService)

	// Настраиваем роутер
	r := gin.Default()

	// Настройка CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.ForwardedByClientIP = true
	if err := r.SetTrustedProxies(nil); err != nil {
		panic("SetTrustedProxies failed:" + err.Error())
	}

	// Настраиваем роуты
	controllers.SetupRoutes(r, userController)

	// Настройка HTTP сервера
	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Info("Starting server on :8080")

	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(fmt.Sprintf("Failed to start Service: %s", err))
	}
}
