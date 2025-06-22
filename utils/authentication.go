package utils

import (
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
