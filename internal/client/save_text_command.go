package client

import (
	"errors"
	"fmt"
	"github.com/desepticon55/gophkeeper/internal/model"
	"github.com/desepticon55/gophkeeper/pkg/crypto"
	"github.com/go-resty/resty/v2"
	"time"
)

// Command to save TEXT secret
type SaveTextCommand struct {
	secretName string
	data       string
}

func NewSaveTextCommand(args map[string]string) (*SaveTextCommand, error) {
	secretName, ok1 := args["name"]
	data, ok2 := args["data"]
	if !ok1 || !ok2 || secretName == "" || data == "" {
		return nil, errors.New("secretName and data are required")
	}

	return &SaveTextCommand{
		secretName: secretName,
		data:       data,
	}, nil
}

func (cmd *SaveTextCommand) Execute(config Config) error {
	token, err := readTokenFromFile()
	if err != nil {
		return fmt.Errorf("can`t find auth data: %w", err)
	}

	encryptedData, err := crypto.EncryptData([]byte(cmd.data), []byte(config.EncryptionKey))
	if err != nil {
		return fmt.Errorf("error during encrypt data: %w", err)
	}

	secretPayload := &model.Secret{
		Name:    cmd.secretName,
		Type:    model.TextSecretType,
		Content: encryptedData,
	}

	client := resty.New().SetTimeout(10 * time.Second)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+token).
		SetBody(secretPayload).
		Post(config.ServerAddress + "/api/user/secret")

	if err != nil {
		return fmt.Errorf("error during send request: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("can`t save secret. Reason: %s", resp.String())
	}

	fmt.Println("Secret saved successfully")
	return nil
}

// Fabric to create TEXT secret command
type SaveTextCommandFactory struct{}

func (f *SaveTextCommandFactory) Create(args map[string]string) (Command, error) {
	return NewSaveTextCommand(args)
}
