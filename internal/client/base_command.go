package client

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Base interface to command
type Command interface {
	Execute(config Config) error
}

// Factory to create command with arguments
type CommandFactory interface {
	Create(args map[string]string) (Command, error)
}

// Flag definition
type FlagDef struct {
	Name         string
	DefaultValue string
	Description  string
}

// Command registry
type CommandRegistry struct {
	config   Config
	commands map[string]CommandFactory
	rootCmd  *cobra.Command
}

func NewCommandRegistry(config Config, rootCmd *cobra.Command) *CommandRegistry {
	return &CommandRegistry{
		config:   config,
		commands: make(map[string]CommandFactory),
		rootCmd:  rootCmd,
	}
}

// Register new command
func (cr *CommandRegistry) Register(name string, factory CommandFactory, flags []FlagDef) {
	cr.commands[name] = factory

	cmd := &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("%s command", name),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags := make(map[string]string)
			cmd.Flags().VisitAll(func(f *pflag.Flag) {
				flags[f.Name] = f.Value.String()
			})

			command, err := factory.Create(flags)
			if err != nil {
				return err
			}
			return command.Execute(cr.config)
		},
	}

	for _, flag := range flags {
		cmd.Flags().String(flag.Name, flag.DefaultValue, flag.Description)
	}

	cr.rootCmd.AddCommand(cmd)
}
