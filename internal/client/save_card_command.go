package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/desepticon55/gophkeeper/internal/model"
	"github.com/desepticon55/gophkeeper/pkg/crypto"
	"github.com/go-resty/resty/v2"
	"time"
)

// Command to save CARD secret
type SaveCardCommand struct {
	secretName     string
	cardRequisites model.Card
}

func NewSaveCardCommand(args map[string]string) (*SaveCardCommand, error) {
	secretName, ok1 := args["name"]
	number, ok2 := args["number"]
	date, ok3 := args["date"]
	code, ok4 := args["code"]
	holder, ok5 := args["holder"]
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || secretName == "" || number == "" || date == "" || code == "" || holder == "" {
		return nil, errors.New("card requisites are required")
	}

	return &SaveCardCommand{
		secretName: secretName,
		cardRequisites: model.Card{
			Number: number,
			Date:   date,
			Code:   code,
			Holder: holder,
		},
	}, nil
}

func (cmd *SaveCardCommand) Execute(config Config) error {
	content, err := json.Marshal(cmd.cardRequisites)
	if err != nil {
		return fmt.Errorf("can`t serialise card requisites: %w", err)
	}

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
		Type:    model.CardSecretType,
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

// Fabric to create CARD secret command
type SaveCardCommandFactory struct{}

func (f *SaveCardCommandFactory) Create(args map[string]string) (Command, error) {
	return NewSaveCardCommand(args)
}
