package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	var cfg Config

	configPath, err := getConfigFilePath()
	if err != nil {
		return cfg, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	return write(*c)
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configPath := homeDir + "/" + configFileName

	return configPath, nil
}

func write(cfg Config) error {
	// get file path which we are going to write to
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	// convert config (go struct) to json
	data, err := json.MarshalIndent(cfg, "", "")
	if err != nil {
		return err
	}

	// write json data to file
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
