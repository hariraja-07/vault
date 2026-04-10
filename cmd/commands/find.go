package commands

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"vault/internal/storage"
)

var findSelect bool
var findGroup string
var findLimit int
var findCopy bool

var FindCmd = &cobra.Command{
	Use:   "find <terms...>",
	Short: "Find keys by fuzzy search",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		data := storage.LoadData()
		terms := args

		matches := findKeys(data, terms, findGroup)

		if len(matches) == 0 {
			fmt.Printf("No keys found matching '%s'\n", strings.Join(terms, ", "))
			return
		}

		sort.Strings(matches)

		limit := findLimit
		if limit == 0 {
			limit = len(matches)
		}

		displayLimit := 10
		if limit < 10 {
			displayLimit = limit
		}

		shown := 0
		scanner := bufio.NewScanner(os.Stdin)

		for shown < limit && shown < len(matches) {
			end := shown + displayLimit
			if end > limit {
				end = limit
			}
			if end > len(matches) {
				end = len(matches)
			}

			for i := shown; i < end; i++ {
				fmt.Printf("[%d] %s\n", i+1, matches[i])
			}

			shown = end

			if shown < limit && shown < len(matches) {
				fmt.Printf("(%d/%d) more (y/N): ", shown, len(matches))

				if !scanner.Scan() {
					break
				}
				response := strings.ToLower(strings.TrimSpace(scanner.Text()))

				if response != "y" && response != "yes" {
					break
				}
			}
		}
	},
}

func findKeys(data map[string]interface{}, terms []string, groupFilter string) []string {
	var matches []string

	if groupFilter != "" {
		groupMap, ok := data[groupFilter].(map[string]interface{})
		if !ok {
			return matches
		}
		for key := range groupMap {
			if fuzzyMatch(key, terms) {
				matches = append(matches, fmt.Sprintf("%s/%s", groupFilter, key))
			}
		}
	} else {
		for key, value := range data {
			if fuzzyMatch(key, terms) {
				matches = append(matches, key)
			}

			if groupMap, ok := value.(map[string]interface{}); ok {
				for subKey := range groupMap {
					if fuzzyMatch(subKey, terms) {
						matches = append(matches, fmt.Sprintf("%s/%s", key, subKey))
					}
				}
			}
		}
	}

	return matches
}

func fuzzyMatch(key string, terms []string) bool {
	keyLower := strings.ToLower(key)
	for _, term := range terms {
		if strings.Contains(keyLower, strings.ToLower(term)) {
			return true
		}
	}
	return false
}

func init() {
	FindCmd.Flags().BoolVarP(&findSelect, "select", "s", false, "Interactive selection mode")
	FindCmd.Flags().StringVarP(&findGroup, "group", "g", "", "Search in specific group")
	FindCmd.Flags().IntVarP(&findLimit, "limit", "l", 0, "Limit number of results")
	FindCmd.Flags().BoolVarP(&findCopy, "copy", "c", false, "Copy selected value to clipboard")
}
