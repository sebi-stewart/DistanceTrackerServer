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

type EmailExtractor struct {
	Email string `json:"email"`
}

func AuthenticateRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		tokenString, err := ctx.Cookie("token")

		if path == "/login" || path == "/register" {
			if err == nil {
				email, verifyErr := verifyToken(tokenString)
				if verifyErr != nil {
					rejectRequest(ctx, http.StatusUnauthorized, "Invalid token")
					return
				}
				rejectRequest(ctx, http.StatusForbidden, "Already logged in, please logout first", email)
				return
			}
			ctx.Set("email", "NEW_USER")
			return
		}

		if err == nil {
			email, verifyErr := verifyToken(tokenString)
			if verifyErr != nil {
				rejectRequest(ctx, http.StatusUnauthorized, "Invalid token")
				return
			}
			ctx.Set("email", email)
			return
		}

		rejectRequest(ctx, http.StatusUnauthorized, "Missing credentials, please provide a valid token or login")
	}
}

func rejectRequest(ctx *gin.Context, statusCode int, message string, user ...string) {
	if len(user) > 0 {
		user = append(user, "anyonymous")
	}
	sugar, _ := sugarFromContext(ctx)
	ctx.JSON(statusCode, gin.H{"error": message})
	ctx.Abort()

	method := ctx.Request.Method
	endpoint := ctx.Request.URL.String()

	sugar.Infow("<-- Request Rejected -->",
		zap.String("user", user[0]),
		zap.String("method", method),
		zap.String("endpoint", endpoint),
		zap.Int("status_code", statusCode),
		zap.String("message", message),
	)
}
