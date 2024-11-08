package main

import (
	"fmt"
	"github.com/desepticon55/gophkeeper/internal/client"
	"github.com/desepticon55/gophkeeper/pkg/logger"
	"github.com/desepticon55/gophkeeper/pkg/version"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

func main() {
	log := logger.InitLogger()
	config := client.ParseConfig()
	rootCmd := &cobra.Command{}

	fmt.Println(version.MakeBuildInfo(log))

	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Auth commands",
	}

	authRegistry := client.NewCommandRegistry(config, authCmd)
	authRegistry.Register("register", &client.UserRegisterCommandFactory{}, []client.FlagDef{
		{Name: "username", DefaultValue: "", Description: "User email for login"},
		{Name: "password", DefaultValue: "", Description: "User password for login"},
	})
	authRegistry.Register("login", &client.UserLoginCommandFactory{}, []client.FlagDef{
		{Name: "username", DefaultValue: "", Description: "User email for login"},
		{Name: "password", DefaultValue: "", Description: "User password for login"},
	})

	secretCmd := &cobra.Command{
		Use:   "secret",
		Short: "Secret commands",
	}

	secretRegistry := client.NewCommandRegistry(config, secretCmd)
	secretRegistry.Register("read", &client.ReadCommandFactory{}, []client.FlagDef{
		{Name: "name", DefaultValue: "", Description: "Secret name"},
	})
	secretRegistry.Register("delete", &client.DeleteCommandFactory{}, []client.FlagDef{
		{Name: "name", DefaultValue: "", Description: "Secret name"},
	})

	secretCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Secret create commands",
	}

	secretCreateRegistry := client.NewCommandRegistry(config, secretCreateCmd)
	secretCreateRegistry.Register("credentials", &client.SaveCredentialsCommandFactory{}, []client.FlagDef{
		{Name: "name", DefaultValue: "", Description: "Secret name"},
		{Name: "username", DefaultValue: "", Description: "User email for login"},
		{Name: "password", DefaultValue: "", Description: "User password for login"},
	})
	secretCreateRegistry.Register("card", &client.SaveCardCommandFactory{}, []client.FlagDef{
		{Name: "name", DefaultValue: "", Description: "Secret name"},
		{Name: "number", DefaultValue: "", Description: "Card number"},
		{Name: "date", DefaultValue: "", Description: "Card expire date"},
		{Name: "code", DefaultValue: "", Description: "CVC code"},
		{Name: "holder", DefaultValue: "", Description: "Holder"},
	})
	secretCreateRegistry.Register("text", &client.SaveTextCommandFactory{}, []client.FlagDef{
		{Name: "name", DefaultValue: "", Description: "Secret name"},
		{Name: "data", DefaultValue: "", Description: "Data"},
	})

	secretCmd.AddCommand(secretCreateCmd)
	rootCmd.AddCommand(authCmd, secretCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Error("Error during execute command", zap.Error(err))
		os.Exit(1)
	}
}
