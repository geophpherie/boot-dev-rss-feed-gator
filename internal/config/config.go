package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const CONFIG_FILE = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(userName string) error {
	c.CurrentUserName = userName

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	homeDir, err := getUserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, CONFIG_FILE)
	err = os.WriteFile(configPath, data, 600)
	if err != nil {
		return err
	}
	return nil
}

func Read() (Config, error) {
	homeDir, err := getUserHomeDir()
	if err != nil {
		return Config{}, err
	}
	configPath := filepath.Join(homeDir, CONFIG_FILE)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}
	config := Config{}
	json.Unmarshal(data, &config)

	return config, nil
}

func getUserHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home, nil
}
