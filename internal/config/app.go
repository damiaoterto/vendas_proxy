package config

import (
	"fmt"
	"os"
	"strings"
)

type MongoDB struct {
	URI string
}

type AppConfig struct {
	MongoDB
}

func NewConfig() (*AppConfig, error) {
	mongoURI := os.Getenv("MONGODB_URI")
	if strings.Trim(mongoURI, " ") == "" {
		return nil, fmt.Errorf("invalid value for MONGODB_URI")

	}

	return &AppConfig{
		MongoDB: MongoDB{
			URI: mongoURI,
		},
	}, nil
}
