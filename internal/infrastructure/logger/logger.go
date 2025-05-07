package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	l   *zap.Logger
	ctx *context.Context
}

var Logger *logger

func Init() error {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	l, err := config.Build()
	if err != nil {
		return err
	}
	defer l.Sync()

	Logger = &logger{
		l: l,
	}

	return nil
}

// TODO
func WithContextBackground(ctx context.Context) *logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	l, _ := config.Build()

	defer l.Sync()

	return &logger{
		l:   l,
		ctx: &ctx,
	}
}

func Info(msg string, fields ...zap.Field) {
	Logger.l.Info(msg, fields...)
}
func Error(msg string, fields ...zap.Field) {
	Logger.l.Error(msg, fields...)
}
func Debug(msg string, fields ...zap.Field) {
	Logger.l.Debug(msg, fields...)
}
func Warn(msg string, fields ...zap.Field) {
	Logger.l.Warn(msg, fields...)
}
func Fatal(msg string, fields ...zap.Field) {
	Logger.l.Fatal(msg, fields...)
}
func Panic(msg string, fields ...zap.Field) {
	Logger.l.Panic(msg, fields...)
}
