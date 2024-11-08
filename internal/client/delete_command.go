package client

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
)

// Command to delete secret
type DeleteCommand struct {
	secretName string
}

func NewDeleteCommand(args map[string]string) (*DeleteCommand, error) {
	secretName, ok1 := args["name"]
	if !ok1 || secretName == "" {
		return nil, errors.New("secret name are required")
	}

	return &DeleteCommand{
		secretName: secretName,
	}, nil
}

func (cmd *DeleteCommand) Execute(config Config) error {
	token, err := readTokenFromFile()
	if err != nil {
		return fmt.Errorf("can`t find auth data: %w", err)
	}

	client := resty.New().SetTimeout(10 * time.Second)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+token).
		Delete(config.ServerAddress + "/api/user/secret/" + cmd.secretName)

	if err != nil {
		return fmt.Errorf("error during send request: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("can`t delete secret. Reason: %s", resp.String())
	}

	fmt.Printf("Secret with name \"%s\" was deleted", cmd.secretName)
	return nil
}

// Fabric to create secret delete command
type DeleteCommandFactory struct{}

func (f *DeleteCommandFactory) Create(args map[string]string) (Command, error) {
	return NewDeleteCommand(args)
}
