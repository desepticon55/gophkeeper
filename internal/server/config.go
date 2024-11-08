package server

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	ServerAddress      string
	DatabaseConnString string
	AuthKey            string
	ExpirationMinutes  int
}

func ParseConfig() Config {
	defaultAddress := "localhost:8080"
	if envAddr, exists := os.LookupEnv("ADDRESS"); exists {
		defaultAddress = envAddr
	}
	address := flag.String("a", defaultAddress, "Server address")

	defaultDatabaseConnString := "postgres://postgres:postgres@localhost:5432/postgres"
	if envDatabaseConnString, exists := os.LookupEnv("DATABASE_URI"); exists {
		defaultDatabaseConnString = envDatabaseConnString
	}
	databaseConnString := flag.String("d", defaultDatabaseConnString, "Database connection string")

	defaultAuthKey := "xiuw1bi4r98vd1(&*6"
	if envAuthKey, exists := os.LookupEnv("AUTH_KEY"); exists {
		defaultAuthKey = envAuthKey
	}
	authKey := flag.String("k", defaultAuthKey, "Auth key")

	defaultExpirationMinutes := 25
	if envExpirationMinutes, exists := os.LookupEnv("AUTH_EXPIRATION_TIME"); exists {
		if parsedExpirationMinutes, err := strconv.Atoi(envExpirationMinutes); err == nil {
			defaultExpirationMinutes = parsedExpirationMinutes
		}
	}
	expirationMinutes := flag.Int("e", defaultExpirationMinutes, "Expiration time (minutes)")

	flag.Parse()
	return Config{
		ServerAddress:      *address,
		DatabaseConnString: *databaseConnString,
		AuthKey:            *authKey,
		ExpirationMinutes:  *expirationMinutes,
	}
}
