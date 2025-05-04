package utils

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func ParseRequestBody(r *http.Request, v interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}

	err = r.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to close request body: %w", err)
	}

	err = json.Unmarshal(body, &v)
	if err != nil {
		return fmt.Errorf("Failed to parse request body: %w", err)
	}
	return nil
}

func WriteSimpleErrorResponse(w http.ResponseWriter, sugar *zap.SugaredLogger, statusCode int, message error) {
	sugar.Errorf("Error: %v", message)
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(message.Error()))
	if err != nil {
		sugar.Errorf("Failed to write response: %v", err)
	}
}

func WriteSimpleResponse(w http.ResponseWriter, sugar *zap.SugaredLogger, statusCode int, message string) {
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(message))
	if err != nil {
		sugar.Errorf("Failed to write response: %v", err)
	}
}
