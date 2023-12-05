package cmd

import (
	"fmt"
	"github.com/livioz/password-manager/internal/vault"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"log"
	"syscall"
)

var username string
var other string

var addCmd = &cobra.Command{
	Use:   "add entry-name",
	Short: "Add an entry in the database with the specified name and optional fields",
	Long:  "Add entry in the database. Specify fields using flags",
	Args:  cobra.ExactArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		exitIfVaultDoesNotExist()
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		var password string

		fmt.Print("Insert new entry's password: ")
		_, err := fmt.Scanln(&password)
		if err != nil {
			log.Fatalln(err)
		}

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

		err = vault.AddEntry(dbPath, name, username, password, other)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("Entry '%s' added successfully\n", name)

		err = vault.LockVault(dbPath, vaultKey)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&username, "username", "u", "", "username or email")
	addCmd.Flags().StringVarP(&other, "other", "o", "", "anything else")
}
