package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigPath() string {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("error getting home directory")
		os.Exit(1)
	}
	return home_dir + "/" + configFileName
}

func Read() Config {
	config_file, err := os.ReadFile(getConfigPath())
	if err != nil {
		fmt.Printf("error reading config")
		os.Exit(1)
	}
	var config Config
	err = json.Unmarshal(config_file, &config)
	if err != nil {
		fmt.Printf("error unmarshalling config")
		os.Exit(1)
	}
	return config
}

func (cfg *Config) SetUser(username string) {
	cfg.CurrentUserName = username
	jsonConfig, err := json.Marshal(cfg)
	if err != nil {
		fmt.Printf("error setting config")
		os.Exit(1)
	}
	os.WriteFile(getConfigPath(), jsonConfig, 0644)
}
