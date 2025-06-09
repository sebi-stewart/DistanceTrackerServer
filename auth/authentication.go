package auth

import (
	"DistanceTrackerServer/models"
	"DistanceTrackerServer/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

var (
	validateRegistration = ValidateRegistration
	register             = Register
	linkAccounts         = LinkAccounts
	login                = Login
)

func RegisterHandler(ctx *gin.Context) {
	sugar, err := utils.SugarFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL SERVER ERROR"})
		return
	}

	newUser := models.UserRegister{}
	err = ctx.BindJSON(&newUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErr := validateRegistration(newUser)
	if validationErr != nil {
		sugar.Errorw("Error", validationErr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	userID, err := register(ctx, newUser)
	if err != nil {
		sugar.Errorw("registration error",
			zap.String("Error", err.Error()),
			zap.String("User", newUser.ToString()),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user registered successfully", "user_id": userID})
}

func LoginHandler(ctx *gin.Context) {
	sugar, err := utils.SugarFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL SERVER ERROR"})
		return
	}

	loginData := models.UserLogin{}
	err = ctx.BindJSON(&loginData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := login(ctx, loginData.Email, loginData.Password)
	if err != nil {
		sugar.Errorw("Error", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "successfully logged in", "user_id": userID})
}

func LogoutHandler(ctx *gin.Context) {
	sugar, err := utils.SugarFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL SERVER ERROR"})
		return
	}

	sugar.Info("LOGGED OUT")

	ctx.JSON(http.StatusOK, gin.H{"message": "LOGGED OUT"})
}

func AccountLinkHandler(ctx *gin.Context) {
	sugar, err := utils.SugarFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL SERVER ERROR"})
		return
	}

	accountLink := models.AccountLink{}
	err = ctx.BindJSON(&accountLink)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	linkingErr := linkAccounts(ctx, accountLink)
	if linkingErr != nil {
		sugar.Errorw("linking error",
			zap.String("Error", linkingErr.Error()),
			zap.String("AccountLink", accountLink.ToString()),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": linkingErr.Error()})
		return
	}

	sugar.Info("ACCOUNTS LINKED")

	ctx.JSON(http.StatusOK, gin.H{"message": "ACCOUNT LINKED"})
}
