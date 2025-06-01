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
		Port    string
		Version string `yaml:"version"`
	} `yaml:"app"`
	Database struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
	}
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
	config.App.Port = os.Getenv("API_PORT")
	config.Database.Port = os.Getenv("DB_PORT")
	config.Database.Name = os.Getenv("DB_NAME")
	config.Database.User = os.Getenv("DB_USER")
	config.Database.Password = os.Getenv("DB_PASSWORD")

	return &config, nil
}
