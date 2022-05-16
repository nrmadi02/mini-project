package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm/logger"
	"time"
)

type MyWriter struct {
	log *logrus.Logger
}

func (m *MyWriter) Printf(format string, v ...interface{}) {
	strlog := fmt.Sprintf(format, v...)
	m.log.Info(strlog)
}

func NewMyWriter(connectMongo *mongo.Database) *MyWriter {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	mw := NewMongoWriter(connectMongo)
	log.SetOutput(mw)
	return &MyWriter{log: log}
}

func SlowLoggerGorm(cm *mongo.Database) logger.Interface {
	slowLogger := logger.New(
		NewMyWriter(cm),
		logger.Config{
			SlowThreshold: time.Millisecond,
			LogLevel:      logger.Warn,
		},
	)
	return slowLogger
}
