package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func Load() (*AppConfig, error) {
	cfg := &AppConfig{}

	_ = godotenv.Load()

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("ошибка парсинга ENV: %w", err)
	}

	/*yamlFile, err := os.ReadFile("configs/config.yml")
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать config.yaml: %w", err)
	}

	if err := yaml.Unmarshal(yamlFile, cfg); err != nil {
		return nil, fmt.Errorf("ошибка парсинга config.yaml: %w", err)
	}*/

	return cfg, nil
}
