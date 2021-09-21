package logger

import "context"

type Interface interface {
	Infof(context.Context, string, ...interface{})
	Warnf(context.Context, string, ...interface{})
	Errorf(context.Context, string, ...interface{})
	Fatalf(context.Context, string, ...interface{})
}

var logger Interface

func SetLogger(loggerT Interface) {
	if loggerT != nil {
		logger = loggerT
	}
}

func Infof(ctx context.Context, format string, v ...interface{}) {
	logger.Infof(ctx, format, v...)
}
func Warnf(ctx context.Context, format string, v ...interface{}) {
	logger.Warnf(ctx, format, v...)
}
func Errorf(ctx context.Context, format string, v ...interface{}) {
	logger.Errorf(ctx, format, v...)
}
func Fatalf(ctx context.Context, format string, v ...interface{}) {
	logger.Fatalf(ctx, format, v...)
}
