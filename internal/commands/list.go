package commands

import (
	"fmt"
	"sort"
	"strconv"

	"vault/internal/storage"
)

func HandleList(args []string) {
	full := false
	recent := false
	recentLimit := 0
	group := ""

	for _, arg := range args[2:] {
		if arg == "--full" || arg == "-f" {
			full = true
		} else if arg == "--recent" {
			recent = true
			recentLimit = storage.GetRecentLimit()
		} else if arg == "--help" || arg == "-h" {
			showListHelp()
			return
		} else {
			parsedLimit, err := strconv.Atoi(arg)
			if err == nil && parsedLimit > 0 && recent {
				recentLimit = parsedLimit
			} else {
				group = arg
			}
		}
	}

	data := storage.LoadData()

	if recent {
		showRecentKeys(recentLimit)
		return
	}

	if group != "" {
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

		return
	}

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

		if groupMap, ok := value.(map[string]interface{}); ok {
			fmt.Printf("%s%s/\n", prefix, key)
			if full {
				var subKeys []string
				for subKey := range groupMap {
					subKeys = append(subKeys, subKey)
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
			fmt.Printf("%s%s\n", prefix, key)
		}
	}
}

func showRecentKeys(limit int) {
	recentKeys := storage.GetRecentKeys(limit)
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

func showListHelp() {
	fmt.Println("vault list - List all secrets or secrets in a specific group")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  vault list                List all secrets")
	fmt.Println("  vault list <group>       List secrets in a specific group")
	fmt.Println("  vault list --full        List all secrets with nested keys")
	fmt.Println("  vault list --full -f     Short form for --full")
	fmt.Println("  vault list --recent      List recent keys")
	fmt.Println("  vault list --recent 5    List 5 recent keys")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  vault list")
	fmt.Println("  vault list --full")
	fmt.Println("  vault list work")
	fmt.Println("  vault list --recent")
	fmt.Println("  vault list --recent 5")
}
