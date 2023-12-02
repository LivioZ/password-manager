package cmd

import (
	"fmt"
	"github.com/livioz/password-manager/internal/vault"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"log"
	"syscall"
)

var decryptCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Decrypt the vault",
	Long:  "The database file gets decrypted with AES-256 GCM using the master password",
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
		err = vault.UnlockVault(dbPath, vaultKey)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("Vault unlocked successfully")
	},
}

func init() {
	// for testing purposes only
	// rootCmd.AddCommand(decryptCmd)
}
