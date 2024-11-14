package client

import (
	"errors"
	"fmt"
	"github.com/desepticon55/gophkeeper/internal/model"
	"github.com/desepticon55/gophkeeper/pkg/crypto"
	"github.com/go-resty/resty/v2"
	"time"
)

// Command to read secret
type ReadCommand struct {
	secretName string
}

func NewReadCommand(args map[string]string) (*ReadCommand, error) {
	secretName, ok1 := args["name"]
	if !ok1 || secretName == "" {
		return nil, errors.New("secret name are required")
	}

	return &ReadCommand{
		secretName: secretName,
	}, nil
}

func (cmd *ReadCommand) Execute(config Config) error {
	token, err := readTokenFromFile()
	if err != nil {
		return fmt.Errorf("can`t find auth data: %w", err)
	}

	var secretPayload = &model.Secret{}

	client := resty.New().SetTimeout(10 * time.Second)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+token).
		SetResult(secretPayload).
		Get(config.ServerAddress + "/api/user/secret/" + cmd.secretName)

	if err != nil {
		return fmt.Errorf("error during send request: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("can`t save secret. Reason: %s", resp.String())
	}

	data, err := crypto.DecryptData(secretPayload.Content, []byte(config.EncryptionKey))
	if err != nil {
		return fmt.Errorf("error during decrypt content: %w", err)
	}
	fmt.Printf("Your secret name: %s\nYour secret value: %s\n", secretPayload.Name, data)
	return nil
}

// Fabric to create secret read command
type ReadCommandFactory struct{}

func (f *ReadCommandFactory) Create(args map[string]string) (Command, error) {
	return NewReadCommand(args)
}
