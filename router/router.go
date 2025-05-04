package router

import (
	"DistanceTrackerServer/auth"
	"DistanceTrackerServer/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var (
	log                 *zap.Logger
	logRequest          = LogRequest
	addRouterMiddleware = AddRouterMiddleware
	sugarFromContext    = utils.SugarFromContext
	register            = auth.RegisterHandler
	login               = auth.LoginHandler
	logout              = auth.LogoutHandler
	healthCheckHandler  = HealthCheckHandler
)

func LogRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		method := ctx.Request.Method
		endpoint := ctx.Request.URL.String()
		sugar, _ := sugarFromContext(ctx)
		user, _ := ctx.Get("user")

		sugar.Infow("-->",
			zap.String("user", fmt.Sprintf("%v", user)),
			zap.String("method", method),
			zap.String("endpoint", endpoint),
		)

		start := time.Now()
		ctx.Next()

		duration := time.Since(start)
		statusCode := ctx.Writer.Status()
		sugar.Infow("<--",
			zap.String("method", method),
			zap.String("endpoint", endpoint),
			zap.Int("status_code", statusCode),
			zap.Duration("duration", duration),
		)
	}
}

func AddRouterMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := uuid.New()
		logger := log.With(zap.String("request_id", id.String()))
		ctx.Set("logger", logger)
		ctx.Set("sugar", logger.Sugar())
		ctx.Set("request_id", id.String())
		ctx.Next()
	}
}

func HealthCheckHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "All good here :D")
}

func Init(logger *zap.Logger) *gin.Engine {
	log = logger
	sugar := log.Sugar()

	router := gin.New()
	err := router.SetTrustedProxies(nil)
	if err != nil {
		sugar.Fatal("Failed to set trusted proxies: ", err)
	}
	router.Use(addRouterMiddleware())
	router.Use(logRequest())

	router.GET("/healthcheck", healthCheckHandler)
	router.POST("/register", register)
	router.POST("/login", login)
	router.DELETE("/logout", logout)

	return router
}
