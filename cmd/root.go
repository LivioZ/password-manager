package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "password-manager",
	Short: "Password manager written in Go",
	Long:  "Password manager written in Go, use command 'init' to get started",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}

func exitIfVaultDoesNotExist() {
	// check if vault does not exist in current directory
	if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Vault doesn't exist in current directory\nRun `password-manager init` to initialize vault")
		os.Exit(1)
	}
}
