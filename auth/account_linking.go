package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func LinkAccounts(ctx *gin.Context, link AccountLink) error {
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
