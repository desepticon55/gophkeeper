package client

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testCommandFactory struct {
	expectedFlags map[string]string
	executed      bool
}

func (f *testCommandFactory) Create(flags map[string]string) (Command, error) {
	if f.expectedFlags != nil {
		for key, value := range f.expectedFlags {
			if flags[key] != value {
				return nil, fmt.Errorf("unexpected flag value: %s=%s", key, flags[key])
			}
		}
	}
	return &testCommand{factory: f}, nil
}

type testCommand struct {
	factory *testCommandFactory
}

func (cmd *testCommand) Execute(config Config) error {
	cmd.factory.executed = true
	return nil
}

func TestCommandRegistry_RegisterAndExecute(t *testing.T) {
	rootCmd := &cobra.Command{Use: "testapp"}
	config := Config{ServerAddress: "http://localhost:8080"}
	registry := NewCommandRegistry(config, rootCmd)

	expectedFlags := map[string]string{
		"flag1": "value1",
		"flag2": "value2",
	}
	factory := &testCommandFactory{expectedFlags: expectedFlags}
	flags := []FlagDef{
		{Name: "flag1", DefaultValue: "default1", Description: "Test flag 1"},
		{Name: "flag2", DefaultValue: "default2", Description: "Test flag 2"},
	}

	registry.Register("test", factory, flags)

	args := []string{"test", "--flag1=value1", "--flag2=value2"}
	rootCmd.SetArgs(args)

	err := rootCmd.Execute()
	assert.NoError(t, err)

	assert.True(t, factory.executed, "Command was not executed")
}
