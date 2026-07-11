package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DbUrl             string `json:"db_url"`
	CurrentUserName   string `json:"current_user_name"`
}

// read fron .gatorconfig.json in the ~/.gatorconfig.json folder
func ReadConfig() (*Config, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homedir, ".gatorconfig.json")
	
	data, err := os.ReadFile(configPath)
	
	if err != nil {
		return nil, err
	}
	
	var config Config
	
	err = json.Unmarshal(data, &config)
	
	if err != nil {
		return nil, err
	}
	
	return &config, nil
}

func (config *Config) SetUser(username string) error {
	homedir, err := os.UserHomeDir()

	if err != nil {
		return err
	}

	configPath := filepath.Join(homedir, ".gatorconfig.json")

	config.CurrentUserName = username

	data, err := json.MarshalIndent(config, "", "  ")

	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, data, 0600)
	
	if err != nil {
		return err
	}

	return nil
}