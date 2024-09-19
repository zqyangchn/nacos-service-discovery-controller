package logger

import (
	"os"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var routerLogger *zap.Logger

func newGinDebugZapCore() zapcore.Core {
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), // 打印到控制台或文件
		zap.NewAtomicLevelAt(zapcore.DebugLevel),                // 日志级别
	)
}

func GetZapRouterLogger() *zap.Logger {
	return routerLogger
}

// GinDebugPrintRouteZapLoggerFunc print gin route function
func GinDebugPrintRouteZapLoggerFunc(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	routerLogger.Debug("",
		zap.Strings("Gin Debug", []string{httpMethod, absolutePath, handlerName, strconv.Itoa(nuHandlers)}),
	)
}
