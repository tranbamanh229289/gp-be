package logger

import "go.uber.org/zap"

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

type ZapLogger struct {
	logger *zap.Logger
}

func (zl *ZapLogger) Debug(msg string, fields ...zap.Field) {
	zl.logger.Debug(msg, fields...)
}

func (zl *ZapLogger) Info(msg string, fields ...zap.Field){
	zl.logger.Info(msg, fields...)
}

func (zl *ZapLogger) Warn(msg string, fields ...zap.Field) {
	zl.logger.Warn(msg, fields...)
}

func (zl *ZapLogger) Error(msg string, fields ...zap.Field){
	zl.logger.Error(msg, fields...)
}

func(zl *ZapLogger) Fatal(msg string, fields ...zap.Field) {
	zl.logger.Fatal(msg, fields...)
}
