package auth

import (
	"DistanceTrackerServer/utils"
	"crypto/rand"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	createToken        = CreateToken
	randomHash  []byte = generateRandomHash()
)

func randomString(length int) string {
	b := make([]byte, length+2)
	_, err := rand.Read(b)
	if err != nil {
		panic(fmt.Errorf("failed to generate random string: %w", err))
	}
	return fmt.Sprintf("%x", b)[2 : length+2]
}

func generateRandomHash() []byte {
	// This function generates a random hash for password comparison
	// It is used to ensure that we do not leak information based on the runtime of the request
	randomPassword := randomString(10)
	hash, err := bcrypt.GenerateFromPassword([]byte(randomPassword), bcrypt.DefaultCost)
	if err != nil {
		panic(fmt.Errorf("failed to generate random hash: %w", err))
	}
	return hash
}

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
		// Hash and compare password to a random value to ensure we don't leak information based on the runtime of the request
		_ = bcrypt.CompareHashAndPassword(randomHash, []byte(password))

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
