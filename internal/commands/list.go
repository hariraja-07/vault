package commands

import (
	"fmt"
	"sort"

	"vault/internal/storage"
)

func HandleList(args []string) {
	full := false
	group := ""

	for _, arg := range args[2:] {
		if arg == "--full" || arg == "-f" {
			full = true
		} else {
			group = arg
		}
	}

	data := storage.LoadData()

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
