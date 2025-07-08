package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	MongoDB struct {
		URI        string
		Database   string
		Collection string
	}
}

func Load() (*AppConfig, error) {
	cwdDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("fail on get workdir")

	}

	envDir := filepath.Join(cwdDir, ".env")
	_, err = os.Stat(envDir)
	if os.IsNotExist(err) {
		log.Println("The .env file not exists")
	} else {
		godotenv.Load()
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if strings.Trim(mongoURI, " ") == "" {
		return nil, fmt.Errorf("invalid value for MONGODB_URI")

	}

	config := &AppConfig{}
	config.MongoDB.URI = mongoURI
	config.MongoDB.Database = os.Getenv("MONGODB_DATABASE")
	config.MongoDB.Collection = os.Getenv("MONGODB_COLLECTION")

	return config, nil
}
