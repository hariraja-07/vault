package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"vault/internal/crypto"
	"vault/internal/storage"
)

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

		var val interface{}

		if strings.Contains(key, "/") {
			if strings.Count(key, "/") > 1 {
				fmt.Println("Error: Only one level grouping allowed (group/key)")
				return
			}

			parts := strings.SplitN(key, "/", 2)
			group := parts[0]
			subKey := parts[1]

			groupMap, ok := data[group].(map[string]interface{})
			if !ok {
				fmt.Println("Group not found")
				return
			}
			val, ok = groupMap[subKey]
			if !ok {
				fmt.Println("Key not found")
				return
			}
		} else {
			var ok bool
			val, ok = data[key]
			if !ok {
				fmt.Println("Key not found")
				return
			}
		}

		// Check if value is encrypted
		if encryptedVal, ok := val.(map[string]interface{}); ok {
			password := askPassword("Enter password: ")
			ev := &crypto.EncryptedValue{
				Ciphertext: encryptedVal["ciphertext"].(string),
				Nonce:      encryptedVal["nonce"].(string),
			}
			decrypted, err := crypto.Decrypt(ev, password)
			if err != nil {
				fmt.Println("Error: Decryption failed. Wrong password?")
				return
			}
			fmt.Println(decrypted)
		} else {
			fmt.Println(val)
		}

		storage.TrackKeyUsage(key)
	},
}
