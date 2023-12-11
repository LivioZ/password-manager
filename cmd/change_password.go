package cmd

import (
	"bytes"
	"fmt"
	"github.com/livioz/password-manager/internal/vault"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"log"
	"os"
	"syscall"
	"unicode/utf8"
)

var changePWD = &cobra.Command{
	Use:   "change-pwd [key-path]",
	Short: "Change master password",
	Long:  "Change the master password by re-encrypting the key file",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		exitIfVaultDoesNotExist()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 2 {
			keyPath = args[1]
		}
		fmt.Print("Insert current master password: ")
		currentMasterPassword, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Input the new master password: ")
		newMasterPassword, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			log.Fatal(err)
		}
		if !utf8.Valid(newMasterPassword) {
			log.Fatal("The master password is not a valid UTF-8 encoded string")
		}
		if bytes.Equal(currentMasterPassword, newMasterPassword) {
			log.Fatal("The current master password and the new master password are the same.")
		}

		fmt.Printf("Repeat the new master password: ")
		newMasterPasswordRepeat, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			log.Fatal(err)
		}
		if !bytes.Equal(newMasterPassword, newMasterPasswordRepeat) {
			log.Fatal("The inserted passwords do not match\n")
		}

		// actual master password change flow
		currentMasterKey, err := vault.DeriveMasterKey(string(currentMasterPassword))
		if err != nil {
			log.Fatal(err)
		}

		encryptedKey, err := os.ReadFile(keyPath)
		if err != nil {
			log.Fatal(err)
		}
		key, err := vault.Aes256GCMDecrypt(encryptedKey, currentMasterKey)
		if err != nil {
			log.Fatal("The current master password is not correct")
		}

		newMasterKey, err := vault.DeriveMasterKey(string(newMasterPassword))
		if err != nil {
			log.Fatal(err)
		}

		protectedKey, err := vault.Aes256GCMEncrypt(key, newMasterKey)
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(keyPath, protectedKey, 0600)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Master password change successfully")
	},
}

func init() {
	rootCmd.AddCommand(changePWD)
}
