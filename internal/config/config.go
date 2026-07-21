package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username

	return write(*c)
}

func write(cfg Config) error {
	cfgFilePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("problem detecting user home directory")
	}

	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("problem marshal json: %w", err)
	}

	if err = os.WriteFile(cfgFilePath, jsonData, 0644); err != nil {
		return fmt.Errorf("problem writing file: %w", err)
	}

	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed detecting user home directory")
	}

	filePath := filepath.Join(homeDir, configFileName)
	return filePath, nil
}

// ////////////////////////////////////////////
// ////////////////////////////////////////////
// ////////////////////////////////////////////

func Read() (Config, error) {
	cfgFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(cfgFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("problem reading file: %w", err)
	}

	var config Config
	if err = json.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("problem unmarshal json: %w", err)
	}

	return config, nil
}
