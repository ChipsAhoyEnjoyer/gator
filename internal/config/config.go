package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	gatorConfigFile = ".gatorconfig.json"
)

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUsername string `json:"current_user_name"`
}

func Read() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error retreiving home directory from function 'os.UserHomeDir()': \n%v", err)
	}
	configJsonPath := homeDir + "/" + gatorConfigFile
	jsonData, err := os.ReadFile(configJsonPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file '%v': \n%v", jsonData, err)
	}
	c := &Config{}
	err = json.Unmarshal(jsonData, c)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling json data to Config struct: \n%v", err)
	}
	return c, nil
}

func (cfg *Config) SetUser() error {
	// {
	// 	"db_url": "postgres://example"
	// }
	jsonData, err := json.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("error marshaling json data from Config struct: \n%v", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error retreiving home directory from function 'os.UserHomeDir()': \n%v", err)
	}
	configJsonPath := homeDir + "/" + gatorConfigFile
	err = os.WriteFile(configJsonPath, jsonData, 0666)
	if err != nil {
		return fmt.Errorf(
			"error writing Config struct to config file; File maybe be empty or partially overwritten: \n%v",
			err,
		)
	}
	return nil
}
