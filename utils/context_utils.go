package utils

import (
	"context"
	"fmt"
	"go.uber.org/zap"
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
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:      false,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stdout", "logfile"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build()
	if err != nil {
		return ctx, err
	}

	sugar := logger.Sugar()
	ctx = saveLoggersFunc(ctx, logger, sugar)
	return ctx, nil
}
