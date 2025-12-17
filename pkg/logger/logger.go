package logger

import (
	"be/config"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.Logger
}

func NewLogger(cfg *config.Config) (*ZapLogger, error) {
	var level, err = zapcore.ParseLevel(cfg.Zap.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level %w", cfg.Zap.Level)
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "function",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var cores []zapcore.Core

	//Console cores
	consoleSyncer := zapcore.AddSync(os.Stdout)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleSyncer, level)
	cores = append(cores, consoleCore)

	// // //Fluent cores
	// if fluent != nil {
	// 		fluentEncoder := zapcore.NewJSONEncoder(encoderConfig)
	// 		fluentCore := zapcore.NewCore(fluentEncoder, fluent, level)
	// 		cores = append(cores, fluentCore)
	// }

	logger := zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddCallerSkip(0),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	logger = logger.With(
		zap.String("app_name", cfg.App.Name),
		zap.String("version", cfg.App.Version),
		zap.String("environment", cfg.App.Environment),
	)
	return &ZapLogger{logger}, nil
}

func (zl *ZapLogger) GetLogger() *zap.Logger {
	return zl.logger
}
func (zl *ZapLogger) Sync() error {
	return zl.logger.Sync()
}

func (zl *ZapLogger) Info(msg string, fields ...zap.Field) {
	zl.logger.Info(msg, fields...)
}

func (zl *ZapLogger) Debug(msg string, fields ...zap.Field) {
	zl.logger.Debug(msg, fields...)
}

func (zl *ZapLogger) Warn(msg string, fields ...zap.Field) {
	zl.logger.Warn(msg, fields...)
}

func (zl *ZapLogger) Error(msg string, fields ...zap.Field) {
	zl.logger.Error(msg, fields...)
}

func (zl *ZapLogger) DPanic(msg string, fields ...zap.Field) {
	zl.logger.DPanic(msg, fields...)
}

func (zl *ZapLogger) Panic(msg string, fields ...zap.Field) {
	zl.logger.Panic(msg, fields...)
}

func (zl *ZapLogger) Fatal(msg string, fields ...zap.Field) {
	zl.logger.Fatal(msg, fields...)
}
