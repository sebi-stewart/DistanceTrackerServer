package auth

import (
	"DistanceTrackerServer/models"
	"DistanceTrackerServer/utils"
	"github.com/google/uuid"

	"fmt"
	"github.com/gin-gonic/gin"
	//"github.com/google/uuid"
)

func LinkAccounts(ctx *gin.Context, link models.AccountLink) error {
	// This function should link the accounts in the database
	// For example, you might want to update a user record to include the new account ID
	//dbConn, err := utils.DBConnFromContext(ctx)
	//if err != nil {
	//	return err
	//}

	// Find the corresponding user IDs for the account link by email and password and then UUID
	// First login to the initiator account and get the user ID
	initiatorUserID, err := login(ctx, link.Email, link.Password)
	if err != nil {
		return err
	}

	fmt.Print(initiatorUserID)

	return nil

}

func CreateUuidLink(ctx *gin.Context, origin models.UserLogin) (models.AccountLink, error) {
	//First we login to the origin account to verify credentials
	initiatorUserID, err := login(ctx, origin.Email, origin.Password)
	if err != nil {
		return models.AccountLink{}, err
	}

	pairUUID := uuid.New()
	dbConn, err := utils.DBConnFromContext(ctx)
	if err != nil {
		return models.AccountLink{}, err
	}

	_, err = dbConn.Exec("DELETE FROM link_code WHERE user_id = ?", initiatorUserID)
	if err != nil {
		return models.AccountLink{}, fmt.Errorf("failed to delete existing link code: %w", err)
	}
	_, err = dbConn.Exec("INSERT INTO link_code (user_id, code) VALUES (?, ?)", initiatorUserID, pairUUID)
	if err != nil {
		return models.AccountLink{}, fmt.Errorf("failed to insert new link code: %w", err)
	}

	link := models.AccountLink{
		Email:    origin.Email,
		Password: origin.Password,
		PairUUID: pairUUID,
	}

	return link, nil
}
