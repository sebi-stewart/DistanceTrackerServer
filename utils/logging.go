package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var (
	Logger        *zap.Logger
	Sugar         *zap.SugaredLogger
	createLoggers = CreateLoggers
)

func init() {
	createLoggers()
}

func CreateLoggers() {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(zapcore.Lock(os.Stdout)), zapcore.DebugLevel))

	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.PanicLevel))
	Sugar = Logger.Sugar()

	Logger.Info("Logger initialized")
	Sugar.Info("Sugared logger initialized")
}

func LoggerFromContext(ctx *gin.Context) (*zap.Logger, error) {
	logger, ok := ctx.Value("logger").(*zap.Logger)
	if ok && logger != nil {
		return logger, nil
	}
	return nil, fmt.Errorf("failed to retrieve logger from context")
}

func SugarFromContext(ctx *gin.Context) (*zap.SugaredLogger, error) {
	sugar, ok := ctx.Value("sugar").(*zap.SugaredLogger)
	if ok && sugar != nil {
		return sugar, nil
	}
	return nil, fmt.Errorf("failed to retrieve sugar from context")
}
