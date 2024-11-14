package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func saveTokenToFile(token string) error {
	filePath := "tokens.json"
	tokenData := map[string]string{
		"token": token,
	}
	fileContent, err := json.MarshalIndent(tokenData, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling token data: %w", err)
	}

	err = os.WriteFile(filePath, fileContent, 0644)
	if err != nil {
		return fmt.Errorf("error writing token to file: %w", err)
	}

	return nil
}

func readTokenFromFile() (string, error) {
	filePath := "tokens.json"
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading token from file: %w", err)
	}

	var tokenData map[string]string
	err = json.Unmarshal(data, &tokenData)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling token data: %w", err)
	}

	token, exists := tokenData["token"]
	if !exists || token == "" {
		return "", errors.New("token not found in file")
	}

	return token, nil
}
