package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"vault/cmd"
	"vault/internal/storage"
)

var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove all expired keys",
	Run: func(cmd *cobra.Command, args []string) {
		data := storage.LoadData()
		initialCount := len(data)

		storage.CleanupExpired(data)
		storage.SaveData(data)

		removed := initialCount - len(data)

		if removed > 0 {
			fmt.Printf("Cleaning expired secrets...\n")
			fmt.Printf("Done. %d expired secret(s) removed.\n", removed)
		} else {
			fmt.Println("No expired secrets to clean.")
		}
	},
}

func init() {
	cmd.RootCmd.AddCommand(CleanCmd)
}
