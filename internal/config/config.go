package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	gatorConfigFile = ".gatorconfig.json"
)

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUsername string `json:"current_user_name"`
}

func Read() (*Config, error) {
	configJsonPath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}
	jsonData, err := os.ReadFile(configJsonPath)
	if err != nil {
		return nil, fmt.Errorf("error: cannot read config file '%v': \n%v", jsonData, err)
	}
	c := &Config{}
	err = json.Unmarshal(jsonData, c)
	if err != nil {
		return nil, fmt.Errorf("error: cannot unmarshal json data to Config struct: \n%v", err)
	}
	return c, nil
}

func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUsername = username
	err := write(cfg)
	if err != nil {
		return fmt.Errorf("error: cannot setting new user: \n%v", err)
	}
	return nil
}

func write(cfg *Config) error {
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error: cannot marshal json data from Config struct: \n%v", err)
	}
	configJsonPath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	err = os.WriteFile(configJsonPath, jsonData, 0666)
	if err != nil {
		return fmt.Errorf(
			"error: cannot write Config struct to config file; File maybe be empty or partially overwritten: \n%v",
			err,
		)
	}
	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error: cannot retreive home directory\n%v", err)
	}
	return filepath.Join(homeDir, gatorConfigFile), nil
}
