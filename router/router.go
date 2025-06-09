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
	authenticateRequest = AuthenticateRequest
	sugarFromContext    = utils.SugarFromContext
	register            = auth.RegisterHandler
	login               = auth.LoginHandler
	loginFunc           = auth.Login
	logout              = auth.LogoutHandler
	verifyToken         = auth.VerifyToken
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

func AuthenticateRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := ctx.Cookie("token")
		if err == nil {
			verificationError := verifyToken(tokenString)
			if verificationError != nil {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Credentials"})
				ctx.Abort()
			}
			return
		}

		user, password, ok := ctx.Request.BasicAuth()
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing credentials"})
			ctx.Abort()
			return
		}

		sugar, _ := sugarFromContext(ctx)
		sugar.Infow("Authenticating user", zap.String("user", user))
		sugar.Infow("Authenticating password", zap.String("password", password))

		_, err = loginFunc(ctx, user, password)
		if err != nil {
			sugar.Errorw("Authentication error", zap.String("user", user), zap.Error(err))
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Credentials"})
			ctx.Abort()
			return
		}
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

	return router
}
