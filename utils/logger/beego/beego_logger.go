package beego

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/astaxie/beego/logs"

	ctxlog "github.com/Koyomikun/gobot/utils/logger"
)

const (
	logFuncCallDepth = 4
)

type Logger struct {
	beeLogger *logs.BeeLogger
}

func NewLogger(path string) ctxlog.Interface {
	var beeLogger *logs.BeeLogger
	beeLogger = logs.GetBeeLogger()
	beeLogger.EnableFuncCallDepth(true)
	beeLogger.SetLogFuncCallDepth(logFuncCallDepth)

	if path != "" {
		dirPath := filepath.Dir(path)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			fmt.Printf("mkdir %s failed: %v\n", dirPath, err)
		}

		if err := beeLogger.SetLogger(
			logs.AdapterFile,
			fmt.Sprintf(
				`{"filename": "%s", "level":7, "daily":true, "maxdays":30, "perm":"0666"}`,
				path,
			)); err != nil {
			panic(fmt.Sprintf("fail to set log: %v", err))
		}
		beeLogger.SetLogger(logs.AdapterConsole)
	}
	return &Logger{
		beeLogger: beeLogger,
	}
}

func (l *Logger) Infof(ctx context.Context, format string, v ...interface{}) {
	l.beeLogger.Info(ctxlog.CtxSprintf(ctx, format, v...))
}

func (l *Logger) Warnf(ctx context.Context, format string, v ...interface{}) {
	l.beeLogger.Warning(ctxlog.CtxSprintf(ctx, format, v...))
}

func (l *Logger) Errorf(ctx context.Context, format string, v ...interface{}) {
	l.beeLogger.Error(ctxlog.CtxSprintf(ctx, format, v...))
}

func (l *Logger) Fatalf(ctx context.Context, format string, v ...interface{}) {
	l.beeLogger.Critical(ctxlog.CtxSprintf(ctx, format, v...))
	os.Exit(1)
}
