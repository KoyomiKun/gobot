package main

import (
	"context"

	"github.com/Koyomikun/gobot/pkg/logger"
	"github.com/Koyomikun/gobot/pkg/logger/zaplogger"
)

func main() {
	logger.SetGlobalLogger(zaplogger.NewZapLogger(
		zaplogger.WithEnv(zaplogger.Dev),
		zaplogger.WithLevel(zaplogger.Debug),
		zaplogger.WithEncoder("msg", "level", "ts", "file"),
		zaplogger.WithRotate(20, 10, 30, false), // 20MB, 30Days
	))

	ctx := context.WithValue(context.Background(), logger.RequestID, logger.GetUUID())
	logger.Debug(ctx, "ahh", "test", "123")
}
