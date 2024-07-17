package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type APIConfig struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

type Config struct {
	ENV map[string]string

	Apis struct {
		Paprika  APIConfig
		CoinBase APIConfig
	}
}

func Crear() (Config, error) {

	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("error cargando las variable de entorno: %w", err)
	}

	configFile, err := os.Open(os.Getenv("CONFIG_FILE"))
	if err != nil {
		return Config{}, fmt.Errorf("failed to open config file: %v", err)
	}
	defer configFile.Close()

	var config Config
	if err := json.NewDecoder(configFile).Decode(&config.Apis); err != nil {
		return Config{}, fmt.Errorf("failed to decode config file: %v", err)
	}

	config.ENV = make(map[string]string)
	config.ENV["DB_USER"] = os.Getenv("DB_USER")
	config.ENV["DB_USER"] = os.Getenv("DB_USER")
	config.ENV["DB_PASS"] = os.Getenv("DB_PASS")
	config.ENV["DB_HOST"] = os.Getenv("DB_HOST")
	config.ENV["DB_PORT"] = os.Getenv("DB_PORT")
	config.ENV["DB_NAME"] = os.Getenv("DB_NAME")

	return config, nil

}
