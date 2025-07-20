package partner

import (
	"DistanceTrackerServer/models"
	"DistanceTrackerServer/utils"
	"database/sql"
	"fmt"
)

func Information(dbConn *sql.DB, userEmail string) (models.UserInformation, error) {
	userId, err := utils.GetUserIdByEmail(dbConn, userEmail)
	if err != nil {
		return models.UserInformation{}, fmt.Errorf("failed to retrieve user ID: %w", err)
	}

	partnerId, err := utils.GetPartnerIdByUserId(dbConn, userId)
	if err != nil {
		return models.UserInformation{}, fmt.Errorf("failed to retrieve partner ID: %w", err)
	}

	return getUserInformation(dbConn, partnerId)
}

func getUserInformation(dbConn *sql.DB, userId int) (models.UserInformation, error) {
	var userInfo models.UserInformation
	err := dbConn.QueryRow("SELECT email, name FROM users WHERE id = ?", userId).Scan(&userInfo.Email, &userInfo.FirstName)
	if err != nil {
		return models.UserInformation{}, fmt.Errorf("failed to retrieve user information: %w", err)
	}
	return userInfo, nil
}
