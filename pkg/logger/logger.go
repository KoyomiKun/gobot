package logger

import (
	"context"

	uuid "github.com/satori/go.uuid"
)

const (
	RequestID = "RequestID"
)

var logger Logger

func SetGlobalLogger(loggerT Logger) {
	if logger != nil {
		return
	}

	logger = loggerT
}

type Logger interface {
	Debug(context.Context, string, ...interface{})
	Info(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
	Fatal(context.Context, string, ...interface{})
}

func Debug(ctx context.Context, format string, params ...interface{}) {
	logger.Debug(ctx, format, params...)
}

func Info(ctx context.Context, format string, params ...interface{}) {
	logger.Info(ctx, format, params...)
}
func Warn(ctx context.Context, format string, params ...interface{}) {
	logger.Warn(ctx, format, params...)
}
func Error(ctx context.Context, format string, params ...interface{}) {
	logger.Error(ctx, format, params...)
}
func Fatal(ctx context.Context, format string, params ...interface{}) {
	logger.Fatal(ctx, format, params...)
}

func GetUUID() string {
	return uuid.NewV4().String()
}

func CtxParams(ctx context.Context) []interface{} {
	if ctx == nil {
		return []interface{}{}
	}
	pms := make([]interface{}, 0)
	if requestID := ctx.Value(RequestID); requestID != nil {
		pms = append(pms, []interface{}{RequestID, requestID.(string)}...)
	}
	return pms
}
