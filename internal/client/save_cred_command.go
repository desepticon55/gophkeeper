package client

import (
	"errors"
	"fmt"
	"github.com/desepticon55/gophkeeper/internal/model"
	"github.com/desepticon55/gophkeeper/pkg/crypto"
	"github.com/go-resty/resty/v2"
	"time"
)

// Command to save CREDENTIALS secret
type SaveCredentialsCommand struct {
	secretName string
	username   string
	password   string
}

func NewSaveCredentialsCommand(args map[string]string) (*SaveCredentialsCommand, error) {
	secretName, ok1 := args["name"]
	username, ok2 := args["username"]
	password, ok3 := args["password"]
	if !ok1 || !ok2 || !ok3 || secretName == "" || username == "" || password == "" {
		return nil, errors.New("secretName and username and password are required")
	}

	return &SaveCredentialsCommand{
		secretName: secretName,
		username:   username,
		password:   password,
	}, nil
}

func (cmd *SaveCredentialsCommand) Execute(config Config) error {
	content := []byte(cmd.username + ":" + cmd.password)

	token, err := readTokenFromFile()
	if err != nil {
		return fmt.Errorf("can`t find auth data: %w", err)
	}

	encryptedData, err := crypto.EncryptData(content, []byte(config.EncryptionKey))
	if err != nil {
		return fmt.Errorf("error during encrypt data: %w", err)
	}

	secretPayload := &model.Secret{
		Name:    cmd.secretName,
		Type:    model.CredentialsSecretType,
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

// Fabric to create CREDENTIALS secret command
type SaveCredentialsCommandFactory struct{}

func (f *SaveCredentialsCommandFactory) Create(args map[string]string) (Command, error) {
	return NewSaveCredentialsCommand(args)
}
