package cmd

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/livioz/password-manager/internal/vault"
	"github.com/spf13/cobra"
	"golang.design/x/clipboard"
	"golang.org/x/term"
	"log"
	"os"
	"syscall"
)

var searchTerm string

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "list all entries and select one",
	Long:  "list all entries, choose one and then decide to copy Username, Password or Other field",
	Run: func(cmd *cobra.Command, args []string) {
		err := clipboard.Init()
		if err != nil {
			log.Fatal(err)
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

		var entries []vault.VaultEntry
		if len(searchTerm) > 0 {
			entries, err = vault.SearchVaultEntries(dbPath, searchTerm)
		} else {
			entries, err = vault.ListVaultEntries(dbPath)
		}
		if err != nil {
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

		fmt.Print("Choose an ID to copy the respective password: ")
		var chosenId int
		_, err = fmt.Scanln(&chosenId)
		if err != nil {
			log.Fatal(err)
		}

		copied := false
		for _, entry := range entries {
			if entry.Id == chosenId {
				clipboard.Write(clipboard.FmtText, []byte(entry.Password))
				fmt.Printf("Password for entry '%s' copied to clipboard\n", entry.Name)
				copied = true
				break
			}
		}
		if !copied {
			fmt.Println("Please choose a valid ID")
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&searchTerm, "filter", "f", "", "specify a search term to filter entries")
}
