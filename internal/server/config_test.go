package server

import (
	"flag"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name                      string
		envAddress                string
		envDatabaseConn           string
		envAuthKey                string
		envExpirationMinutes      int
		cmdArgs                   []string
		expectedAddress           string
		expectedConnString        string
		expectedAuthKey           string
		expectedExpirationMinutes int
	}{
		{
			name:                      "Default values",
			envAddress:                "",
			envDatabaseConn:           "",
			envAuthKey:                "",
			envExpirationMinutes:      0,
			cmdArgs:                   []string{},
			expectedAddress:           "localhost:8080",
			expectedConnString:        "postgres://postgres:postgres@localhost:5432/postgres",
			expectedAuthKey:           "xiuw1bi4r98vd1(&*6",
			expectedExpirationMinutes: 25,
		},
		{
			name:                      "Environment variables only",
			envAddress:                "127.0.0.1:9090",
			envDatabaseConn:           "postgres://user:pass@localhost:5432/mydb",
			envAuthKey:                "some_key",
			envExpirationMinutes:      4,
			cmdArgs:                   []string{},
			expectedAddress:           "127.0.0.1:9090",
			expectedConnString:        "postgres://user:pass@localhost:5432/mydb",
			expectedAuthKey:           "some_key",
			expectedExpirationMinutes: 4,
		},
		{
			name:                      "Command-line flags only",
			envAddress:                "",
			envDatabaseConn:           "",
			envAuthKey:                "",
			envExpirationMinutes:      0,
			cmdArgs:                   []string{"-a", "192.168.1.1:8081", "-d", "postgres://admin:admin@db:5432/testdb", "-k", "some_key", "-e", "4"},
			expectedAddress:           "192.168.1.1:8081",
			expectedConnString:        "postgres://admin:admin@db:5432/testdb",
			expectedAuthKey:           "some_key",
			expectedExpirationMinutes: 4,
		},
		{
			name:                      "Environment variables and command-line flags",
			envAddress:                "127.0.0.1:9090",
			envDatabaseConn:           "postgres://user:pass@localhost:5432/mydb",
			envAuthKey:                "some_key",
			envExpirationMinutes:      4,
			cmdArgs:                   []string{"-a", "192.168.1.1:8081", "-d", "postgres://admin:admin@db:5432/testdb", "-k", "some_key1", "-e", "5"},
			expectedAddress:           "192.168.1.1:8081",
			expectedConnString:        "postgres://admin:admin@db:5432/testdb",
			expectedAuthKey:           "some_key1",
			expectedExpirationMinutes: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envAddress != "" {
				os.Setenv("ADDRESS", tt.envAddress)
				defer os.Unsetenv("ADDRESS")
			}
			if tt.envDatabaseConn != "" {
				os.Setenv("DATABASE_URI", tt.envDatabaseConn)
				defer os.Unsetenv("DATABASE_URI")
			}
			if tt.envAuthKey != "" {
				os.Setenv("AUTH_KEY", tt.envAuthKey)
				defer os.Unsetenv("AUTH_KEY")
			}
			if tt.envExpirationMinutes != 0 {
				os.Setenv("AUTH_EXPIRATION_TIME", fmt.Sprint(tt.envExpirationMinutes))
				defer os.Unsetenv("AUTH_EXPIRATION_TIME")
			}

			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ExitOnError)
			os.Args = append([]string{"cmd"}, tt.cmdArgs...)

			config := ParseConfig()
			assert.Equal(t, tt.expectedAddress, config.ServerAddress)
			assert.Equal(t, tt.expectedConnString, config.DatabaseConnString)
			assert.Equal(t, tt.expectedAuthKey, config.AuthKey)
			assert.Equal(t, tt.expectedExpirationMinutes, config.ExpirationMinutes)
		})
	}
}
