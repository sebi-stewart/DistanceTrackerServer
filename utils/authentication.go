package utils

import (
	"DistanceTrackerServer/constants"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func EmailFromContext(ctx *gin.Context) (string, error) {
	email, ok := ctx.Get("email")
	if !ok {
		return "", fmt.Errorf("failed to retrieve email from context")
	}

	emailStr, ok := email.(string)
	if !ok {
		return "", fmt.Errorf("email in context is not a string")
	}

	return emailStr, nil
}

func GetUserIdByEmail(dbConn *sql.DB, email string) (int, error) {
	var userID int
	err := dbConn.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("user with email %s not found", email)
		}
		return 0, fmt.Errorf("failed to query user ID by email: %w", err)
	}
	return userID, nil
}

func LogRejectedRequest(ctx *gin.Context, sugar *zap.SugaredLogger, statusCode int, reason string, user string) error {
	sugar.Infow("Logging Rejected Request", user, reason)
	dbConn, err := DBConnFromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve database connection: %w", err)
	}

	clientIp := ctx.ClientIP()
	loggingQuery := `INSERT INTO rejected_requests(user_email, status_code, reason, ip_address) VALUES (?, ?, ?, ?)`
	_, err = dbConn.Exec(loggingQuery, user, statusCode, reason, clientIp)
	if err != nil {
		return fmt.Errorf("failed to log rejected request: %w", err)
	}

	// View how often the current ip address has been rejected within the last 24 hours and ban if necessary
	var count int

	countRejectionQuery := `SELECT COUNT(*) FROM rejected_requests WHERE ip_address = ? AND created_at >= datetime('now', '-24 hours')`
	err = dbConn.QueryRow(countRejectionQuery, clientIp).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to count rejected requests for IP %s: %w", clientIp, err)
	}

	if count >= constants.RequestsUntilBan {
		err = BanRequestIp(dbConn, sugar, clientIp, count)
		if err != nil {
			return fmt.Errorf("failed to ban IP %s after too many rejected requests: %w", clientIp, err)
		}
	}

	sugar.Info("Successfully Logged Rejected request")
	return nil
}

func BanRequestIp(dbConn *sql.DB, sugar *zap.SugaredLogger, clientIp string, count int) error {
	sugar.Infow("Attempting to ban IP", "ip", clientIp, "failed_requests", count)

	// Check if the IP has already been banned before and if so how many times
	var alreadyBanned bool
	checkBanQuery := `SELECT EXISTS(SELECT 1 FROM banned_ips WHERE ip_address = ?)`
	err := dbConn.QueryRow(checkBanQuery, clientIp).Scan(&alreadyBanned)
	if err != nil {
		return fmt.Errorf("failed to check if IP %s is already banned: %w", clientIp, err)
	}

	// Create a new ban entry if the IP is not already banned
	if !alreadyBanned {
		banQuery := `INSERT INTO banned_ips (ip_address, reason, banned_length, banned_until) VALUES (?, ?, 0.25, datetime('now', '+15 minutes'))` // 0.25 hours = 15 minutes
		_, err = dbConn.Exec(banQuery, clientIp, fmt.Sprintf("Too many rejected requests - %d failed requests", count))
		if err != nil {
			return fmt.Errorf("failed to ban IP %s: %w", clientIp, err)
		}
		sugar.Warnf("Initially Banned IP %s due to too many rejected requests - 15 minute ban", clientIp)
		return nil
	}

	// If the IP is already banned, update the ban length and increment the ban count
	var newBanLength float64
	updateBanQuery := `UPDATE banned_ips SET 
                      banned_length = banned_length * 2, 
                      banned_times = banned_times + 1, 
                      last_banned_at = datetime('now'),
                      banned_until = datetime('now', '+' || (banned_length * 2) || ' hours'),
                      reason = reason || ' | Banned again due to too many rejected requests - ' || ? || ' failed requests'
                  WHERE ip_address = ?
                  RETURNING banned_length`
	err = dbConn.QueryRow(updateBanQuery, count, clientIp).Scan(&newBanLength)
	if err != nil {
		return fmt.Errorf("failed to update ban for IP %s: %w", clientIp, err)
	}
	sugar.Warnf("Updated ban for IP %s due to too many rejected requests - new ban length: %f hours", clientIp, newBanLength)
	return nil
}

func IsIpBanned(ctx *gin.Context, ip string) (bool, sql.NullTime, error) {
	var bannedUntil sql.NullTime

	dbConn, err := DBConnFromContext(ctx)
	if err != nil {
		return false, bannedUntil, fmt.Errorf("failed to retrieve database connection: %w", err)
	}

	query := `SELECT banned_until FROM banned_ips WHERE ip_address = ? AND banned_until > datetime('now') LIMIT 1`
	queryErr := dbConn.QueryRow(query, ip).Scan(&bannedUntil)
	if queryErr != nil {
		if errors.Is(queryErr, sql.ErrNoRows) {
			return false, bannedUntil, nil // IP is not banned
		}
		return false, bannedUntil, fmt.Errorf("failed to check if IP %s is banned: %w", ip, queryErr)
	}

	if !bannedUntil.Valid {
		return false, bannedUntil, fmt.Errorf("banned_until field is not valid for IP %s", ip)
	}
	return true, bannedUntil, nil // IP is banned until the specified time
}
