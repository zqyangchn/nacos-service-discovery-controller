package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var once sync.Once

func Setup() {
	once.Do(setup)
}

func setup() {
	logger := zap.New(newAllFieldsZapCore(), zap.AddCaller(), zap.AddCallerSkip(1))
	defer logger.Sync()
	zap.ReplaceGlobals(logger)
	logger.Info("已初始化 全局 logger .")

	routerLogger = zap.New(newGinDebugZapCore())
	defer routerLogger.Sync()
	logger.Info("已初始化 gin logger .")
}

func newAllFieldsZapCore() zapcore.Core {
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), // 打印到控制台或文件
		zap.NewAtomicLevelAt(zapcore.InfoLevel),                 // 日志级别
	)
}

func Debug(msg string, fields ...zap.Field) {
	zap.L().Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	zap.L().Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	zap.L().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	zap.L().Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	zap.L().Fatal(msg, fields...)
}
