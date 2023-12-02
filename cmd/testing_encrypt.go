package cmd

import (
	"fmt"
	"github.com/livioz/password-manager/internal/vault"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"log"
	"syscall"
)

var encryptCmd = &cobra.Command{
	Use:   "lock",
	Short: "Encrypt the vault",
	Long:  "The database file gets encrypted with AES-256 GCM using the master password",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Insert master password: ")
		masterPassword, err := term.ReadPassword(syscall.Stdin)
		fmt.Println()
		if err != nil {
			log.Fatalln(err)
		}

		vaultKey, err := vault.DeriveToVaultKey(masterPassword, keyPath)
		if err != nil {
			log.Fatalln(err)
		}
		err = vault.LockVault(dbPath, vaultKey)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("Vault locked successfully")
	},
}

func init() {
	// for testing purposes only
	// rootCmd.AddCommand(encryptCmd)
}
