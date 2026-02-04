package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const keychainService = "com.rayui.encryption"
const keychainAccount = "master-key"

// DeriveKey retrieves or creates a 256-bit key from the platform keystore.
func DeriveKey() ([]byte, error) {
	switch runtime.GOOS {
	case "darwin":
		return deriveKeyDarwin()
	case "windows":
		return deriveKeyWindows()
	default:
		return deriveKeyLinux()
	}
}

// deriveKeyDarwin uses macOS Keychain via the `security` CLI.
func deriveKeyDarwin() ([]byte, error) {
	// Try to read existing key.
	out, err := exec.Command("security", "find-generic-password",
		"-s", keychainService, "-a", keychainAccount, "-w").Output()
	if err == nil {
		raw := strings.TrimSpace(string(out))
		h := sha256.Sum256([]byte(raw))
		return h[:], nil
	}

	// Generate a new random secret and store it.
	secret := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, secret); err != nil {
		return nil, fmt.Errorf("generate random key: %w", err)
	}
	encoded := base64.StdEncoding.EncodeToString(secret)

	if err := exec.Command("security", "add-generic-password",
		"-s", keychainService, "-a", keychainAccount, "-w", encoded, "-U").Run(); err != nil {
		return nil, fmt.Errorf("store key in keychain: %w", err)
	}

	h := sha256.Sum256([]byte(encoded))
	return h[:], nil
}

// deriveKeyWindows uses a deterministic fallback based on machine-id.
// A proper implementation would use DPAPI via syscall.
func deriveKeyWindows() ([]byte, error) {
	return deriveKeyFromMachineID()
}

// deriveKeyLinux uses machine-id as a fallback.
func deriveKeyLinux() ([]byte, error) {
	return deriveKeyFromMachineID()
}

func deriveKeyFromMachineID() ([]byte, error) {
	paths := []string{"/etc/machine-id", "/var/lib/dbus/machine-id"}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err == nil {
			h := sha256.Sum256(data)
			return h[:], nil
		}
	}
	// Final fallback: hostname.
	name, err := os.Hostname()
	if err != nil {
		return nil, errors.New("unable to derive encryption key: no machine-id or hostname")
	}
	h := sha256.Sum256([]byte("rayui:" + name))
	return h[:], nil
}

// Encrypt encrypts plaintext with AES-256-GCM and returns a base64 string.
func Encrypt(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a base64 AES-256-GCM ciphertext.
func Decrypt(ciphertext string, key []byte) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, sealed := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, sealed, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
