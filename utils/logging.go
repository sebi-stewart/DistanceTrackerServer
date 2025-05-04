package utils

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var (
	saveLoggersFunc = SaveLoggersToContext
)

func LoggerFromContext(ctx context.Context) (*zap.Logger, error) {
	logger, ok := ctx.Value("logger").(*zap.Logger)
	if ok && logger != nil {
		return logger, nil
	}
	return nil, fmt.Errorf("failed to retrieve logger from context")
}

func SugarFromContext(ctx context.Context) (*zap.SugaredLogger, error) {
	sugar, ok := ctx.Value("sugar").(*zap.SugaredLogger)
	if ok && sugar != nil {
		return sugar, nil
	}
	return nil, fmt.Errorf("failed to retrieve sugar from context")
}

func SaveLoggersToContext(ctx context.Context, logger *zap.Logger, sugar *zap.SugaredLogger) context.Context {
	ctx = context.WithValue(ctx, "logger", logger)
	ctx = context.WithValue(ctx, "sugar", sugar)
	return ctx
}

func CreateAndSaveLoggers(ctx context.Context) (context.Context, error) {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(zapcore.Lock(os.Stdout)), zapcore.DebugLevel))

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.PanicLevel))
	sugar := logger.Sugar()
	ctx = saveLoggersFunc(ctx, logger, sugar)
	return ctx, nil
}
