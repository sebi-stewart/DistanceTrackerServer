package auth

import (
	"DistanceTrackerServer/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

var (
	verifyToken      = VerifyToken
	sugarFromContext = utils.SugarFromContext
)

func AuthenticateRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Escape early if the request is the health check, login or register endpoint
		if ctx.Request.URL.Path == "/login" || ctx.Request.URL.Path == "/register" {
			return
		}
		tokenString, err := ctx.Cookie("token")
		if err == nil {
			email, verificationError := verifyToken(tokenString)
			if verificationError != nil {
				rejectRequest(ctx, http.StatusUnauthorized, "Invalid token")
				return
			}
			ctx.Set("email", email)
			return
		}
		rejectRequest(ctx, http.StatusUnauthorized, "Missing credentials, please provide a valid token or login")
		return
	}
}

func rejectRequest(ctx *gin.Context, statusCode int, message string) {
	sugar, _ := sugarFromContext(ctx)
	ctx.JSON(statusCode, gin.H{"error": message})
	ctx.Abort()

	method := ctx.Request.Method
	endpoint := ctx.Request.URL.String()

	sugar.Infow("<-- Request Rejected -->",
		zap.String("user", "anyonymous"),
		zap.String("method", method),
		zap.String("endpoint", endpoint),
		zap.Int("status_code", statusCode),
		zap.String("message", message),
	)
}
