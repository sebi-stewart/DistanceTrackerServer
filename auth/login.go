package auth

import (
	"DistanceTrackerServer/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	createToken = CreateToken
)

func Login(ctx *gin.Context, email string, password string) (int, error) {
	// First we get the password hash from the database
	// Then we compare the password hash with the password
	// If they match, we return the user ID
	// If they don't match, we return an error

	dbConn, err := utils.DBConnFromContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get database connection: %w", err)
	}

	var userID int
	var passwordHash string
	err = dbConn.QueryRow("SELECT id, password FROM users WHERE email = ?", email).Scan(&userID, &passwordHash)
	if err != nil {
		return 0, fmt.Errorf("failed to get user from database: %w", err)
	}

	// Check if the password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return 0, fmt.Errorf("password is incorrect: %w", err)
	}

	tokenString, err := createToken(email)
	if err != nil {
		return 0, fmt.Errorf("failed to create token: %w", err)
	}
	ctx.SetCookie("token", tokenString, 3600, "/", "", false, true)
	return userID, nil
}
