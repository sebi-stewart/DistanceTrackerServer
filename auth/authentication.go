package auth

import (
	"DistanceTrackerServer/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	validateRegistration = ValidateRegistration
)

type UserRegister struct {
	Email           string `json:"email"`
	FirstName       string `json:"first_name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func RegisterHandler(ctx *gin.Context) {
	sugar, err := utils.SugarFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL SERVER ERROR"})
		return
	}

	newUser := UserRegister{}
	err = ctx.BindJSON(&newUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	validationErr := validateRegistration(newUser)
	if validationErr != nil {
		sugar.Errorw("Error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func LoginHandler(ctx *gin.Context) {
	sugar, err := utils.SugarFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL SERVER ERROR"})
		return
	}

	sugar.Info("LOGGED IN")

	ctx.JSON(http.StatusOK, gin.H{"message": "LOGGED IN"})
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
