package client

import (
	"errors"
	"fmt"
	"github.com/desepticon55/gophkeeper/internal/model"
	"github.com/go-resty/resty/v2"
	"time"
)

// Command to register new user
type UserRegisterCommand struct {
	username string
	password string
}

func NewUserRegisterCommand(args map[string]string) (*UserRegisterCommand, error) {
	username, ok1 := args["username"]
	password, ok2 := args["password"]
	if !ok1 || !ok2 || username == "" || password == "" {
		return nil, errors.New("username and password are required")
	}

	return &UserRegisterCommand{username: username, password: password}, nil
}

func (cmd *UserRegisterCommand) Execute(config Config) error {
	registerPayload := &model.User{
		Username: cmd.username,
		Password: cmd.password,
	}

	client := resty.New().SetTimeout(10 * time.Second)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(registerPayload).
		Post(config.ServerAddress + "/api/user/register")

	if err != nil {
		return fmt.Errorf("error during send request: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("error during register new user. Reason: %s. Status = %d", resp.String(), resp.StatusCode())
	}

	authHeader := resp.Header().Get("Authorization")
	if authHeader == "" {
		return errors.New("authorization token missing from response")
	}

	token := authHeader[len("Bearer "):]
	err = saveTokenToFile(token)
	if err != nil {
		return fmt.Errorf("failed to save token to file: %w", err)
	}

	fmt.Println("User logged in successfully")
	return nil
}

// Fabric to create user register command
type UserRegisterCommandFactory struct{}

func (f *UserRegisterCommandFactory) Create(args map[string]string) (Command, error) {
	return NewUserRegisterCommand(args)
}
