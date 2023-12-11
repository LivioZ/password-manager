package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/argon2"
	"io"
	"os"
	"runtime"
)

type Params struct {
	// The amount of memory used by the algorithm (in kibibytes).
	Memory uint32

	// The number of iterations over the memory.
	Iterations uint32

	// The number of threads (or lanes) used by the algorithm.
	// Recommended value is between 1 and runtime.NumCPU().
	Parallelism uint8

	// Length of the random salt. 16 bytes is recommended for password hashing.
	SaltLength uint32

	// Length of the generated key. 16 bytes or more is recommended.
	KeyLength uint32
}

var Argon2idParams = &Params{
	Memory:      64 * 1024, // KiB * 1024
	Iterations:  4,         // aka time
	Parallelism: uint8(runtime.NumCPU()),
	SaltLength:  16,
	KeyLength:   32, // bytes
}

func InitVault(masterPassword string, keyFilePath string, dbFilePath string) error {
	masterKey, err := DeriveMasterKey(masterPassword)
	if err != nil {
		return err
	}

	key, err := GenerateEncryptionKey()
	if err != nil {
		return err
	}

	protectedKey, err := Aes256GCMEncrypt(key, masterKey)
	if err != nil {
		return err
	}

	err = os.WriteFile(keyFilePath, protectedKey, 0600)
	if err != nil {
		return err
	}

	err = CreateDB(dbFilePath)
	if err != nil {
		return err
	}

	// encrypt database file
	databasePlaintext, err := os.ReadFile(dbFilePath)
	if err != nil {
		return err
	}

	databaseCiphertext, err := Aes256GCMEncrypt(databasePlaintext, (*[32]byte)(key))
	if err != nil {
		return err
	}

	err = os.WriteFile(dbFilePath, databaseCiphertext, 0600)
	if err != nil {
		return err
	}

	return nil
}

func DeriveMasterKey(masterPassword string) (*[32]byte, error) {
	// 256-bit master key derivation with argon2id from the master password
	masterKey := argon2.IDKey(
		[]byte(masterPassword),
		[]byte(""),
		Argon2idParams.Iterations,
		Argon2idParams.Memory,
		Argon2idParams.Parallelism,
		Argon2idParams.KeyLength,
	)

	return (*[32]byte)(masterKey), nil
}

func GenerateEncryptionKey() ([]byte, error) {
	encryptionKey := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, encryptionKey)
	if err != nil {
		return nil, err
	}

	return encryptionKey, nil
}

func Aes256GCMEncrypt(plaintext []byte, key *[32]byte) (ciphertext []byte, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func Aes256GCMDecrypt(ciphertext []byte, key *[32]byte) (plainText []byte, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
}

func DeriveToVaultKey(masterPassword []byte, keyPath string) (*[32]byte, error) {
	masterKey, err := DeriveMasterKey(string(masterPassword))
	if err != nil {
		return nil, errors.New("error reading master password from input")
	}

	// decrypt protected encryption key
	encryptedKey, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("error reading key file ('%s')", keyPath)
	}
	key, err := Aes256GCMDecrypt(encryptedKey, masterKey)
	if err != nil {
		return nil, errors.New("error decrypting vault key (wrong master password or tampered key file)")
	}

	return (*[32]byte)(key), nil
}

func UnlockVault(dbPath string, key *[32]byte) error {
	// decrypt database file
	databaseCiphertext, err := os.ReadFile(dbPath)
	if err != nil {
		return err
	}
	databasePlaintext, err := Aes256GCMDecrypt(databaseCiphertext, key)
	if err != nil {
		return err
	}

	err = os.WriteFile(dbPath, databasePlaintext, 0600)
	if err != nil {
		return err
	}

	return nil
}

func LockVault(dbPath string, key *[32]byte) error {
	// encrypt database file
	databasePlaintext, err := os.ReadFile(dbPath)
	if err != nil {
		return err
	}
	databaseCiphertext, err := Aes256GCMEncrypt(databasePlaintext, key)
	if err != nil {
		return err
	}

	err = os.WriteFile(dbPath, databaseCiphertext, 0600)
	if err != nil {
		return err
	}

	return nil
}
