package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gopasspw/clipboard"
	"github.com/spf13/cobra"
	"vault/internal/crypto"
	"vault/internal/storage"
)

var copyFlag bool

func init() {
	GetCmd.Flags().BoolVarP(&copyFlag, "copy", "c", false, "Copy the secret to clipboard")
}

var GetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a secret",
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completeKeys(cmd, args, toComplete)
	},
	Run: func(cmd *cobra.Command, args []string) {
		data := storage.LoadData()
		key := args[0]
		var isOnce bool
		var groupPath string

		if strings.Contains(key, "/") {
			if strings.Count(key, "/") > 1 {
				fmt.Println("Error: Only one level grouping allowed (group/key)")
				return
			}

			parts := strings.SplitN(key, "/", 2)
			group := parts[0]
			subKey := parts[1]
			groupPath = group

			groupMap, ok := data[group].(map[string]interface{})
			if !ok {
				fmt.Println("Group not found")
				return
			}
			if _, ok := groupMap[subKey]; !ok {
				fmt.Println("Key not found")
				return
			}
			val := groupMap[subKey]

			if m, ok := val.(map[string]interface{}); ok {
				if expires, ok := m["expires"].(float64); ok && int64(expires) > 0 {
					if time.Now().Unix() >= int64(expires) {
						delete(groupMap, subKey)
						storage.SaveData(data)
						return
					}
				}
				if once, ok := m["once"].(bool); ok {
					isOnce = once
				}
			}
		} else {
			if _, ok := data[key]; !ok {
				fmt.Println("Key not found")
				return
			}

			if m, ok := data[key].(map[string]interface{}); ok {
				if expires, ok := m["expires"].(float64); ok && int64(expires) > 0 {
					if time.Now().Unix() >= int64(expires) {
						delete(data, key)
						storage.SaveData(data)
						return
					}
				}
				if once, ok := m["once"].(bool); ok {
					isOnce = once
				}
			}
		}

		val := getValue(data, key)

		var result string

		if m, ok := val.(map[string]interface{}); ok {
			if _, hasCiphertext := m["ciphertext"]; hasCiphertext {
				password := askPassword("Enter password: ")
				ev := &crypto.EncryptedValue{
					Ciphertext: m["ciphertext"].(string),
					Nonce:      m["nonce"].(string),
				}
				decrypted, err := crypto.Decrypt(ev, password)
				if err != nil {
					fmt.Println("Error: Decryption failed. Wrong password?")
					return
				}
				result = decrypted
			} else if v, ok := m["value"].(string); ok {
				result = v
			} else {
				fmt.Println("Error: Invalid value format")
				return
			}
		} else if s, ok := val.(string); ok {
			result = s
		} else {
			result = fmt.Sprintf("%v", val)
		}

		fmt.Println(result)

		if copyFlag {
			if err := clipboard.WriteAllString(context.Background(), result); err != nil {
				fmt.Println("Error: Failed to copy to clipboard")
				return
			}
			fmt.Printf("Copied: %s\n", key)
			if isOnce {
				deleteKey(data, key, groupPath)
				storage.SaveData(data)
				fmt.Printf("%s has been removed (one-time)\n", key)
			}
		} else if isOnce {
			deleteKey(data, key, groupPath)
			storage.SaveData(data)
			fmt.Printf("%s has been removed (one-time)\n", key)
		}

		storage.TrackKeyUsage(key)
	},
}

func getValue(data map[string]interface{}, key string) interface{} {
	if strings.Contains(key, "/") {
		parts := strings.SplitN(key, "/", 2)
		group := parts[0]
		subKey := parts[1]
		groupMap := data[group].(map[string]interface{})
		return groupMap[subKey]
	}
	return data[key]
}

func deleteKey(data map[string]interface{}, key string, groupPath string) {
	if groupPath != "" {
		groupMap := data[groupPath].(map[string]interface{})
		delete(groupMap, strings.SplitN(key, "/", 2)[1])
	} else {
		delete(data, key)
	}
}
