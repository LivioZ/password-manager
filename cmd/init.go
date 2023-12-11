package cmd

import (
	"bytes"
	"fmt"
	"github.com/livioz/password-manager/internal/vault"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"log"
	"syscall"
	"unicode/utf8"
)

var dbPath = "vault.db"
var keyPath = "key.bin"

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize vault in current working directory",
	Long:  "Initialize vault by choosing the master password. It creates the database and encryption key files",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Input the master password: ")
		masterPassword, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			log.Fatal(err)
		}
		if !utf8.Valid(masterPassword) {
			log.Fatal("The master password is not a valid UTF-8 encoded string")
		}

		fmt.Printf("Repeat the master password: ")
		masterPasswordRepeat, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			log.Fatal(err)
		}
		if !bytes.Equal(masterPassword, masterPasswordRepeat) {
			log.Fatal("The inserted passwords do not match\n")
		}

		err = vault.InitVault(string(masterPassword), keyPath, dbPath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
