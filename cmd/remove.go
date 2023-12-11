package cmd

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/livioz/password-manager/internal/vault"
	"golang.org/x/term"
	"log"
	"os"
	"syscall"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an entry from vault",
	Long:  "Delete an entry from vault.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		exitIfVaultDoesNotExist()
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Insert master password: ")
		masterPassword, err := term.ReadPassword(int(syscall.Stdin))
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

		var entries []vault.VaultEntry
		if len(searchTerm) > 0 {
			entries, err = vault.SearchVaultEntries(dbPath, searchTerm)
		} else {
			entries, err = vault.ListVaultEntries(dbPath)
		}
		if err != nil {
			errLock := vault.LockVault(dbPath, vaultKey)
			if errLock != nil {
				log.Println(err)
			}
			log.Fatalln(err)
		}

		err = vault.LockVault(dbPath, vaultKey)
		if err != nil {
			log.Fatalln(err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"id", "Name", "Username", "Password", "Other"})
		rows := make([]table.Row, 0)
		for _, entry := range entries {
			rows = append(rows, table.Row{entry.Id, entry.Name, entry.Username, "********", entry.Other})
		}
		t.AppendRows(rows)
		t.AppendSeparator()
		t.Render()

		fmt.Print("Choose an ID to delete the respective entry: ")
		var chosenId int
		_, err = fmt.Scanln(&chosenId)
		if err != nil {
			log.Fatal(err)
		}

		deleted := false
		for _, entry := range entries {
			if entry.Id == chosenId {
				err = vault.UnlockVault(dbPath, vaultKey)
				if err != nil {
					log.Fatalln(err)
				}
				err = vault.DeleteEntry(dbPath, chosenId)
				if err != nil {
					errLock := vault.LockVault(dbPath, vaultKey)
					if errLock != nil {
						log.Println(err)
					}
					log.Fatalln(err)
				}
				err = vault.LockVault(dbPath, vaultKey)
				if err != nil {
					log.Fatalln(err)
				}
				fmt.Printf("Entry with name '%s' deleted from vault\n", entry.Name)
				deleted = true
				break
			}
		}
		if !deleted {
			fmt.Println("Please choose a valid ID")
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVarP(&searchTerm, "filter", "f", "", "specify a search term to filter entries")
}
