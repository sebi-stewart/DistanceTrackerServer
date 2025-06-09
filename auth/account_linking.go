package auth

import (
	"DistanceTrackerServer/models"
	"DistanceTrackerServer/utils"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"time"

	"fmt"
	"github.com/gin-gonic/gin"
	//"github.com/google/uuid"
)

func LinkAccounts(ctx *gin.Context, link models.AccountLink) error {
	// First login to the initiator account and get the user ID
	initiatorUserID, err := login(ctx, link.Email, link.Password)
	if err != nil {
		return err
	}

	dbConn, err := utils.DBConnFromContext(ctx)
	if err != nil {
		return err
	}

	// Now we find the account to link to
	var linkUserID int
	var creationTime time.Time
	err = dbConn.QueryRow("SELECT user_id, created_at FROM link_code WHERE code = ?", link.PairUUID).Scan(&linkUserID, &creationTime)
	if err != nil {
		return fmt.Errorf("failed to find link code: %w", err)
	}

	// Check if the link code is expired (valid for 15 minutes)
	if time.Since(creationTime) > 15*time.Minute {
		// If the link code is expired, we delete it from the database
		_, err = dbConn.Exec("DELETE FROM link_code WHERE code = ?", link.PairUUID)
		return fmt.Errorf("link code expired")
	}

	// First we remove any links either account has to any other account
	_, err = dbConn.Exec("DELETE FROM link_code WHERE user_id = ? OR user_id = ?", initiatorUserID, linkUserID)
	if err != nil {
		return fmt.Errorf("failed to remove existing links: %w", err)
	}

	// Next we link the accounts
	_, err1 := dbConn.Exec("UPDATE users SET linked_account = ? WHERE id = ?", initiatorUserID, linkUserID)
	_, err2 := dbConn.Exec("UPDATE users SET linked_account = ? WHERE id = ?", linkUserID, initiatorUserID)
	if err1 != nil || err2 != nil {
		return fmt.Errorf("failed to link accounts: %w | %w", err1, err2)
	}

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

	// Delete any existing link code for the initiator user
	_, err = dbConn.Exec("DELETE FROM link_code WHERE user_id = ?", initiatorUserID)
	if err != nil {
		return models.AccountLink{}, fmt.Errorf("failed to delete existing link code: %w", err)
	}

	// Remove any existing linked accounts for the initiator user, or the potential link user
	// This ensures that the initiator user can only have one link code at a time
	var existingLinkUserID int
	err = dbConn.QueryRow("SELECT linked_account FROM users WHERE id = ?", initiatorUserID).Scan(&existingLinkUserID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.AccountLink{}, fmt.Errorf("failed to get existing linked account: %w", err)
	}

	if existingLinkUserID != 0 {
		_, err = dbConn.Exec("UPDATE users SET linked_account = NULL WHERE id = ?", existingLinkUserID)
		if err != nil {
			return models.AccountLink{}, fmt.Errorf("failed to remove existing linked account: %w", err)
		}
	}
	_, err = dbConn.Exec("UPDATE users SET linked_account = NULL WHERE id = ?", initiatorUserID)
	if err != nil {
		return models.AccountLink{}, fmt.Errorf("failed to remove existing linked account: %w", err)
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
