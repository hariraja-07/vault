package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const vaultFolderName = ".vault"
const vaultFileName = "data.json"

func GetVaultFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Cannot find home directory:", err)
		os.Exit(1)
	}

	vaultDir := filepath.Join(home, vaultFolderName)

	if _, err := os.Stat(vaultDir); os.IsNotExist(err) {
		os.Mkdir(vaultDir, 0755)
	}

	return filepath.Join(vaultDir, vaultFileName)
}

func LoadData() map[string]interface{} {
	filePath := GetVaultFilePath()
	data := make(map[string]interface{})

	file, err := os.ReadFile(filePath)

	if err != nil {
		if os.IsNotExist(err) {
			return data
		}
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	if len(file) == 0 {
		return data
	}

	err = json.Unmarshal(file, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}
	return data
}

func SaveData(data map[string]interface{}) {
	filePath := GetVaultFilePath()

	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		os.Exit(1)
	}

	err = os.WriteFile(filePath, file, 0644)

	if err != nil {
		fmt.Println("Error writing file:", err)
		os.Exit(1)
	}
}
