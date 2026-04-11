package logging

import (
	"goster/config"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func ConfigureLogger() {
	// Функция для конфигурации логгера

	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})

	level, err := logrus.ParseLevel(config.LogLevel)

	if err != nil {
		level = logrus.InfoLevel
	}

	logrus.SetLevel(level)

	logrus.SetOutput(os.Stdout)
}

func GetLogger(function string) *logrus.Entry {
	// Функция чтобы в логе указывалась функция пишущая логи
	return logrus.WithField("function", function)
}
