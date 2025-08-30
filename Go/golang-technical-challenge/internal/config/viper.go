package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	v := viper.New()

	rootPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	if filepath.Base(rootPath) == "web" {
		rootPath = filepath.Dir(filepath.Dir(rootPath))
	}

	envPath := filepath.Join(rootPath, ".env")
	v.SetConfigFile(envPath)

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read configuration file at %s: %v", envPath, err)
	}

	log.Printf("Loaded configuration from: %s", v.ConfigFileUsed())
	return v
}
