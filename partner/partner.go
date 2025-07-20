package partner

import (
	"DistanceTrackerServer/models"
	"DistanceTrackerServer/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// DistanceHandler /**
// This handler is responsible for processing distance-related requests.
//
//	We will only have one endpoint, as we will use the endpoint to submit the current location of the user and return the
//	distance the partner is away.
//
// */
func DistanceHandler(ctx *gin.Context) {
	sugar, err := utils.SugarFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL SERVER ERROR"})
		return
	}

	dbConn, err := dbConnFromContext(ctx)
	if err != nil {
		sugar.Errorw("Error retrieving database connection from context", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL SERVER ERROR"})
		return
	}

	location := models.Location{}
	err = ctx.BindJSON(&location)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userEmail, err := emailFromContext(ctx)
	if err != nil {
		sugar.Errorw("Error retrieving email from context", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL SERVER ERROR"})
		return
	}

	userId, err := utils.GetUserIdByEmail(dbConn, userEmail)
	if err != nil {
		sugar.Errorw("Error retrieving user ID by email", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL SERVER ERROR"})
		return
	}

	validationErr := validateDistanceRequest(location, dbConn, userId)
	if validationErr != nil {
		sugar.Errorw("Distance validation failed, inserting location into database as invalid", "error", validationErr)
	} else {
		sugar.Info("Successfully validated Location Request, saving to db and returning partner location to user")
	}

	err = insertLocationToDB(location, dbConn, userId, validationErr == nil)
	if err != nil {
		sugar.Errorw("Error inserting location into database", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL SERVER ERROR"})
		return
	}
	if validationErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	partnerLocation, err := retrievePartnerLocation(dbConn, userId)
	if err != nil {
		sugar.Errorw("Error retrieving partner location", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL SERVER ERROR"})
		return
	}

	distance := calculateDistance(location, partnerLocation.ToLocation())
	sugar.Infow("Successfully calculated distance", "distance", distance)
	ctx.JSON(http.StatusOK, gin.H{
		"distance": distance,
	})
}

func InformationHandler(ctx *gin.Context) {
	sugar, err := utils.SugarFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL SERVER ERROR"})
		return
	}

	dbConn, err := dbConnFromContext(ctx)
	if err != nil {
		sugar.Errorw("Error retrieving database connection from context", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL SERVER ERROR"})
		return
	}

	userEmail, err := emailFromContext(ctx)
	if err != nil {
		sugar.Errorw("Error retrieving email from context", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL SERVER ERROR"})
		return
	}

	info, retrievalErr := Information(dbConn, userEmail)
	if retrievalErr != nil {
		sugar.Errorw("Error retrieving partner information", "error", retrievalErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL SERVER ERROR"})
		return
	}
	sugar.Infow("Successfully retrieved partner information", "info", info)
	ctx.JSON(http.StatusOK, info)
}
