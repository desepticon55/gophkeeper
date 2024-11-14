package client

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress string
	EncryptionKey string
}

func ParseConfig() Config {
	defaultAddress := "http://localhost:8080"
	if envAddr, exists := os.LookupEnv("ADDRESS"); exists {
		defaultAddress = envAddr
	}
	address := flag.String("a", defaultAddress, "Server address")

	defaultEncryptionKey := "WYJcWgkItShq513L21E1CFuz6uQWDy3p"
	if envEncryptionKey, exists := os.LookupEnv("ENCRYPTION_KEY"); exists {
		defaultEncryptionKey = envEncryptionKey
	}
	key := flag.String("k", defaultEncryptionKey, "Encryption key")

	flag.Parse()
	return Config{
		ServerAddress: *address,
		EncryptionKey: *key,
	}
}
