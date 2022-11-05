package zaplogger

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/Koyomikun/gobot/pkg/logger"
)

type LogEnv int
type LogLevel int

const (
	Prod LogEnv = iota
	Dev
)

const (
	Debug LogLevel = iota
	Info
	Warn
	Error
	Dpanic
	Panic
	Fatal
)

type ZapLogger struct {
	config  zap.Config
	level   zapcore.Level
	encoder zapcore.EncoderConfig
	suger   *zap.SugaredLogger
	rotate  *lumberjack.Logger

	logFile string
}

func NewZapLogger(opts ...ZapLoggerOption) *ZapLogger {
	zl := &ZapLogger{
		config:  zap.NewDevelopmentConfig(),
		level:   zapcore.DebugLevel,
		encoder: zapcore.EncoderConfig{},
		logFile: "./log/app.log",
	}

	for _, opt := range opts {
		opt(zl)
	}

	zl.build()

	return zl
}

func (zl *ZapLogger) build() {
	zl.config.Level.SetLevel(zl.level)
	zl.config.EncoderConfig = zl.encoder

	logDir := filepath.Dir(zl.logFile)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.Mkdir(logDir, os.ModePerm)
	}

	if zl.rotate != nil {
		zl.buildWithRotate()
		return
	}

	zl.config.OutputPaths = []string{zl.logFile, "stdout"}
	logs, err := zl.config.Build()
	if err != nil {
		log.Fatalf("build zaplogger failed: %v", err)
	}

	zl.suger = logs.Sugar()
}

func (zl *ZapLogger) buildWithRotate() {

	zl.rotate.Filename = zl.logFile

	logs, err := zl.config.Build(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		w := zapcore.AddSync(zl.rotate)
		core := zapcore.NewCore(zapcore.NewJSONEncoder(zl.encoder), w, zl.level)
		return zapcore.NewTee(c, core)
	}))

	if err != nil {
		log.Fatalf("build zaplogger failed: %v", err)
	}

	zl.suger = logs.Sugar()

}

func (zl *ZapLogger) Debug(ctx context.Context, format string, kvs ...interface{}) {
	kvs = append(kvs, logger.CtxParams(ctx)...)
	zl.suger.Debugw(format, kvs...)
}

func (zl *ZapLogger) Info(ctx context.Context, format string, kvs ...interface{}) {
	kvs = append(kvs, logger.CtxParams(ctx)...)
	zl.suger.Infow(format, kvs...)
}
func (zl *ZapLogger) Warn(ctx context.Context, format string, kvs ...interface{}) {
	kvs = append(kvs, logger.CtxParams(ctx)...)
	zl.suger.Warnw(format, kvs...)
}
func (zl *ZapLogger) Error(ctx context.Context, format string, kvs ...interface{}) {

	kvs = append(kvs, logger.CtxParams(ctx)...)
	zl.suger.Errorw(format, kvs...)
}
func (zl *ZapLogger) Fatal(ctx context.Context, format string, kvs ...interface{}) {

	kvs = append(kvs, logger.CtxParams(ctx)...)
	zl.suger.Fatalw(format, kvs...)
}
