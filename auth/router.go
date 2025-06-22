package auth

import (
	"DistanceTrackerServer/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	loginFunc        = Login
	verifyToken      = VerifyToken
	sugarFromContext = utils.SugarFromContext
)

func AuthenticateRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Escape early if the request is the health check, login or register endpoint
		if ctx.Request.URL.Path == "/healthcheck" || ctx.Request.URL.Path == "/login" || ctx.Request.URL.Path == "/register" {
			return
		}
		tokenString, err := ctx.Cookie("token")
		if err == nil {
			email, verificationError := verifyToken(tokenString)
			if verificationError != nil {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Credentials"})
				ctx.Abort()
			}
			ctx.Set("email", email)
			return
		}

		_, _, ok := ctx.Request.BasicAuth()
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing credentials, please provide a valid token or login"})
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Basic authentication is not supported, please use token-based authentication"})
		ctx.Abort()
		return

		//sugar, _ := sugarFromContext(ctx)
		//sugar.Infow("Authenticating user", zap.String("user", user))
		//sugar.Infow("Authenticating password", zap.String("password", password))
		//
		//_, err = loginFunc(ctx, user, password)
		//if err != nil {
		//	sugar.Errorw("Authentication error", zap.String("user", user), zap.Error(err))
		//	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Credentials"})
		//	ctx.Abort()
		//	return
		//}
	}
}
