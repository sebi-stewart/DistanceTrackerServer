package auth

import (
	"DistanceTrackerServer/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

var (
	loginFunc        = Login
	verifyToken      = VerifyToken
	sugarFromContext = utils.SugarFromContext
)

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
