package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = "config.json"

type Config struct {
	RecentLimit int `json:"recentLimit"`
}

var defaultConfig = Config{
	RecentLimit: 10,
}

func GetConfigFilePath() string {
	return filepath.Join(GetConfigDir(), configFileName)
}

func LoadConfig() Config {
	configFile := GetConfigFilePath()

	file, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultConfig
		}
		return defaultConfig
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return defaultConfig
	}

	if config.RecentLimit <= 0 {
		config.RecentLimit = 10
	}

	return config
}

func SaveConfig(config Config) {
	configFile := GetConfigFilePath()

	if config.RecentLimit <= 0 {
		config.RecentLimit = 10
	}

	file, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println("Error encoding config:", err)
		return
	}

	err = os.WriteFile(configFile, file, 0644)
	if err != nil {
		fmt.Println("Error writing config:", err)
	}
}

func SetRecentLimit(limit int) {
	config := LoadConfig()
	config.RecentLimit = limit
	SaveConfig(config)
}

func GetRecentLimit() int {
	config := LoadConfig()
	return config.RecentLimit
}
