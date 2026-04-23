package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func Init(env string){
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var err error
	log, err = config.Build()
	if err != nil {
		os.Exit(1)
	}
}

func GetLogger() *zap.Logger {
	if log == nil {
		Init("development")		
	}
	return log
}

func Info(msg string, fields ...zap.Field){
	GetLogger().Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field){
	GetLogger().Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field){
	GetLogger().Warn(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field){
	GetLogger().Fatal(msg, fields...)
}

func Debug(msg string, fields ...zap.Field){
	GetLogger().Debug(msg, fields...)
}

func Sync()  {
	if log != nil {
		_ = log.Sync()
	}
}
