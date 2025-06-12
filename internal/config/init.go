package config

import (
	"fmt"
	"os"
	"tfidf/docs"
)

func Initialize() (*AppConfig, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить конфигурацию: %v", err)
	}

	docs.SwaggerInfo.Version = cfg.App.Version
	docs.SwaggerInfo.Host = os.Getenv("API_HOST") + ":" + cfg.App.Port
	return cfg, nil
}
