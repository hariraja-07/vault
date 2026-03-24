package main

import (
	"encoding/json"
	"os"
)

func loadData() map[string]string {
	data := make(map[string]string)

	file, err := os.ReadFile("data.json")

	if err != nil {
		return data
	}

	json.Unmarshal(file, &data)
	return data
}

func saveData(data map[string]string) {
	file, err := json.MarshalIndent(data, "", "  ")

	if err != nil {
		return
	}

	os.WriteFile("data.json", file, 0644)
}
