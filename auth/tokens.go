package auth

import (
	"DistanceTrackerServer/constants"
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/ssh"
	"os"
	"time"
)

var (
	secretKey *rsa.PrivateKey
)

func init() {
	bytes, err := os.ReadFile(constants.JwtSecretkey)
	if err != nil {
		panic(err)
	}
	key, err := ssh.ParseRawPrivateKey(bytes)
	if err != nil {
		panic(err)
	}
	secretKey = key.(*rsa.PrivateKey)
}

func CreateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256,
		jwt.MapClaims{
			"email": email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey.Public(), nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}
	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}
	email, ok := token.Claims.(jwt.MapClaims)["email"].(string)
	if !ok {
		return "", fmt.Errorf("failed to extract email from token")
	}
	return email, nil
}
