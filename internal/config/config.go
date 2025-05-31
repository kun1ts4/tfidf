package config

import (
	"github.com/lpernett/godotenv"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"os"
)

type AppConfig struct {
	App struct {
		Port    string `yaml:"port"`
		Version string `yaml:"version"`
	} `yaml:"app"`
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
		User     string
		Password string
	} `yaml:"database"`
}

func LoadConfig() (*AppConfig, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("ошибка загрузки файла .env: %v", err)
	}

	configFile, err := os.Open("config.yaml")
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	data, err := io.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	var config AppConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	config.Database.User = os.Getenv("DB_USER")
	config.Database.Password = os.Getenv("DB_PASSWORD")

	return &config, nil
}
