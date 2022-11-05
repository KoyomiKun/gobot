package zaplogger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ZapLoggerOption func(*ZapLogger)

func WithEnv(env LogEnv) ZapLoggerOption {
	return func(zl *ZapLogger) {
		zl.config = map[LogEnv]zap.Config{
			Prod: zap.NewProductionConfig(),
			Dev:  zap.NewDevelopmentConfig(),
		}[env]
	}
}

func WithLevel(level LogLevel) ZapLoggerOption {
	return func(zl *ZapLogger) {
		zl.level = map[LogLevel]zapcore.Level{
			Debug:  zapcore.DebugLevel,
			Info:   zapcore.InfoLevel,
			Warn:   zapcore.WarnLevel,
			Error:  zapcore.ErrorLevel,
			Dpanic: zapcore.DPanicLevel,
			Panic:  zapcore.PanicLevel,
			Fatal:  zapcore.FatalLevel,
		}[level]
	}
}

func WithEncoder(msgKey, levelKey, timeKey, callerKey string) ZapLoggerOption {
	return func(zl *ZapLogger) {
		zl.encoder = zapcore.EncoderConfig{
			MessageKey:  msgKey,
			LevelKey:    levelKey,
			EncodeLevel: zapcore.CapitalLevelEncoder,
			TimeKey:     timeKey,
			CallerKey:   callerKey,
			EncodeTime: func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
				pae.AppendString(t.Format("2006-01-02 15:04:05"))
			},
			EncodeDuration: func(d time.Duration, pae zapcore.PrimitiveArrayEncoder) {
				pae.AppendInt64(int64(d) / 1e6)
			},
			EncodeCaller: zapcore.ShortCallerEncoder,
		}
	}
}

func WithLogFile(filepath string) ZapLoggerOption {
	return func(zl *ZapLogger) {
		zl.logFile = filepath
	}
}

func WithRotate(size, age, backups int, compress bool) ZapLoggerOption {
	return func(zl *ZapLogger) {
		zl.rotate = &lumberjack.Logger{
			MaxSize:    size,
			MaxBackups: backups,
			MaxAge:     age,
			Compress:   compress,
		}
	}
}
