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

func CheckIfIpIsBanned() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		sugar, _ := sugarFromContext(ctx)

		isBanned, bannedUntil, err := utils.IsIpBanned(ctx, ip)
		if err != nil {
			sugar.Errorw("Error checking if IP is banned",
				zap.String("ip", ip),
				zap.Error(err),
			)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			ctx.Abort()
			return
		}

		if isBanned {
			if bannedUntil.Valid {
				bannedUntilTime := bannedUntil.Time
				sugar.Infow("--> Banned IP detected <--",
					zap.String("ip", ip),
					zap.Time("banned_until", bannedUntilTime))
				ctx.JSON(http.StatusForbidden, gin.H{"error": "Your IP is banned until " + bannedUntilTime.Format("2006-01-02 15:04:05") + " UTC"})
			} else {
				sugar.Infow("--> Banned IP detected with no specific end time <--",
					zap.String("ip", ip))
				ctx.JSON(http.StatusForbidden, gin.H{"error": "Your IP is banned"})
			}
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func AuthenticateRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		tokenString, err := ctx.Cookie("token")

		if err != nil { // no token provdided
			if path == "/login" || path == "/register" {
				ctx.Set("email", "NEW_USER")
				return
			}
			rejectRequest(ctx, http.StatusUnauthorized, "Missing credentials, please provide a valid token or login")
			return
		}

		email, verifyErr := verifyToken(tokenString)
		if verifyErr != nil {
			rejectRequest(ctx, http.StatusUnauthorized, "Invalid token")
			return
		}

		ctx.Set("email", email)

		if path == "/login" || path == "/register" {
			rejectRequest(ctx, http.StatusForbidden, "Already logged in, please logout first")
			return
		}

	}
}

func rejectRequest(ctx *gin.Context, statusCode int, reason string) {
	userEmail, err := utils.EmailFromContext(ctx)

	user := "unknown"
	if err == nil && userEmail != "" {
		user = userEmail
	}

	sugar, _ := sugarFromContext(ctx)

	loggingErr := utils.LogRejectedRequest(ctx, sugar, statusCode, reason, user)
	if loggingErr != nil {
		sugar.Error("Error logging rejected request", loggingErr)
	}

	ctx.JSON(statusCode, gin.H{"error": reason})
	ctx.Abort()

	method := ctx.Request.Method
	endpoint := ctx.Request.URL.String()

	sugar.Infow("<-- Request Rejected -->",
		zap.String("user", user),
		zap.String("method", method),
		zap.String("endpoint", endpoint),
		zap.Int("status_code", statusCode),
		zap.String("reason", reason),
	)
}
