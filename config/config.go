package config

import (
	"os"

	"github.com/joho/godotenv"
)

var keys = [4]string{"token", "playlist", "clientId", "clientSecret"}

func GetConfig() (map[string]string, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return nil, err
	}

	config := map[string]string{}

	for _, key := range keys {
		value := os.Getenv(key)
		config[key] = value
	}

	return config, nil
}