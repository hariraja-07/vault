package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "v0.3.0"
)

// RootCmd is the root command
var RootCmd = &cobra.Command{
	Use:   "vault",
	Short: "vault - CLI data manager for developers",
	Long: `vault is a simple CLI tool for storing and managing secrets.
Store, retrieve, and organize your data right from the terminal.`,
	Version: version,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
