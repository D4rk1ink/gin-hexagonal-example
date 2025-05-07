package logger

import (
	"context"

	"go.uber.org/zap"
)

type logger struct {
	z   *zap.Logger
	ctx *context.Context
}

var Logger *logger

func Init() error {
	z, err := zap.NewProduction()
	if err != nil {
		return err
	}
	defer z.Sync()

	Logger = &logger{
		z: z,
	}

	return nil
}

func WithContextBackground(ctx context.Context) *logger {
	z, _ := zap.NewProduction()
	return &logger{
		z:   z,
		ctx: &ctx,
	}
}

func Info(msg string, fields ...zap.Field) {
	Logger.z.Info(msg, fields...)
}
func Error(msg string, fields ...zap.Field) {
	Logger.z.Error(msg, fields...)
}
func Debug(msg string, fields ...zap.Field) {
	Logger.z.Debug(msg, fields...)
}
func Warn(msg string, fields ...zap.Field) {
	Logger.z.Warn(msg, fields...)
}
func Fatal(msg string, fields ...zap.Field) {
	Logger.z.Fatal(msg, fields...)
}
func Panic(msg string, fields ...zap.Field) {
	Logger.z.Panic(msg, fields...)
}
