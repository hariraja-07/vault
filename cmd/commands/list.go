package commands

import (
	"fmt"
	"sort"
	"time"

	"github.com/spf13/cobra"
	"vault/internal/storage"
)

var (
	listFull   bool
	listRecent bool
)

var ListCmd = &cobra.Command{
	Use:   "list [group]",
	Short: "List all secrets",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		data := storage.LoadData()
		storage.CleanupExpired(data)
		storage.SaveData(data)

		if listRecent {
			showRecentKeys()
			return
		}

		if len(args) == 1 {
			group := args[0]
			showGroup(group, data)
			return
		}

		showAll(data)
	},
}

func showAll(data map[string]interface{}) {
	if len(data) == 0 {
		fmt.Println("No data stored")
		return
	}

	var keys []string
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	fmt.Println("Vault")
	for i, key := range keys {
		value := data[key]
		isLast := i == len(keys)-1

		prefix := "├── "
		if isLast {
			prefix = "└── "
		}

		if isGroup(value) {
			fmt.Printf("%s%s/\n", prefix, key)
			if listFull {
				groupMap := value.(map[string]interface{})
				var subKeys []string
				for subKey, subValue := range groupMap {
					if !storage.IsExpired(subValue) && !isGroup(subValue) {
						subKeys = append(subKeys, subKey)
					}
				}
				sort.Strings(subKeys)

				for j, subKey := range subKeys {
					isSubLast := j == len(subKeys)-1

					if isLast {
						if isSubLast {
							fmt.Printf("    └── %s\n", subKey)
						} else {
							fmt.Printf("    ├── %s\n", subKey)
						}
					} else {
						if isSubLast {
							fmt.Printf("│   └── %s\n", subKey)
						} else {
							fmt.Printf("│   ├── %s\n", subKey)
						}
					}
				}
			}
		} else {
			expiry := getExpiryMarker(value)
			if expiry != "" {
				fmt.Printf("%s%s %s\n", prefix, key, expiry)
			} else {
				fmt.Printf("%s%s\n", prefix, key)
			}
		}
	}
}

func isGroup(value interface{}) bool {
	m, ok := value.(map[string]interface{})
	if !ok {
		return false
	}
	_, hasValue := m["value"]
	_, hasCiphertext := m["ciphertext"]
	return !hasValue && !hasCiphertext
}

func getExpiryMarker(value interface{}) string {
	if m, ok := value.(map[string]interface{}); ok {
		remaining := ""
		if expires, ok := m["expires"].(float64); ok && int64(expires) > 0 {
			rem := int64(expires) - time.Now().Unix()
			if rem > 0 {
				if rem < 60 {
					remaining = fmt.Sprintf("%ds", rem)
				} else if rem < 3600 {
					remaining = fmt.Sprintf("%dm", rem/60)
				} else if rem < 86400 {
					remaining = fmt.Sprintf("%dh", rem/3600)
				} else {
					remaining = fmt.Sprintf("%dd", rem/86400)
				}
			}
		}
		isOnce, _ := m["once"].(bool)
		if remaining != "" && isOnce {
			return fmt.Sprintf("[%s][once]", remaining)
		}
		if remaining != "" {
			return fmt.Sprintf("[%s]", remaining)
		}
		if isOnce {
			return "[once]"
		}
	}
	return ""
}

func showGroup(group string, data map[string]interface{}) {
	groupMap, ok := data[group].(map[string]interface{})
	if !ok {
		fmt.Println("Group not found")
		return
	}

	if len(groupMap) == 0 {
		fmt.Println("Group is empty")
		return
	}

	var keys []string
	for key := range groupMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	fmt.Println(group + "/")

	for i, key := range keys {
		isLast := i == len(keys)-1

		prefix := "├── "
		if isLast {
			prefix = "└── "
		}

		fmt.Printf("%s%s\n", prefix, key)
	}
}

func showRecentKeys() {
	recentLimit := storage.GetRecentLimit()
	recentKeys := storage.GetRecentKeys(recentLimit)
	allKeys := storage.GetAllKeys()

	if len(allKeys) == 0 {
		fmt.Println("No keys stored")
		return
	}

	if len(recentKeys) == 0 {
		fmt.Println("No recent keys")
		return
	}

	fmt.Println("Recent keys:")
	for i, key := range recentKeys {
		isLast := i == len(recentKeys)-1
		prefix := "├── "
		if isLast {
			prefix = "└── "
		}
		fmt.Printf("%s%s\n", prefix, key)
	}
}

func init() {
	ListCmd.Flags().BoolVarP(&listFull, "full", "f", false, "Show nested keys")
	ListCmd.Flags().BoolVar(&listRecent, "recent", false, "Show recent keys")
}
