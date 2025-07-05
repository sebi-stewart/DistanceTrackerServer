package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

func DBConnFromContext(ctx *gin.Context) (*sql.DB, error) {
	dbConn, ok := ctx.Value("dbConn").(*sql.DB)
	if ok && dbConn != nil {
		return dbConn, nil
	}
	return nil, fmt.Errorf("failed to retrieve database connection from context")
}

func GetPartnerIdByUserId(dbConn *sql.DB, userId int) (int, error) {
	var partnerId int
	err := dbConn.QueryRow("SELECT linked_account FROM users WHERE id = ?", userId).Scan(&partnerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("no linked account found for user ID %d", userId)
		}
		return 0, fmt.Errorf("failed to query partner ID by user ID: %w", err)
	}
	return partnerId, nil
}
