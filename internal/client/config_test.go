package client

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name                  string
		envAddress            string
		envEncryptionKey      string
		cmdArgs               []string
		expectedAddress       string
		expectedEncryptionKey string
	}{
		{
			name:                  "Default values",
			envAddress:            "",
			envEncryptionKey:      "",
			cmdArgs:               []string{},
			expectedAddress:       "http://localhost:8080",
			expectedEncryptionKey: "WYJcWgkItShq513L21E1CFuz6uQWDy3p",
		},
		{
			name:                  "Environment variables only",
			envAddress:            "http://localhost:9090",
			envEncryptionKey:      "WYJcWgkItShq513L21E1CFuz6uQWDy5p",
			cmdArgs:               []string{},
			expectedAddress:       "http://localhost:9090",
			expectedEncryptionKey: "WYJcWgkItShq513L21E1CFuz6uQWDy5p",
		},
		{
			name:                  "Command-line flags only",
			envAddress:            "",
			envEncryptionKey:      "",
			cmdArgs:               []string{"-a", "http://192.168.1.1:8081", "-k", "some_key"},
			expectedAddress:       "http://192.168.1.1:8081",
			expectedEncryptionKey: "some_key",
		},
		{
			name:                  "Environment variables and command-line flags",
			envAddress:            "http://localhost:9090",
			envEncryptionKey:      "WYJcWgkItShq513L21E1CFuz6uQWDy5p",
			cmdArgs:               []string{"-a", "http://192.168.1.1:8081", "-k", "some_key"},
			expectedAddress:       "http://192.168.1.1:8081",
			expectedEncryptionKey: "some_key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envAddress != "" {
				os.Setenv("ADDRESS", tt.envAddress)
				defer os.Unsetenv("ADDRESS")
			}
			if tt.expectedEncryptionKey != "" {
				os.Setenv("ENCRYPTION_KEY", tt.expectedEncryptionKey)
				defer os.Unsetenv("ENCRYPTION_KEY")
			}

			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ExitOnError)
			os.Args = append([]string{"cmd"}, tt.cmdArgs...)

			config := ParseConfig()
			assert.Equal(t, tt.expectedAddress, config.ServerAddress)
			assert.Equal(t, tt.expectedEncryptionKey, config.EncryptionKey)
		})
	}
}
