package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const vaultFolderName = ".vault"
const vaultFileName = "data.json"

var initialized = false

func ensureInitialized() {
	if !initialized {
		InitConfig()
		initialized = true
	}
}

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
	ensureInitialized()
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

	migrateData(data)
	return data
}

func migrateData(data map[string]interface{}) {
	for _, value := range data {
		if m, ok := value.(map[string]interface{}); ok {
			if _, hasExpires := m["expires"]; !hasExpires {
				if _, hasCiphertext := m["ciphertext"]; hasCiphertext {
					m["expires"] = float64(0)
					m["once"] = false
				}
			}
			if _, hasOnce := m["once"]; !hasOnce {
				if _, hasCiphertext := m["ciphertext"]; hasCiphertext {
					m["once"] = false
				}
			}
		}
	}
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

func HasKey(data map[string]interface{}, key string) bool {
	_, exists := data[key]
	return exists
}

func HasGroup(data map[string]interface{}, key string) bool {
	if val, exists := data[key]; exists {
		_, isGroup := val.(map[string]interface{})
		return isGroup
	}
	return false
}

func IsGroup(value interface{}) bool {
	_, isGroup := value.(map[string]interface{})
	return isGroup
}

// EncryptedData represents encrypted data stored in JSON
type EncryptedData struct {
	Ciphertext string `json:"ciphertext"`
	Nonce      string `json:"nonce"`
}

// IsEncrypted checks if a value is encrypted (has ciphertext and nonce)
func IsEncrypted(value interface{}) bool {
	if m, ok := value.(map[string]interface{}); ok {
		_, hasCiphertext := m["ciphertext"]
		_, hasNonce := m["nonce"]
		return hasCiphertext && hasNonce
	}
	return false
}

// StoredValue represents a value with expiry metadata
type StoredValue struct {
	Value      interface{} `json:"value,omitempty"`
	Ciphertext string      `json:"ciphertext,omitempty"`
	Nonce      string      `json:"nonce,omitempty"`
	Expires    int64       `json:"expires"` // Unix timestamp, 0 = no expiry
	Once       bool        `json:"once"`    // Delete after first read
}

// GetExpires returns the expiry timestamp and a boolean indicating if it expires
func GetExpires(value interface{}) (int64, bool) {
	if m, ok := value.(map[string]interface{}); ok {
		if expires, ok := m["expires"].(float64); ok {
			return int64(expires), expires > 0
		}
	}
	return 0, false
}

// IsOnce returns true if the value should be deleted after first read
func IsOnce(value interface{}) bool {
	if m, ok := value.(map[string]interface{}); ok {
		if once, ok := m["once"].(bool); ok {
			return once
		}
	}
	return false
}

// IsExpired checks if a value has expired
func IsExpired(value interface{}) bool {
	if expires, ok := GetExpires(value); ok {
		return time.Now().Unix() >= expires
	}
	return false
}

// CleanupExpired removes all expired keys from data (silently)
func CleanupExpired(data map[string]interface{}) {
	for key, value := range data {
		if IsExpired(value) {
			delete(data, key)
			continue
		}
		if groupMap, ok := value.(map[string]interface{}); ok {
			for subKey, subValue := range groupMap {
				if IsExpired(subValue) {
					delete(groupMap, subKey)
				}
			}
			if len(groupMap) == 0 {
				delete(data, key)
			}
		}
	}
}
