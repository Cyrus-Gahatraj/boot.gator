package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func write(cfg Config) error {
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	return os.WriteFile(path, jsonData, 0666)
}

func (cfg *Config) SetUser(name string) error {
	cfg.CurrentUserName = name
	return write(*cfg)
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", nil
	}

	path := home + "/" + configFileName
	return path, nil
}

func Read() Config {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}
	}

	data, err := os.Open(path)
	if err != nil {
		return Config{}
	}

	decoder := json.NewDecoder(data)
	var cfg Config
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}
	}

	return cfg
}
