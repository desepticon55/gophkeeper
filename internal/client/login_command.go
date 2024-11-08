package client

import (
	"errors"
	"fmt"
	"github.com/desepticon55/gophkeeper/internal/model"
	"github.com/go-resty/resty/v2"
	"time"
)

// Command to login
type UserLoginCommand struct {
	username string
	password string
}

func NewUserLoginCommand(args map[string]string) (*UserLoginCommand, error) {
	username, ok1 := args["username"]
	password, ok2 := args["password"]
	if !ok1 || !ok2 || username == "" || password == "" {
		return nil, errors.New("username and password are required")
	}

	return &UserLoginCommand{username: username, password: password}, nil
}

func (cmd *UserLoginCommand) Execute(config Config) error {
	loginPayload := &model.User{
		Username: cmd.username,
		Password: cmd.password,
	}

	client := resty.New().SetTimeout(10 * time.Second)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(loginPayload).
		Post(config.ServerAddress + "/api/user/login")

	if err != nil {
		return fmt.Errorf("error during send request: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("error during execute command. Reason: %s", resp.String())
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

// Fabric to create user login command
type UserLoginCommandFactory struct{}

func (f *UserLoginCommandFactory) Create(args map[string]string) (Command, error) {
	return NewUserLoginCommand(args)
}
