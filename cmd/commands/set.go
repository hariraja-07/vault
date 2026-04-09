package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"vault/internal/crypto"
	"vault/internal/storage"
)

func askPassword(prompt string) string {
	fmt.Print(prompt)
	var password string
	fmt.Scanln(&password)
	return password
}

var setForce bool
var setSecure bool

var SetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a key-value pair",
	Args:  cobra.ExactArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completeKeys(cmd, args, toComplete)
	},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]
		force := setForce
		secure := setSecure

		data := storage.LoadData()

		// Encrypt value if --secure flag is set
		if secure {
			password := askPassword("Enter password for this key: ")
			if password == "" {
				fmt.Println("Error: Password cannot be empty")
				return
			}

			confirm := askPassword("Confirm password: ")
			if password != confirm {
				fmt.Println("Error: Passwords do not match")
				return
			}

			encrypted, err := crypto.Encrypt(value, password)
			if err != nil {
				fmt.Printf("Error: Failed to encrypt value: %v\n", err)
				return
			}

			value = ""
			setValue := map[string]interface{}{
				"ciphertext": encrypted.Ciphertext,
				"nonce":      encrypted.Nonce,
			}
			storeEncryptedValue(data, key, force, setValue)
		} else {
			storeValue(data, key, value, force)
		}

		storage.SaveData(data)
		storage.TrackKeyUsage(key)
		fmt.Println("Saved:", key)
	},
}

func storeValue(data map[string]interface{}, key string, value string, force bool) {
	if strings.Contains(key, "/") {
		parts := strings.SplitN(key, "/", 2)
		group := parts[0]
		subKey := parts[1]

		if existingGroup, exists := data[group]; exists {
			if !storage.IsGroup(existingGroup) {
				if !force {
					fmt.Printf("Error: Key '%s' already exists.\n", group)
					fmt.Println("Use --force or -F to overwrite.")
					return
				}
				delete(data, group)
				fmt.Printf("Warning: overwriting key '%s'\n", group)
			} else {
				groupMap := existingGroup.(map[string]interface{})
				if _, subKeyExists := groupMap[subKey]; subKeyExists {
					if !force {
						fmt.Printf("Error: Subkey '%s' already exists in group '%s'.\n", subKey, group)
						fmt.Println("Use --force or -F to overwrite.")
						return
					}
					fmt.Printf("Warning: overwriting subkey '%s'\n", subKey)
				}
			}
		}

		if _, exists := data[group]; !exists {
			data[group] = map[string]interface{}{}
		}

		groupMap := data[group].(map[string]interface{})
		groupMap[subKey] = value
	} else {
		if existingValue, exists := data[key]; exists {
			if !force {
				if storage.IsGroup(existingValue) {
					groupMap := existingValue.(map[string]interface{})
					count := len(groupMap)

					fmt.Printf("Error: Group '%s' already exists with %d nested key(s).\n", key, count)
					fmt.Println("Use --force or -F to delete all nested keys and overwrite.")
					return
				}
				fmt.Printf("Error: Key '%s' already exists.\n", key)
				fmt.Println("Use --force or -F to overwrite.")
				return
			}

			if storage.IsGroup(existingValue) {
				groupMap := existingValue.(map[string]interface{})
				count := len(groupMap)

				fmt.Printf("Error: Group '%s' already exists with %d nested key(s).\n", key, count)
				fmt.Println("Use --force or -F to delete all nested keys and overwrite.")
				return
			}

			if storage.IsGroup(existingValue) {
				groupMap := existingValue.(map[string]interface{})
				count := len(groupMap)
				fmt.Printf("Warning: overwriting group '%s' (%d nested key(s) will be deleted)\n", key, count)
			} else {
				fmt.Printf("Warning: overwriting key '%s'\n", key)
			}
			delete(data, key)
		}

		data[key] = value
	}
}

func storeEncryptedValue(data map[string]interface{}, key string, force bool, encryptedValue map[string]interface{}) {
	if strings.Contains(key, "/") {
		parts := strings.SplitN(key, "/", 2)
		group := parts[0]
		subKey := parts[1]

		if _, exists := data[group]; !exists {
			data[group] = map[string]interface{}{}
		}

		groupMap := data[group].(map[string]interface{})
		if _, exists := groupMap[subKey]; exists && !force {
			fmt.Printf("Error: Subkey '%s' already exists in group '%s'.\n", subKey, group)
			fmt.Println("Use --force or -F to overwrite.")
			return
		}

		groupMap[subKey] = encryptedValue
	} else {
		if existingValue, exists := data[key]; exists {
			if !force {
				if storage.IsGroup(existingValue) {
					groupMap := existingValue.(map[string]interface{})
					count := len(groupMap)
					fmt.Printf("Error: Group '%s' already exists with %d nested key(s).\n", key, count)
					fmt.Println("Use --force or -F to delete all nested keys and overwrite.")
					return
				}
				fmt.Printf("Error: Key '%s' already exists.\n", key)
				fmt.Println("Use --force or -F to overwrite.")
				return
			}
			delete(data, key)
		}

		data[key] = encryptedValue
	}
}

func init() {
	SetCmd.Flags().BoolVarP(&setForce, "force", "F", false, "Force overwrite existing key")
	SetCmd.Flags().BoolVarP(&setSecure, "secure", "s", false, "Encrypt this value with a password")
	RegisterKeyCompletion(SetCmd)
}
