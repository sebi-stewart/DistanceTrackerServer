package utils

import (
	"database/sql"
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
