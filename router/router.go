package router

import (
	"DistanceTrackerServer/auth"
	"DistanceTrackerServer/constants"
	"DistanceTrackerServer/database"
	"DistanceTrackerServer/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var (
	log                 *zap.Logger
	logRequest          = LogRequest
	addRouterMiddleware = AddRouterMiddleware
	authenticateRequest = auth.AuthenticateRequest
	sugarFromContext    = utils.SugarFromContext
	register            = auth.RegisterHandler
	login               = auth.LoginHandler
	logout              = auth.LogoutHandler
	accountLinkCreation = auth.AccountLinkCreationHandler
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

func AddRouterMiddleware(dbConn *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := uuid.New()
		logger := log.With(zap.String("request_id", id.String()))
		ctx.Set("logger", logger)
		ctx.Set("sugar", logger.Sugar())
		ctx.Set("request_id", id.String())
		ctx.Set("dbConn", dbConn)
		ctx.Next()
	}
}

func HealthCheckHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "All good here :D")
}

func Init(logger *zap.Logger) *gin.Engine {
	log = logger
	sugar := log.Sugar()

	sugar.Info("intializing sql connection")
	db, err := sql.Open("sqlite3", constants.DatabaseFile)
	if err != nil {
		sugar.Fatal("Failed to open database: ", err)
	}
	sugar.Info("Initializing database")
	err = database.InitDatabase(db)
	if err != nil {
		sugar.Fatal("Failed to initialize database: ", err)
	}

	sugar.Info("Initializing router")
	router := gin.New()
	err = router.SetTrustedProxies(nil)
	if err != nil {
		sugar.Fatal("Failed to set trusted proxies: ", err)
	}
	router.Use(addRouterMiddleware(db))
	router.Use(logRequest())
	router.Use(authenticateRequest())

	sugar.Info("Registering routes")
	router.GET("/healthcheck", healthCheckHandler)
	router.POST("/register", register)
	router.POST("/login", login)
	router.DELETE("/logout", logout)
	router.POST("/account-link-creation", accountLinkCreation)

	return router
}
