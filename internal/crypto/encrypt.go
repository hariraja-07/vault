package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

// deriveKey creates a 32-byte key from a password using SHA-256
func deriveKey(password string) ([]byte, error) {
	hash := sha256.Sum256([]byte(password))
	return hash[:], nil
}

// EncryptedValue represents an encrypted value in storage
type EncryptedValue struct {
	Ciphertext string `json:"ciphertext"`
	Nonce      string `json:"nonce"`
}

// Encrypt encrypts plaintext using AES-256-GCM with the given password
func Encrypt(plaintext, password string) (*EncryptedValue, error) {
	// Derive key from password
	key, err := deriveKey(password)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %w", err)
	}

	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	return &EncryptedValue{
		Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
		Nonce:      base64.StdEncoding.EncodeToString(nonce),
	}, nil
}

// Decrypt decrypts an encrypted value using the given password
func Decrypt(ev *EncryptedValue, password string) (string, error) {
	// Derive key from password
	key, err := deriveKey(password)
	if err != nil {
		return "", fmt.Errorf("failed to derive key: %w", err)
	}

	// Decode ciphertext and nonce
	ciphertext, err := base64.StdEncoding.DecodeString(ev.Ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	nonce, err := base64.StdEncoding.DecodeString(ev.Nonce)
	if err != nil {
		return "", fmt.Errorf("failed to decode nonce: %w", err)
	}

	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// IsEncrypted checks if a value is an encrypted value (map with ciphertext/nonce)
func IsEncrypted(value interface{}) bool {
	_, ok := value.(map[string]interface{})
	return ok
}
