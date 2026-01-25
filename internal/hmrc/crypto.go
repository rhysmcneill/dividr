package hmrc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// Encrypt encrypts a plain text string using AES-GCM and returns a base64 encoded string.
// The key must be exactly 32 bytes (for AES-256).
func Encrypt(plaintext string, keyString string) (string, error) {
	if len(keyString) != 32 {
		return "", errors.New("encryption key must be 32 bytes")
	}
	key := []byte(keyString)

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

	// Seal: encrypts and authenticates
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Return as Base64 string for easy storage in Postgres TEXT columns
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt takes a base64 encoded ciphertext and returns the original string.
func Decrypt(cryptoText string, keyString string) (string, error) {
	if len(keyString) != 32 {
		return "", errors.New("encryption key must be 32 bytes")
	}
	key := []byte(keyString)

	data, err := base64.StdEncoding.DecodeString(cryptoText)
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

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
