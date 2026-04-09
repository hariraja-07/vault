package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const configFolderName = "vault"
const historyFileName = "history"
const keysFileName = "keys.txt"
const maxHistoryEntries = 100

type HistoryEntry struct {
	Timestamp int64
	Key       string
}

func GetConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Cannot find home directory:", err)
		os.Exit(1)
	}

	configDir := filepath.Join(home, ".config", configFolderName)

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.MkdirAll(configDir, 0755)
	}

	return configDir
}

func GetHistoryFilePath() string {
	return filepath.Join(GetConfigDir(), historyFileName)
}

func GetKeysFilePath() string {
	return filepath.Join(GetConfigDir(), keysFileName)
}

func TrackKeyUsage(key string) {
	historyFile := GetHistoryFilePath()

	entries, err := loadHistoryEntries(historyFile)
	if err != nil {
		entries = []HistoryEntry{}
	}

	entries = removeKeyFromHistory(entries, key)

	newEntry := HistoryEntry{
		Timestamp: time.Now().Unix(),
		Key:       key,
	}
	entries = append([]HistoryEntry{newEntry}, entries...)

	if len(entries) > maxHistoryEntries {
		entries = entries[:maxHistoryEntries]
	}

	saveHistoryEntries(historyFile, entries)

	updateKeysFile()
}

func loadHistoryEntries(historyFile string) ([]HistoryEntry, error) {
	var entries []HistoryEntry

	file, err := os.Open(historyFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}

		timestamp, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			continue
		}

		entries = append(entries, HistoryEntry{
			Timestamp: timestamp,
			Key:       parts[1],
		})
	}

	return entries, scanner.Err()
}

func saveHistoryEntries(historyFile string, entries []HistoryEntry) {
	file, err := os.Create(historyFile)
	if err != nil {
		fmt.Println("Error creating history file:", err)
		return
	}
	defer file.Close()

	for _, entry := range entries {
		fmt.Fprintf(file, "%d %s\n", entry.Timestamp, entry.Key)
	}
}

func removeKeyFromHistory(entries []HistoryEntry, key string) []HistoryEntry {
	var result []HistoryEntry
	for _, entry := range entries {
		if entry.Key != key {
			result = append(result, entry)
		}
	}
	return result
}

func GetRecentKeys(limit int) []string {
	entries, err := loadHistoryEntries(GetHistoryFilePath())
	if err != nil {
		return []string{}
	}

	if limit <= 0 {
		limit = 10
	}

	if len(entries) > limit {
		entries = entries[:limit]
	}

	var keys []string
	for _, entry := range entries {
		keys = append(keys, entry.Key)
	}

	return keys
}

func GetAllKeys() []string {
	data := LoadData()
	var keys []string

	for key, value := range data {
		if IsGroup(value) {
			keys = append(keys, key+"/")
		} else {
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)
	return keys
}

func updateKeysFile() {
	keysFile := GetKeysFilePath()

	data := make(map[string]interface{})
	filePath := GetVaultFilePath()
	file, err := os.ReadFile(filePath)
	if err == nil && len(file) > 0 {
		json.Unmarshal(file, &data)
	}

	var keys []string
	for key, value := range data {
		if !IsGroup(value) {
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)

	f, err := os.Create(keysFile)
	if err != nil {
		return
	}
	defer f.Close()

	for _, key := range keys {
		fmt.Fprintln(f, key)
	}
}

func InitConfig() {
	GetConfigDir()
	updateKeysFile()
}
