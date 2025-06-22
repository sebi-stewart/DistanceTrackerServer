package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
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
